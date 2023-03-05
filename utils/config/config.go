package config

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/namsral/flag"
)

var (
	BindAddr string
	PostgresUri string
	JwtSecret *jwtauth.JWTAuth
	CommissionAddress string
	CallbackUrl string
	MoneroPay string
)

func Load() {
	flag.StringVar(&BindAddr, "bind", "localhost:8080", "Bind address:port")
	flag.StringVar(&PostgresUri, "postgres", "postgresql://metronero:m3tr0n3r0@localhost:5432/metronero", "PostgreSQL connection string")
	var jwtSecretStr string
	flag.StringVar(&jwtSecretStr, "jwt-secret", "aabbccddeeffgg", "JWT secret")
	flag.StringVar(&CommissionAddress, "commission-address", "46VGoe3bKWTNuJdwNjjr6oGHLVtV1c9QpXFP9M2P22bbZNU7aGmtuLe6PEDRAeoc3L7pSjfRHMmqpSF5M59eWemEQ2kwYuw", "Monero address for commissions")
	flag.StringVar(&CallbackUrl, "callback-url", "http://localhost:8080/callback", "Incoming callback URL")
	flag.StringVar(&MoneroPay, "moneropay", "http://localhost:5000", "MoneroPay instance")
	flag.Parse()
	JwtSecret = jwtauth.New("HS256", []byte(jwtSecretStr), nil)
}
