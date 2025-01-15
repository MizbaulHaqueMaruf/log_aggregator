package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mizbaulhaquemaruf/log_aggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf(("usage: %s <name> <url>"), cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}

	fmt.Printf("feed created: %v\n", feed)

	fmt.Printf("User: %v\n", feedFollow.UserName)
	fmt.Printf("Feed: %v\n", feedFollow.FeedName)

	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf(("usage: %s <name>"), cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}

	for _, feed := range feeds {
		printFeed(feed)
		fmt.Println("===================================")
	}

	return nil
}

func printFeed(feed database.GetFeedsRow) {
	fmt.Printf("feed name: %v\n", feed.Name)
	fmt.Printf("feed url: %v\n", feed.Url)
	fmt.Printf("feed user: %v\n", feed.Username)
}
