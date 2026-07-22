package main

import (
	"context"
	"database/sql"
	"fmt"
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
	inCommands.register("users", handlerListUsers)
	inCommands.register("agg", handlerAgg)
	inCommands.register("feeds", handlerListFeeds)
	inCommands.register("reset", handlerReset)
	inCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	inCommands.register("follow", middlewareLoggedIn(handlerFollow))
	inCommands.register("following", middlewareLoggedIn(handlerFollowing))
	inCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	inCommands.register("browse", middlewareLoggedIn(handlerBrowse))

	commandArgs := os.Args
	if len(commandArgs) < 2 {
		log.Fatal("command not found")
	}

	switch commandArgs[1] {
	case "register", "login", "users", "agg", "addfeed", "feeds", "follow", "following", "unfollow", "browse", "reset":
		cmdName, cmdArg := commandArgs[1], commandArgs[2:]
		err = inCommands.run(inState, command{
			name: cmdName,
			args: cmdArg,
		})
		if err != nil {
			log.Fatalf("could not run command: %v", err)
		}
	default:
		log.Fatalf("could not run command, bad input")
	}
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
