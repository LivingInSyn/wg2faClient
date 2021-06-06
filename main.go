package main

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// base config variables
// TODO: move to env/cli args
var BASE_DOMAIN = "https://dev-318981.oktapreview.com"
var DOMAIN = BASE_DOMAIN + "/oauth2/default/v1'"
var REDIRECT_URL = "http://localhost:8080/implicit/callback"
var CLIENT_ID = "0oawepftsdT43o2CM0h7"

// this is the channel for postToken to say we've gotten a valid token
var tokenChannel chan bool // <- Create var for channel

// struct for the incoming token

func postToken(w http.ResponseWriter, req *http.Request) {
	tokenChannel <- true
}

func startWebserver(wg *sync.WaitGroup) *http.Server {
	r := mux.NewRouter()
	// setup handling of static pages
	r.Path("/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	// token handler
	r.HandleFunc("/token", postToken).Methods("POST")
	// start and return server
	srv := &http.Server{
		Addr: "127.0.0.1:8091",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}
	log.Debug().Msg("Starting http server")
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatal().AnErr("err", err).Msg("error starting web server")
		}
	}()
	return srv
}

// func checkWg2fa() {

// }

func openBrowser(url string) {
	//credit: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("error opening up web browser")
	}
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// init the token channel
	tokenChannel = make(chan bool)
	//start the web server
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	webserver := startWebserver(httpServerExitDone)
	// build secrets and challenges
	verifier, err := TokenUrlSafe(32)
	if err != nil {
		log.Fatal().AnErr("err", err).Msg("couldn't generate random numbers")
	}
	challenge := GetChallenge(verifier)
	state, err := TokenUrlSafe(32)
	if err != nil {
		log.Fatal().AnErr("err", err).Msg("couldn't generate random numbers")
	}
	lparams := LoginParams{
		Domain:        DOMAIN,
		RedirectURI:   REDIRECT_URL,
		ClientID:      CLIENT_ID,
		CodeChallenge: challenge,
		State:         state,
	}
	url := GetLoginUrl(&lparams)
	openBrowser(url)
	// wait for our token from the browser
	good := <-tokenChannel
	log.Debug().Str("token channel", strconv.FormatBool(good)).Msg("output from token channel")
	if err := webserver.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}
