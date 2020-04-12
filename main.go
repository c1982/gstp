package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
)

var (
	configfilepath = flag.String("c", "./config.yaml", "pge config file path")
)

func main() {

	flag.Parse()

	config, err := readConfig(*configfilepath)
	if err != nil {
		log.Fatal(err)
	}

service:
	srv, err := NewGmailService(config.CredentialFile, config.TokenFile, config.Query)
	if err == ErrCredentialFileExists {
		setupCredential()
	} else if err == ErrTokenFileExists {
		setUpToken(srv.config, srv.tokenFile)
		goto service
	} else if err != nil {
		log.Fatal(err)
	}

	run(srv, config)
}

func readConfig(configfilepath string) (config *GstpConfig, err error) {

	cdata, err := ReadFile(configfilepath)
	if err != nil {
		return config, err
	}

	config, err = UnMarshallConfig(cdata)
	if err != nil {
		return config, err
	}

	return config, nil
}

func setupCredential() {
	fmt.Println("error:")
	fmt.Println("\tcredential.json file not found")
	fmt.Println("\tGo to the Google Developers API Console: https://console.developers.google.com/apis/credentials")
	fmt.Println("\tand download credential json file in this directory.")
	os.Exit(1)
}

func setUpToken(config *oauth2.Config, tokenPath string) {

	var code string

	fmt.Println()
	fmt.Println("1. Sign in to your browser with the gmail account you want to track")
	fmt.Println()
	fmt.Println("2. Open the following link and authorise gmail account:")
	fmt.Println(config.AuthCodeURL("pge", oauth2.AccessTypeOffline))
	fmt.Println()
	fmt.Println("3. Enter the authorisation code:")

	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Failed to read authorisation code: %v.", err)
	}

	fmt.Println()

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Failed to exchange authorisation code for token: %v.", err)
	}

	tokenFile, err := os.OpenFile(tokenPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Failed to open token file for writing: %v.", err)
	}

	defer tokenFile.Close()
	if err := json.NewEncoder(tokenFile).Encode(token); err != nil {
		log.Fatalf("Failed to write token: %v.", err)
	}
}
