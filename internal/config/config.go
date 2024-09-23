package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Env        string   `yaml:"env"`
	Server     Server   `yaml:"server"`
	Database   Database `yaml:"database"`
	configPath string   `env:"CONFIG_PATH"`
}

type Server struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

func LoadConfig() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file %v", err)
		}
	}

	configPath := os.Getenv("CONFIG_PATH")
	env := os.Getenv("MY_ENV")
	fmt.Println(configPath)
	if configPath == "" {
		log.Fatalf("CONFIG_PATH environment variable not set")
	}
	switch env {
	case "local":
		configPath = filepath.Join(configPath, "config.local.yaml")
	case "test":
		configPath = filepath.Join(configPath, "config.test.yaml")
	default:
		return nil, errors.New("environment not supported")
	}

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
