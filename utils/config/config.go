package config

import (
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/namsral/flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	BindAddr          string
	PostgresUri       string
	JwtSecret         *jwtauth.JWTAuth
	CallbackUrl       string
	MoneroPay         string
)

func Load() {
	flag.StringVar(&BindAddr, "bind", "localhost:5001", "Bind address:port")
	flag.StringVar(&PostgresUri, "postgres", "postgresql://metronero:m3tr0n3r0@localhost:5432/metronero?sslmode=disable", "PostgreSQL connection string")
	var jwtSecretStr string
	flag.StringVar(&jwtSecretStr, "token-secret", "aabbccddeeffgg", "JWT secret")
	flag.StringVar(&CallbackUrl, "callback-url", "http://localhost:8080/callback", "Incoming callback URL")
	flag.StringVar(&MoneroPay, "moneropay", "http://localhost:5000", "MoneroPay instance")
	var logFormat string
	flag.StringVar(&logFormat, "log-format", "pretty", "Log format (pretty or json)")
	flag.Parse()

	JwtSecret = jwtauth.New("HS256", []byte(jwtSecretStr), nil)
	if logFormat == "pretty" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr,
		    TimeFormat: time.RFC3339})
	}
}
