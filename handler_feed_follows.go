package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mizbaulhaquemaruf/log_aggregator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}

	ffRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}

	fmt.Println("feed followed")
	fmt.Printf("User: %v\n", ffRow.UserName)
	fmt.Printf("Feed: %v\n", ffRow.FeedName)

	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feed follows: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("no feeds followed")
		return nil
	}

	fmt.Printf("feeds followed by %v:\n", user.Name)
	for _, follow := range follows {
		fmt.Printf("  %v\n", follow.FeedName)
	}

	return nil

}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    cmd.Args[0],
	})
	if err != nil {
		return fmt.Errorf("error unfollowing feed: %w", err)
	}

	fmt.Println("feed unfollowed")
	return nil
}
