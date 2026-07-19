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
	fmt.Printf("Username: %s has been set", username)
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
		fmt.Printf("User: %s was registered", data.Name)
	} else {
		return fmt.Errorf("unexpected db error: %w", err)
	}
	return nil
}
