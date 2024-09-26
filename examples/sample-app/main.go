package main

import (
	"encoding/json"
	"fmt"
	"github.com/marnixbouhuis/confpb/pkg/config"
	configv1 "sample-app/gen/config/v1"
)

func main() {
	// Take default values as configuration base
	conf := configv1.ApplicationConfigFromDefault()

	// Get configuration variables set in environment
	envConfig, err := configv1.ApplicationConfigFromEnv()
	if err != nil {
		panic(err)
	}

	// Overlay the config from the environment variables onto the default config
	conf = config.Overlay(conf, envConfig)

	// Load config from file
	fileConfig, err := config.FromFile[configv1.ApplicationConfig]("./config.yaml")
	if err != nil {
		panic(err)
	}

	// Overlay the config again over the values from env / defaults
	conf = config.Overlay(conf, fileConfig)

	// The resulting config uses the following priority:
	// - Config file first
	// - ENV variable
	// - Default value
	// By reordering the merges different behaviour can be achieved,

	// Do something with the resulting config
	prettyPrint(conf)
}

func prettyPrint[T any](t T) {
	s, _ := json.MarshalIndent(t, "", "\t")
	fmt.Println(string(s))
}
