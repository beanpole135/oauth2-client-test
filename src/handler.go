package main

import (
	"golang.org/x/oauth2"
	"net/http"
	"fmt"
	"io/ioutil"
	"context"
	"log"
	"encoding/base64"
	"encoding/json"
	"crypto/rand"
	"crypto/tls"
	"time"
	"html/template"
	"bytes"
)

// Scopes: OAuth 2.0 scopes provide a way to limit the amount of access that is granted to an access token.
var oauthConfig *oauth2.Config
var oauthUrlAPI string = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var templates *template.Template
//Populate oauth config from settings
func generateOauthConfig(S Settings) {
	oauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8000/auth/callback",
		ClientID:     S.ClientID,
		ClientSecret: S.ClientSecret,
		Scopes:       S.Scopes,
		Endpoint:     oauth2.Endpoint{
			AuthURL:	S.Endpoint_authurl,
			TokenURL:	S.Endpoint_tokenurl,
			AuthStyle:	oauth2.AuthStyleAutoDetect,
		},
	}
	oauthUrlAPI = S.UserAPI
	templates = template.Must(template.ParseFiles("templates/error.html", "templates/success.html"))
}

func oauthLogin(w http.ResponseWriter, r *http.Request) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)

	/*
	AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
	validate that it matches the the state query parameter on your redirect callback.
	*/
	u := oauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func oauthCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")
	var err error
	if r.FormValue("state") != oauthState.Value {
		err = fmt.Errorf("invalid oauth state")
	}
	if err == nil {
		errstring := r.FormValue("error")
		errDesc := r.FormValue("error_description")
		if errstring != "" {
			err = fmt.Errorf("OauthError: "+ errstring + "\n" + errDesc)
		}
	}
	var data []byte
	if err == nil {
		data, err = getUserData(r.FormValue("code"))
	}
	fmt.Println("Got Form replies:", r.Form)
	if err != nil {
		log.Println(err.Error())
		templates.ExecuteTemplate(w, "error.html", err.Error())
		//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// GetOrCreate User in your db.
	// Redirect or response with a token.
	//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

	//Simple display of info for testing purposes
	var js bytes.Buffer
	errJ := json.Indent(&js, data, "", "  ")
	if errJ == nil  {
		templates.ExecuteTemplate(w, "success.html", string(js.Bytes()))
	}else{
		templates.ExecuteTemplate(w, "success.html", string(data))
	}

}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserData(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.
 	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    sslcli := &http.Client{Transport: tr}
    ctx := context.TODO()
    ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := sslcli.Get( fmt.Sprintf(oauthUrlAPI, token.AccessToken) )
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}