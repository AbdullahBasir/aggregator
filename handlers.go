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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("command must include feed name and url")
	}
	username := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("could not get user from database: %w", err)
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
	fmt.Printf("%+v\n", feed)
	return nil
}
