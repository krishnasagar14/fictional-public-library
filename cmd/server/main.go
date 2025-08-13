package main

import (
	"fictional-public-library/logging"
	"fictional-public-library/server"
)

func main() {
	err := server.RunServer()
	if err != nil {
		logging.Log.WithField("Error", err).Fatal("Server quit due to error")
	}
}
