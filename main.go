package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var BASE_DOMAIN = "https://dev-318981.oktapreview.com"
var DOMAIN = BASE_DOMAIN + "/oauth2/default/v1'"
var REDIRECT_URL = "http://localhost:8080/implicit/callback"
var CLIENT_ID = "0oawepftsdT43o2CM0h7"

// struct for the incoming token

func postToken(w http.ResponseWriter, req *http.Request) {

}

func handleWebserver() {
	// TODO: https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve
	// handle intentional shutdown of this method
	r := mux.NewRouter()
	// setup handling of static pages
	// https://www.alexedwards.net/blog/serving-static-sites-with-go
	// TODO: handle the headers in the GETs to the static dir, need to set Access-Control-Allow-Origin
	fs := http.FileServer(http.Dir("./static"))
	r.HandleFunc("/", fs)
	// token handler
	r.HandleFunc("/token").Methods("POST")
	srv := &http.Server{
		Addr: "127.0.0.1:8091",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}
	log.Debug().Msg("Starting http server")
	log.Fatal().Msg(srv.ListenAndServe().Error())
}

func checkWg2fa() {

}

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
		log.Fatal(err)
	}
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	//start the web server
	go handleWebserver()
	// build secrets and challenges
	// verifier :=
	// challenge :=
	// state :=
	// url :=

	// openBrowser(url)
	// wait for
}
