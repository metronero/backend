BEGIN;

CREATE TABLE IF NOT EXISTS instance (
	version text PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS accounts (
	account_id UUID PRIMARY KEY,
	username text UNIQUE NOT NULL,
	password text NOT NULL,
	role text NOT NULL DEFAULT 'merchant',
	creation_date timestamp NOT NULL,
	last_login timestamp
);

CREATE TABLE IF NOT EXISTS api_keys (
    key_id UUID PRIMARY KEY,
    key_secret text NOT NULL,
    expiry timestamp NOT NULL,
    account_id UUID NOT NULL REFERENCES accounts ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS merchants (
	account_id UUID PRIMARY KEY REFERENCES accounts ON DELETE CASCADE,
	total_sales bigint NOT NULL DEFAULT 0,
	disabled bool DEFAULT false NOT NULL
);

CREATE TABLE IF NOT EXISTS callback_secrets (
	account_id UUID PRIMARY KEY REFERENCES accounts ON DELETE CASCADE,
	secret_key text NOT NULL,
	valid_until timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS payments (
	payment_id UUID PRIMARY KEY,
	amount bigint NOT NULL,
	fee bigint NOT NULL,
	order_id text DEFAULT '',
	accept_url text DEFAULT '',
	cancel_url text DEFAULT '',
	callback_url text DEFAULT '',
	address text NOT NULL,
	callback_data text,
	status text NOT NULL DEFAULT 'Pending',
	last_update timestamp NOT NULL,
	account_id UUID REFERENCES accounts ON DELETE CASCADE,
	merchant_name text,
	merchant_extra text DEFAULT ''
);

CREATE TABLE IF NOT EXISTS withdrawals (
	withdrawal_id UUID PRIMARY KEY,
	amount bigint NOT NULL,
	withdraw_date timestamp NOT NULL
);

COMMIT;
