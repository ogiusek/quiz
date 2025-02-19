package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"strings"
)

type Config struct {
	RabbitMqUrl string `json:"rabbitmq_url"`
	Port        int    `json:"port"`
}

func (config *Config) Valid() []error {
	var errs []error
	if config.RabbitMqUrl == "" {
		errs = append(errs, errors.New("`url` is missing"))
	}
	if config.Port <= 0 {
		errs = append(errs, errors.New("`port` has to be positive"))
	}
	return errs
}

func GetConfig() Config {
	configPath := flag.String("config", "./env.json", "Path to the configuration file")
	flag.Parse()
	configData := *configPath

	file, err := os.Open(configData)
	if err != nil {
		log.Panic(err.Error())
	}
	defer file.Close()
	var config Config
	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(&config); err != nil {
		log.Panic(err.Error())
	}

	if errors := config.Valid(); len(errors) != 0 {
		var messages []string
		for _, err := range errors {
			messages = append(messages, err.Error())
		}
		log.Panicf("errors:\n```\n%s\n```", strings.Join(messages, "\n"))
	}

	return config
}
