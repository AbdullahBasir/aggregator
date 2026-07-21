package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/AbdullahBasir/aggregator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("could not find username")
	}
	username := cmd.args[0]
	if username == "" {
		return errors.New("could not get username from command")
	}
	_, err := s.db.GetUser(context.Background(), username)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("user %s does not exist", username)
	} else if err != nil {
		return fmt.Errorf("could not get user: %w", err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("could not set username: %w", err)
	}
	fmt.Printf("Username: %s has been set\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("could not find username")
	}
	name := cmd.args[0]
	if name == "" {
		return errors.New("could not get username from command")
	}
	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return errors.New("user already registered")
	} else if errors.Is(err, sql.ErrNoRows) {
		data, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
		})
		if err != nil {
			return fmt.Errorf("could not create user: %w", err)
		}

		err = s.cfg.SetUser(data.Name)
		if err != nil {
			return fmt.Errorf("could not set user data: %w", err)
		}
		fmt.Printf("User: %s was registered\n", data.Name)
	} else {
		return fmt.Errorf("unexpected db error: %w", err)
	}
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	data, err := s.db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not list users: %w", err)
	}
	fmt.Printf("Users have been listed\n")
	for _, name := range data {
		if name.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", name.Name)
		} else {
			fmt.Printf("* %s\n", name.Name)
		}
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUser(context.Background())
	if err != nil {
		return fmt.Errorf("could not reset user: %w", err)
	}
	fmt.Printf("Users have been reset\n")
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("could not fetch the feed: %w", err)
	}
	fmt.Printf("%+v\n", feed)
	return nil
}

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
		return errors.New("bad input no command args needed")
	}
	users, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get feed creator: %w", err)
	}

	fmt.Printf("Feeds have been listed\n")
	for _, user := range users {
		fmt.Printf("* %s\n* %s\n* %s\n", user.Name_2, user.Url, user.Name)
	}
	return nil
}

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
	fmt.Printf("* %s\n* %s\n", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 0 {
		return errors.New("bad input only command arg needed")
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
	if len(cmd.args) < 1 {
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
	fmt.Printf("unfollowed: %s", cmd.args[0])
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		name := s.cfg.CurrentUserName
		if name == "" {
			return fmt.Errorf("no current user name in config")
		}

		user, err := s.db.GetUser(context.Background(), name)
		if err != nil {
			return fmt.Errorf("could not get user: %w", err)
		}
		return handler(s, cmd, user)
	}
}
