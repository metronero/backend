BEGIN;				   
CREATE TABLE IF NOT EXISTS accounts (
	username text PRIMARY KEY,
	password text NOT NULL,
	commission smallint,
	wallet_address character(95) NOT NULL,
	template text
);
CREATE TABLE IF NOT EXISTS payments (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	amount bigint NOT NULL,
	accept_url text NOT NULL,
	cancel_url text NOT NULL,
	callback_url text NOT NULL DEFAULT ' ',
	username text REFERENCES accounts ON DELETE CASCADE,
	address text NOT NULL,
	callback_data text NOT NULL DEFAULT ' ',
	received bool NOT NULL DEFAULT false,
	completed bool NOT NULL DEFAULT false,
	withdrawn bool NOT NULL DEFAULT false
);
COMMIT;
