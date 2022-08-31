package pkg

import (
	"encoding/json"
	"os"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	Resources Resources `json:"resources" validate:"dive,keys,alphanum,endkeys,dive"`
	Root      string    `json:"root" validate:"startswith=/"`
}

func NewConfig() Config {
	return Config{
		Resources: Resources{
			"run": {
				Mode: "0755",
			},
			"shm": {
				Mode:  "0777",
				Tmpfs: true,
				Huge:  "within_size",
			},
			"tmp": {
				Mode: "0777",
			},
		},
		Root: "/var/lib/inaccel",
	}
}

func (config *Config) Read(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}

func (config Config) Validate() error {
	return validator.New().Struct(config)
}

func (config Config) Write(filename string, perm os.FileMode) error {
	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, data, perm); err != nil {
		return err
	}

	return nil
}
