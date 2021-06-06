package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// LoginParams contains the fields required to generate
// a loginurl
type LoginParams struct {
	Domain        string
	RedirectURI   string
	ClientID      string
	CodeChallenge string
	State         string
}

func (l *LoginParams) encodeUrl() string {
	return url.QueryEscape(l.RedirectURI)
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// TokenUrlSafe returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func TokenUrlSafe(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// GetChallenge returns a PKCS challege given a verifier
func GetChallenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))
	encoded := base64.URLEncoding.EncodeToString(h.Sum(nil))
	encoded = strings.ReplaceAll(encoded, "=", "")
	return encoded
}

func GetLoginUrl(l *LoginParams) string {
	/*
			def get_login_url(self, challenge, state):
		        eurl = quote_plus(self.redirect_url)
		        url = f"https://{self.domain}/authorize?"
		        url = url + "response_type=code"
		        url = url + f"&client_id={self.client_id}"
		        url = url + f"&redirect_uri={eurl}"
		        url = url + f"&response_mode=fragment"
		        url = url + f"&code_challenge={challenge}"
		        url = url + f"&code_challenge_method=S256"
		        #url = url + "&scope=openid%20profile%20email&audience=appointments:api&state=xyzABC123"
		        url = url + f"&scope=openid%20profile%20email&state={state}"
		        return url
	*/
	url_template := "https://%s/authorize?response_type=code&client_id=%s&redirect_uri=%s&response_mode=fragment&code_challenge=%s&code_challenge_method=S256&scope=openid%%20profile%%20email&state=%s"
	url_template = fmt.Sprintf(url_template, l.Domain, l.ClientID, l.encodeUrl(), l.CodeChallenge, l.State)
	return url_template
}
