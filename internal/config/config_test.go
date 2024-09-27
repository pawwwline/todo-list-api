package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
	cfg "todo-list-api/internal/config"
)

func CreateTmpFile(content string, env string) (string, string, error) {
	tempDir := os.TempDir()
	tempPath := filepath.Join(tempDir, fmt.Sprintf("config.%s.yaml", env))
	err := os.WriteFile(tempPath, []byte(content), 0644)
	if err != nil {
		return "", "", err
	}
	return tempPath, tempDir, nil

}
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		env            string
		content        string
		expectedConfig cfg.ConfigYaml
		expectError    bool
	}{{name: "Local case",
		env: "local",
		content: `
env: "local"
server:
  host: "localhost"
  port: "8080"
  timeout: "5s"
  idle_timeout: "60s"
database:
  username: "myuser"
  password: "pass"
  host: "localhost"
  port: "5432"
  name: "todolist"
  ssl_mode: "disable"
  `,
		expectedConfig: cfg.ConfigYaml{
			Env: "local",
			Server: cfg.Server{
				Host:        "localhost",
				Port:        "8080",
				Timeout:     5 * time.Second,
				IdleTimeout: 60 * time.Second,
			},
			Database: cfg.Database{
				Username: "myuser",
				Password: "pass",
				Host:     "localhost",
				Port:     "5432",
				Name:     "todolist",
				SSLMode:  "disable",
			},
		},
		expectError: false,
	}, {
		name: "Test case",
		env:  "test",
		content: `
env: "test"
server:
  host: "localhost"
  port: "8080"
  timeout: "5s"
  idle_timeout: "60s"
database:
  username: "myuser"
  password: "pass"
  host: "localhost"
  port: "5432"
  name: "test"
  ssl_mode: "disable"
`,
		expectedConfig: cfg.ConfigYaml{
			Env: "test",
			Server: cfg.Server{
				Host:        "localhost",
				Port:        "8080",
				Timeout:     5 * time.Second,
				IdleTimeout: 60 * time.Second,
			},
			Database: cfg.Database{
				Username: "myuser",
				Password: "pass",
				Host:     "localhost",
				Port:     "5432",
				Name:     "test",
				SSLMode:  "disable",
			},
		},
		expectError: false},
		{name: "Error case",
			expectError: true,
			env:         "error",
			content: `
env: "local"
server:
  host: "localhost"
  port: "8080"
  timeout: "5s"
  idle_timeout: "60s"
database:
  username: "myuser"
  password: "pass"
  host: "localhost"
  port: "5432"
  name: "todolist"
  ssl_mode: "disable"
  `, expectedConfig: cfg.ConfigYaml{
				Env: "local",
				Server: cfg.Server{
					Host:        "localhost",
					Port:        "8080",
					Timeout:     5 * time.Second,
					IdleTimeout: 60 * time.Second,
				},
				Database: cfg.Database{
					Username: "myuser",
					Password: "pass",
					Host:     "localhost",
					Port:     "5432",
					Name:     "todolist",
					SSLMode:  "disable",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, configPath, err := CreateTmpFile(test.content, test.env)
			if err != nil {
				t.Errorf("got error in creating file: %v", err)
			}
			defer os.Remove(path)
			os.Setenv("CONFIG_PATH", filepath.Dir(configPath))
			os.Setenv("MY_ENV", test.env)
			defer os.Unsetenv("CONFIG_PATH")
			defer os.Unsetenv("MY_ENV")
			config, err := cfg.LoadConfig()
			if test.expectError {
				if err == nil {
					t.Errorf("expected error got nil")
				}
				return
			}

			if !test.expectError && err != nil {
				t.Errorf("expected no error got %v", err)
			}
			if config == nil {
				t.Fatalf("configuration is nil")
			}
			if *&config.ConfigYaml != test.expectedConfig {
				t.Errorf("expected %v got %v", test.expectedConfig, *config)
			}

		})
	}
}
