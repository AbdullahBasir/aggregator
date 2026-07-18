package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const (
	configFileName = ".gatorconfig.json"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %v", err)
	}
	filepath := path.Join(homeDir, configFileName)

	return filepath, nil
}

func Read() (Config, error) {
	file, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("could not get filepath: %v", err)
	}
	content, err := os.ReadFile(file)
	if err != nil {
		return Config{}, fmt.Errorf("could not read file: %v", err)
	}

	cfg := Config{}
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("could not decode file content: %v", err)
	}
	return cfg, nil
}

func write(cfg Config) error {
	file, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("could not get filepath: %v", err)
	}

	encode, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("could not encode data to json: %v", err)
	}

	err = os.WriteFile(file, encode, 0640)
	if err != nil {
		return fmt.Errorf("could not write data to file: %v", err)
	}
	return nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	err := write(*cfg)

	if err != nil {
		return fmt.Errorf("could not write username to file: %v", err)
	}
	return nil
}
