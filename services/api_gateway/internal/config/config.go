package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

// Load reads process environment; Compose injects values from the root .env.
func Load() (Config, error) { return FromEnv(os.Getenv) }

func FromEnv(get func(string) string) (Config, error) {
	for _, key := range []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD"} {
		if strings.TrimSpace(get(key)) == "" {
			return Config{}, fmt.Errorf("%s is required", key)
		}
	}
	port, err := strconv.Atoi(get("DB_PORT"))
	if err != nil || port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("DB_PORT must be between 1 and 65535")
	}
	mode := get("DB_SSLMODE")
	if mode == "" {
		mode = "require"
	}
	switch mode {
	case "disable", "require", "verify-ca", "verify-full":
	default:
		return Config{}, fmt.Errorf("unsupported DB_SSLMODE")
	}
	u := &url.URL{Scheme: "postgres", User: url.UserPassword(get("DB_USER"), get("DB_PASSWORD")), Host: net.JoinHostPort(get("DB_HOST"), strconv.Itoa(port)), Path: "/" + get("DB_NAME")}
	q := u.Query()
	q.Set("sslmode", mode)
	q.Set("connect_timeout", "5")
	u.RawQuery = q.Encode()
	addr := get("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	_, httpPort, err := net.SplitHostPort(addr)
	if err != nil {
		return Config{}, fmt.Errorf("HTTP_ADDR must be host:port")
	}
	n, err := strconv.Atoi(httpPort)
	if err != nil || n < 1 || n > 65535 {
		return Config{}, fmt.Errorf("invalid HTTP_ADDR port")
	}
	return Config{HTTPAddr: addr, DatabaseURL: u.String()}, nil
}
