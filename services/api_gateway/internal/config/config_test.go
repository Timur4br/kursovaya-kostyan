package config

import (
	"net/url"
	"testing"
)

func environment() map[string]string {
	return map[string]string{"DB_HOST": "db", "DB_PORT": "5432", "DB_NAME": "coursework", "DB_USER": "app", "DB_PASSWORD": "p@ss:/?#%word", "DB_SSLMODE": "disable"}
}
func getter(values map[string]string) func(string) string {
	return func(k string) string { return values[k] }
}

func TestPasswordIsEscaped(t *testing.T) {
	values := environment()
	cfg, err := FromEnv(getter(values))
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	password, _ := u.User.Password()
	if password != values["DB_PASSWORD"] {
		t.Fatal("password changed in connection URL")
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatal(cfg.HTTPAddr)
	}
}

func TestConfigurationRejectsMissingAndInvalidValues(t *testing.T) {
	for _, key := range []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD"} {
		t.Run(key, func(t *testing.T) {
			values := environment()
			delete(values, key)
			if _, err := FromEnv(getter(values)); err == nil {
				t.Fatal("missing value accepted")
			}
		})
	}
	for _, port := range []string{"0", "65536", "abc"} {
		values := environment()
		values["DB_PORT"] = port
		if _, err := FromEnv(getter(values)); err == nil {
			t.Fatal("invalid database port accepted")
		}
	}
	values := environment()
	values["DB_SSLMODE"] = "unexpected"
	if _, err := FromEnv(getter(values)); err == nil {
		t.Fatal("invalid TLS mode accepted")
	}
	values = environment()
	values["HTTP_ADDR"] = ":0"
	if _, err := FromEnv(getter(values)); err == nil {
		t.Fatal("invalid HTTP port accepted")
	}
}
