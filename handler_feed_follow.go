package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yumski/rss_aggregator/internal/database"
)

func (cfg *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, "Unable to parse json")
		return
	}

	feedFollow, err := cfg.DB.CreateFollowFeed(r.Context(), database.CreateFollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Unable to follow feed: %v", err))
		return
	}

	respondWithJson(w, 200, databaseFeedFollowstoFeedFollows(feedFollow))
}

func (cfg *apiConfig) handlerGetUserFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFollowFeeds(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Unable authenticate user: %v", err))
		return
	}

	respondWithJson(w, 200, databaseFeedsFollowstoFeedsFollows(feedFollows))
}

func (cfg *apiConfig) handlerDeleteUserFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDstr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDstr)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow id: %v", err))
		return
	}

	err = cfg.DB.DeleteFollowFeeds(r.Context(), database.DeleteFollowFeedsParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete feed fllow: %v", err))
		return
	}

	respondWithJson(w, 200, struct{}{})
}
