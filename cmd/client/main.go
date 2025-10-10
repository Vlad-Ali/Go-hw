package main

import (
	"encoding/base64"
	"fmt"
	"httpapp/internal/client/clientapp"
)

var (
	baseURL = "http://127.0.0.1:8080"
)

func main() {
	clientApp := clientapp.NewClientApp(baseURL)

	if err := clientApp.GetVersionEndpoint(); err != nil {
		fmt.Printf("Error getting version endpoint: %s\n", err)
	}

	if err := clientApp.PostDecodeEndpoint(base64.StdEncoding.EncodeToString([]byte("Hello world!"))); err != nil {
		fmt.Printf("Error posting decode endpoint: %s\n", err)
	}

	if err := clientApp.GetHardOpEndpoint(); err != nil {
		fmt.Printf("Error getting hard op endpoint: %s\n", err)
	}

}
