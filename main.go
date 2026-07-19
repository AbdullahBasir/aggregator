package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/AbdullahBasir/aggregator/internal/config"
	"github.com/AbdullahBasir/aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("could not read file: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}

	dbQueries := database.New(db)

	inState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmdMap := map[string]func(*state, command) error{}
	inCommands := &commands{
		cmdNames: cmdMap,
	}
	inCommands.register("login", handlerLogin)
	inCommands.register("register", handlerRegister)
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
