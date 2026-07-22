package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AbdullahBasir/aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("could not find url")
	}
	feed, err := s.db.GetFeedWithUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not get feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed follow: %w", err)
	}
	fmt.Printf("Feed: %s\nUser: %s\n", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 0 {
		return errors.New("following does not take any arguments")
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("could not get follows for user: %w", err)
	}

	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("could not find url")
	}
	feed, err := s.db.GetFeedWithUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not get feed: %w", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not delete feed follow: %w", err)
	}
	fmt.Printf("unfollowed: %s\n", cmd.args[0])
	return nil
}
