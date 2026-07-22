package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AbdullahBasir/aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("command must include feed name and url")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not store feed in database: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    feed.UserID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed follow: %w", err)
	}

	fmt.Printf("* %s\n", feedFollow.FeedName)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("list feeds does not take any arguments")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get feed creator: %w", err)
	}

	fmt.Printf("Feeds have been listed\n")
	for _, feed := range feeds {
		fmt.Printf("* %s\n  URL: %s\n  User: %s\n", feed.Name_2, feed.Url, feed.Name)
	}
	return nil
}
