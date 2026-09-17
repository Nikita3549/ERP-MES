// Package configs provides a struct for managing env variables
package configs

import (
	"errors"
	"fmt"
	"io/fs"
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
	Port              int           `env:"HTTP_PORT" default:"8080"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" default:"5s"`
	ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT" default:"15s"`
	WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT" default:"20s"`
	IdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT" default:"60s"`
	ShutdownTimeout   time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" default:"20s"`
}

type DBConfig struct {
	Host     string `env:"DB_HOST" default:"localhost"`
	Port     int    `env:"DB_PORT" default:"5432"`
	Name     string `env:"DB_NAME" default:"erp-mes"`
	Password string `env:"DB_PASSWORD" default:"postgres"`
	User     string `env:"DB_USER" default:"postgres"`
}

func (d *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", d.Host, d.User, d.Password, d.Name, d.Port)
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("load config: %v", err)
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
