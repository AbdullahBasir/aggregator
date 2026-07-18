package main

import (
	"log"
	"os"

	"github.com/AbdullahBasir/aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("could not read file: %v", err)
	}

	err = cfg.SetUser("Abdullah")
	if err != nil {
		log.Fatalf("could not set username: %v", err)
	}

	inState := &state{}
	inState.config = &cfg

	cmdMap := map[string]func(*state, command) error{}
	inCommands := &commands{
		cmdNames: cmdMap,
	}
	inCommands.register("login", handlerLogin)
	commandArgs := os.Args
	if len(commandArgs) < 2 {
		log.Fatal("command not found")
	}
	cmdName, cmdArg := commandArgs[1], commandArgs[2:]
	err = inCommands.run(inState, command{
		name: cmdName,
		args: cmdArg,
	})
	if err != nil {
		log.Fatalf("could not run command: %v", err)
	}
}
