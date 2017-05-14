package main

import (
	"fmt"
	"os"
)

const (
	EnvironmentError = iota //Error thrown when environment variables are not set properly.
	NetworkError     = iota //Error thrown when there are problems while connecting to external services.
	AuthError        = iota //Error thrown when authenticating to external services.
)

func main() {
	exchange := os.Getenv("exchange")
	algorithm := os.Getenv("algorithm")
	key := os.Getenv("key")
	secret := os.Getenv("secret")

	if exchange == "" {
		exchange = "bittrex"
	}
	if algorithm == "" {
		algorithm = "macd"
	}
	if key == "" || secret == "" {
		fmt.Println("Env variables for public and secret key not found. Exiting...")
		os.Exit(1)
	}

	for {

	}

}
