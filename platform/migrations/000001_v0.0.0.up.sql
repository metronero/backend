BEGIN;

CREATE TABLE IF NOT EXISTS instance (
	version text PRIMARY KEY,
	default_commission bigint DEFAULT 1,
	custodial_mode bool DEFAULT true,
	registrations_allowed bool DEFAULT true
);

CREATE TABLE IF NOT EXISTS accounts (
	id UUID PRIMARY KEY,
	username text UNIQUE NOT NULL,
	password_hash text NOT NULL,
	creation_date timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS merchants (
	account_id UUID REFERENCES accounts ON DELETE CASCADE,
	commission bigint,
	balance bigint NOT NULL DEFAULT 0,
	wallet_address text,
	active_template_id int DEFAULT 0
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

CREATE TABLE IF NOT EXISTS merchant_stats (
	account_id UUID REFERENCES accounts ON DELETE CASCADE,
	total_comm_paid bigint NOT NULL DEFAULT 0,
	last_login timestamp
);

COMMIT;
