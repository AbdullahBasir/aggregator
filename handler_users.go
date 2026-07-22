package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/AbdullahBasir/aggregator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
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
	if len(cmd.args) != 1 {
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

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	if len(cmd.args) > 0 {
		newlimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("could not convert string to int: %w", err)
		}
		limit = newlimit
	} else {
		limit = 2
	}
	userPost, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("could not get posts for user: %w", err)
	}
	for _, post := range userPost {
		title := "no title"
		if post.Title.Valid {
			title = post.Title.String
		}
		publish := time.Time{}
		if post.PublishedAt.Valid {
			publish = post.PublishedAt.Time
		}
		description := "no description"
		if post.Description.Valid {
			description = post.Description.String
		}
		fmt.Printf("--- %v ---\n* %s\n* %v\n%v\n", title, post.Url, publish, description)
	}
	return nil
}
