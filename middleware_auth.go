package main

import (
	"fmt"
	"net/http"

	"github.com/yumski/rss_aggregator/internal/auth"
	"github.com/yumski/rss_aggregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

// utilize an anonmyous function as an http handler to be able to access the database.User variable
func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
