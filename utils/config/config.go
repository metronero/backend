package daemon

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/namsral/flag"
)

type config struct {
	BindAddr string
	postgresCS string
	JwtSecret *jwtauth.JWTAuth
	CommissionAddress string
	callbackURL string
	moneropay string
}

var Config config

func loadConfig() {
	flag.StringVar(&Config.BindAddr, "bind", "localhost:8080", "Bind address:port")
	flag.StringVar(&Config.postgresCS, "postgresql", "postgresql://metronero:m3tr0n3r0@localhost:5432/metronero", "PostgreSQL connection string")
	var jwtSecret string
	flag.StringVar(&jwtSecret, "jwt-secret", "aabbccddeeffgg", "JWT secret")
	flag.StringVar(&Config.CommissionAddress, "commission-address", "46VGoe3bKWTNuJdwNjjr6oGHLVtV1c9QpXFP9M2P22bbZNU7aGmtuLe6PEDRAeoc3L7pSjfRHMmqpSF5M59eWemEQ2kwYuw", "Monero address for commissions")
	flag.StringVar(&Config.callbackURL, "callback-url", "http://localhost:8080/callback", "Incoming callback URL")
	flag.StringVar(&Config.moneropay, "moneropay", "http://localhost:5000", "MoneroPay instance")
	flag.Parse()
	Config.JwtSecret = jwtauth.New("HS256", []byte(jwtSecret), nil)
}
