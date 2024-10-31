package config

import (
	"html/template"
	"net/url"
	"os"
	"time"

	"github.com/namsral/flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const Version = "0.0.0"

var (
	BindAddr               string
	CallbackAddr           string
	PostgresUri            string
	CallbackUrl            string
	MoneroPay              string
	TemplateDir            string
	TemplateMaxSize        int
	CorsOrigin             string
	DefaultPaymentTemplate *template.Template
)

func Load() {
	flag.StringVar(&BindAddr, "bind", "localhost:5001", "Bind address:port")

	// CallbackAddr is passed to MoneroPay on payment request creation.
	// This value should be accessible to the MoneroPay instance.
	flag.StringVar(&CallbackAddr, "callback-addr", "http://localhost:5001", "http(s)://domain:port for MoneroPay callback registration")
	flag.StringVar(&PostgresUri, "postgres", "postgresql://metronero:m3tr0n3r0@localhost:5432/metronero?sslmode=disable", "PostgreSQL connection string")
	flag.StringVar(&CallbackUrl, "callback-url", "http://localhost:8080/callback", "Incoming callback URL")
	flag.StringVar(&MoneroPay, "moneropay", "http://localhost:5000", "MoneroPay instance")
	var logFormat string
	flag.StringVar(&logFormat, "log-format", "pretty", "Log format (pretty or json)")
	flag.StringVar(&TemplateDir, "template-dir", "./data/merchant_templates", "Directory to save merchant templates.")
	flag.StringVar(&CorsOrigin, "cors-origin", "http://*", "Allowed CORS origin")
	flag.IntVar(&TemplateMaxSize, "template-max-size", 20, "Maximum template upload size in MiB")
	flag.Parse()

	if logFormat == "pretty" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr,
			TimeFormat: time.RFC3339})
	}

	templatePath, err := url.JoinPath(TemplateDir, "default.html")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load default payment template")
	}
	if DefaultPaymentTemplate, err = template.ParseFiles(templatePath); err != nil {
		log.Fatal().Err(err).Msg("Failed to load default payment template")
	}
}
