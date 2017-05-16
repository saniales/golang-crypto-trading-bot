package main

import "github.com/AlessandroSanino1994/gobot/botHelpers"

const (
	EnvironmentError = iota //Error thrown when environment variables are not set properly.
	NetworkError     = iota //Error thrown when there are problems while connecting to external services.
	AuthError        = iota //Error thrown when authenticating to external services.
)

func main() {

	botHelpers.InitExchange("bittrex")

	for {

	}

}
