module gitlab.com/metronero/backend

go 1.22.0

toolchain go1.22.6

//replace gitlab.com/metronero/metronero-go => ./metronero-go

require (
	github.com/go-chi/chi/v5 v5.1.0
	github.com/go-chi/cors v1.2.1
	github.com/go-chi/jwtauth/v5 v5.3.1
	github.com/golang-migrate/migrate/v4 v4.18.1
	github.com/google/uuid v1.6.0
	github.com/namsral/flag v1.7.4-pre
	github.com/rs/zerolog v1.33.0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	gitlab.com/moneropay/go-monero v1.1.1
	gitlab.com/moneropay/moneropay/v2 v2.5.1
	golang.org/x/crypto v0.28.0
	golang.org/x/net v0.30.0
)

require (
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/unknwon/com v1.0.1 // indirect
	golang.org/x/sync v0.8.0 // indirect
)

require (
	gitea.com/go-chi/session v0.0.0-20240316035857-16768d98ec96
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.1
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.6 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx/v2 v2.1.1 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)
