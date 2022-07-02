package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

const applicationConfigFile = "application.yaml"

type ApplicationConfig struct {
	Router   RouterConfig `mapstructure:"server"`
	Database DatabaseConfig
}

func main() {
	fail := func(err error) {
		log.Fatalf("Failed to start application: %v", err)
	}
	log.Printf("Starting application..")

	// Read application config from file.
	cfg, err := readApplicationConfig(applicationConfigFile)
	if err != nil {
		fail(err)
	}

	// Connect to database.
	db, err := newDatabase(cfg.Database)
	if err != nil {
		fail(err)
	}

	// Start HTTP server.
	address := fmt.Sprintf("%s:%d", cfg.Router.Host, cfg.Router.Port)
	log.Printf("Starting HTTP server on %q\n", address)
	if err := http.ListenAndServe(address, newRouter(StatusHandler{}, BusinessHandler{BusinessStore{db}})); err != nil {
		fail(err)
	}
	log.Printf("Application stopped\n")
}

func readApplicationConfig(file string) (ApplicationConfig, error) {
	fail := func(err error) (ApplicationConfig, error) {
		return ApplicationConfig{}, fmt.Errorf("readApplicationConfig: failed to read application config from file %q: %w", file, err)
	}

	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return fail(fmt.Errorf("failed to load application config: %w", err))
	}

	var cfg ApplicationConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return fail(fmt.Errorf("failed to deserialise application config: %w", err))
	}

	return cfg, nil
}
