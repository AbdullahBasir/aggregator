package main

import (
	"errors"
	"fmt"
)

type commands struct {
	cmdNames map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	method, ok := c.cmdNames[cmd.name]
	if !ok {
		return errors.New("could not find function")
	}
	err := method(s, cmd)
	if err != nil {
		return fmt.Errorf("could not run function: %w", err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdNames[name] = f
}
