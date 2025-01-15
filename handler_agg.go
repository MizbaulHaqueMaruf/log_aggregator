package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mizbaulhaquemaruf/log_aggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing time between requests: %w", err)
	}

	log.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("error getting next feed to fetch: %w", err)
		return
	}

	log.Printf("Found a feed: %v\n", feed)
	scrapeFeed(s.db, feed)

}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("error marking feed fetched: %v", err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("error fetching feed: %v", err)
		return
	}

	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("error creating post: %v", err)
			continue
		}
	}

	log.Printf("Feed %s fetched, %v posts found.", feed.Name, len(feedData.Channel.Item))
}
