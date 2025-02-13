package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"quizapi/common"
	"quizapi/modules/usersmodule"
	"strings"
)

// config definition
// config getter

// config definition

type Env string

const (
	DeployEnv Env = "deploy"
	DevEnv    Env = "dev"
)

var Envs = []string{string(DeployEnv), string(DevEnv)}

func (e *Env) Valid() []error {
	var errs []error
	var found bool = false
	for _, env := range Envs {
		if env == string(*e) {
			found = true
			break
		}
	}

	if !found {
		errs = append(errs, fmt.Errorf("env has to be one of these: %s", strings.Join(Envs, "")))
	}

	return errs
}

type MainConfig struct {
	Env         Env                    `json:"env"`
	Port        int                    `json:"port"`
	JwtSecret   string                 `json:"jwt_secret"`
	UsersConfig usersmodule.UserConfig `json:"user_config"`
}

func (c *MainConfig) Valid() []error {
	var errs []error
	for _, err := range c.Env.Valid() {
		errs = append(errs, common.ErrPath(err).Property("env"))
	}
	if c.Port <= 0 {
		errs = append(errs, common.ErrPath(errors.New("port has to be positive number")).Property("port"))
	}
	if c.JwtSecret == "" {
		errs = append(errs, common.ErrPath(errors.New("missing jwt secret")).Property("jwt_secret"))
	}
	for _, err := range c.UsersConfig.Valid() {
		errs = append(errs, common.ErrPath(err).Property("user_config"))
	}
	return errs
}

// config getter

func GetValidConfig() MainConfig {
	configPath := flag.String("config", "./env.json", "Path to the configuration file")
	flag.Parse()
	configData := *configPath

	file, err := os.Open(configData)
	if err != nil {
		log.Panic(err.Error())
	}
	defer file.Close()

	var config MainConfig
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
