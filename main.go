package main

import (
	"fmt"
	"log"

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

	modifyCfg, err := config.Read()
	if err != nil {
		log.Fatalf("could not read file: %v", err)
	}

	fmt.Printf("Username: %s, Url: %s", modifyCfg.CurrentUserName, modifyCfg.DbUrl)
}
