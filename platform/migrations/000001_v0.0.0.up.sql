BEGIN;

CREATE TABLE IF NOT EXISTS instance (
	version text PRIMARY KEY,
	default_commission bigint DEFAULT 1,
	custodial_mode bool DEFAULT true,
	registrations_allowed bool DEFAULT true,
	withdrawal_times text DEFAULT 'instant'
);

CREATE TABLE IF NOT EXISTS instance_bootstrap (
	first_run bool DEFAULT false NOT NULL
);

CREATE TABLE IF NOT EXISTS instance_stats (
	wallet_balance bigint DEFAULT 0 NOT NULL,
	total_profits bigint DEFAULT 0 NOT NULL,
	total_merchants bigint DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
	account_id UUID PRIMARY KEY,
	username text UNIQUE NOT NULL,
	password_hash text NOT NULL
);

CREATE TABLE IF NOT EXISTS account_stats (
	account_id UUID PRIMARY KEY REFERENCES accounts ON DELETE CASCADE,
	creation_date timestamp NOT NULL,
	last_login timestamp
);

CREATE TABLE IF NOT EXISTS merchants (
	account_id UUID PRIMARY KEY REFERENCES accounts ON DELETE CASCADE,
	commission bigint,
	wallet_address text,
	active_template_id int DEFAULT 0
);

CREATE TABLE IF NOT EXISTS merchant_stats (
	account_id UUID PRIMARY KEY REFERENCES accounts ON DELETE CASCADE,
	balance bigint NOT NULL DEFAULT 0,
	total_sales bigint NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS api_tokens (
	account_id UUID PRIMARY KEY REFERENCES accounts ON DELETE CASCADE,
	token text NOT NULL,
	valid_until timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS callback_secrets (
	account_id UUID PRIMARY KEY REFERENCES accounts ON DELETE CASCADE,
	secret_key text NOT NULL,
	valid_until timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS templates (
	id SERIAL PRIMARY KEY,
	account_id UUID REFERENCES accounts ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payments (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	amount bigint NOT NULL,
	accept_url text NOT NULL,
	cancel_url text NOT NULL,
	callback_url text NOT NULL DEFAULT ' ',
	account_id UUID REFERENCES accounts ON DELETE CASCADE,
	address text NOT NULL,
	callback_data text NOT NULL DEFAULT ' ',
	status text NOT NULL DEFAULT 'pending'
);

COMMIT;
