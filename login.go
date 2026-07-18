package main

import (
	"errors"
	"fmt"

	"github.com/AbdullahBasir/aggregator/internal/config"
)

type state struct {
	config *config.Config
}

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
	err := s.config.SetUser(username)
	if err != nil {
		return fmt.Errorf("could not set username: %w", err)
	}
	fmt.Printf("Username: %s has been set", username)
	return nil
}
