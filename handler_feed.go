package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yumski/rss_aggregator/internal/database"
)

func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	type response struct {
		Feed        `json:"feed"`
		FeedFollows `json:"feed_follow"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Unable to create feed: %v", err))
		return
	}

	feedFollow, err := cfg.DB.CreateFollowFeed(r.Context(), database.CreateFollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	res := response{
		Feed:        databaseFeedtoFeed(feed),
		FeedFollows: databaseFeedFollowstoFeedFollows(feedFollow),
	}

	respondWithJson(w, 200, res)
}

func (cfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedstoFeeds(feeds))
}
