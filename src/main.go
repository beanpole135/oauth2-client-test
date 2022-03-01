package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"encoding/json"
)

type Settings struct {
	ClientID string		`json:"client_id"`
	ClientSecret string	`json:"client_secret"`
	Scopes []string		`json:"scopes"`
	Endpoint_authurl string		`json:"endpoint_auth_url"`
	Endpoint_tokenurl string		`json:"endpoint_token_url"`
	UserAPI string		`json:"user_api_url"`
}

func main() {
	err := loadSettings()
	if err != nil {
		log.Printf("Error loading config.json: %v", err)
		os.Exit(1)
	}
	mux := http.NewServeMux()
	// Root
	mux.Handle("/",  http.FileServer(http.Dir("templates/")))

	// OauthGoogle
	mux.HandleFunc("/auth/login", oauthLogin)
	mux.HandleFunc("/auth/callback", oauthCallback)

	server := &http.Server{
		Addr: fmt.Sprintf(":8000"),
		Handler: mux,
	}

	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Server Error: %v", err)
	} else {
		log.Println("Server closed!")
	}
}

func loadSettings() error {
	var S Settings
	//Load the file
	dat, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(dat, &S)
	if err != nil {
		return err
	}
	//Now use the settings
	generateOauthConfig(S)
	return nil
}