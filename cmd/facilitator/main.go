package main

import (
	"net/http"

	"github.com/rabbitprincess/x402-facilitator/api"
	"github.com/rs/zerolog/log"
)

func main() {
	s := api.NewServer()

	if err := http.ListenAndServe(":8080", s); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
