// Package configs provides a struct for managing .env variables
package configs

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	*HTTPConfig
	*DBConfig
}

type HTTPConfig struct {
	Port              int           `env:"HTTP_PORT"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT"`
	ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT"`
	WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT"`
	ShutdownTimeout   time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT"`
}

type DBConfig struct {
	DSN string `env:"DSN"`
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("godotenv error: %s", err.Error())
	}

	conf := &Config{}
	errors := ParseConfig(conf)

	if len(errors) > 0 {
		printErrors(errors)
	}

	return conf
}

func printErrors(errs []error) {
	var msg strings.Builder
	msg.WriteString("config errors:\n")
	for _, err := range errs {
		fmt.Fprintf(&msg, "  - %s\n", err.Error())
	}
	log.Fatal(msg.String())
}
