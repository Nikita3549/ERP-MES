package configs_test

import (
	"strings"
	"testing"
	"time"

	"erp-mes/internal/configs"
)

type Inner struct {
	Name string `env:"TEST_CFG_NAME" default:"inner-default"`
}

type testConfig struct {
	Port    int           `env:"TEST_CFG_PORT" default:"8080"`
	Timeout time.Duration `env:"TEST_CFG_TIMEOUT" default:"5s"`
	Debug   bool          `env:"TEST_CFG_DEBUG" default:"false"`
	*Inner
}

func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{"TEST_CFG_PORT", "TEST_CFG_TIMEOUT", "TEST_CFG_DEBUG", "TEST_CFG_NAME"} {
		t.Setenv(k, "")
	}
}

func TestParseConfig(t *testing.T) {
	t.Run("default applied", func(t *testing.T) {
		clearEnv(t)

		var cfg testConfig
		if errs := configs.ParseConfig(&cfg); len(errs) != 0 {
			t.Fatalf("ParseConfig() errors = %v, want none", errs)
		}
		if cfg.Port != 8080 {
			t.Errorf("Port = %d, want 8080", cfg.Port)
		}
		if cfg.Timeout != 5*time.Second {
			t.Errorf("Timeout = %v, want 5s", cfg.Timeout)
		}
		if cfg.Debug {
			t.Errorf("Debug = true, want false")
		}
	})

	t.Run("env overrides default", func(t *testing.T) {
		clearEnv(t)
		t.Setenv("TEST_CFG_PORT", "9090")
		t.Setenv("TEST_CFG_TIMEOUT", "250ms")
		t.Setenv("TEST_CFG_DEBUG", "true")

		var cfg testConfig
		if errs := configs.ParseConfig(&cfg); len(errs) != 0 {
			t.Fatalf("ParseConfig() errors = %v, want none", errs)
		}
		if cfg.Port != 9090 {
			t.Errorf("Port = %d, want 9090", cfg.Port)
		}
		if cfg.Timeout != 250*time.Millisecond {
			t.Errorf("Timeout = %v, want 250ms", cfg.Timeout)
		}
		if !cfg.Debug {
			t.Errorf("Debug = false, want true")
		}
	})

	t.Run("invalid value returns error", func(t *testing.T) {
		clearEnv(t)
		t.Setenv("TEST_CFG_PORT", "not-a-number")

		var cfg testConfig
		errs := configs.ParseConfig(&cfg)
		if len(errs) != 1 {
			t.Fatalf("ParseConfig() errors = %v, want exactly one", errs)
		}
		if !strings.Contains(errs[0].Error(), "TEST_CFG_PORT") {
			t.Errorf("error = %q, want it to mention TEST_CFG_PORT", errs[0])
		}
		if cfg.Timeout != 5*time.Second {
			t.Errorf("Timeout = %v, want 5s (other fields must still be parsed)", cfg.Timeout)
		}
	})

	t.Run("nested pointer is populated", func(t *testing.T) {
		clearEnv(t)

		var cfg testConfig
		if errs := configs.ParseConfig(&cfg); len(errs) != 0 {
			t.Fatalf("ParseConfig() errors = %v, want none", errs)
		}
		if cfg.Inner == nil {
			t.Fatal("Inner = nil, want allocated")
		}
		if cfg.Name != "inner-default" {
			t.Errorf("Name = %q, want %q", cfg.Name, "inner-default")
		}

		t.Setenv("TEST_CFG_NAME", "from-env")
		cfg = testConfig{}
		if errs := configs.ParseConfig(&cfg); len(errs) != 0 {
			t.Fatalf("ParseConfig() errors = %v, want none", errs)
		}
		if cfg.Name != "from-env" {
			t.Errorf("Name = %q, want %q", cfg.Name, "from-env")
		}
	})
}
