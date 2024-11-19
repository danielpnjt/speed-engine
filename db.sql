CREATE TABLE public.users (
	id serial4 NOT NULL,
	username VARCHAR(255) NOT NULL,
	password VARCHAR(255) NOT NULL,
	email VARCHAR(64) NOT NULL,
	name VARCHAR(255) NOT NULL,
	balance INT NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE public.banks (
	id serial4 NOT NULL,
	user_id INT NOT NULL,
	account_name VARCHAR(255) NOT NULL,
	account_number VARCHAR(255) NOT NULL,
	bank_name VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE public.transactions (
	id serial4 NOT NULL,
	user_id INT NOT NULL,
	bank_id INT NOT NULL,
	amount INT NOT NULL,
	type TEXT NOT NULL,
	reference VARCHAR(255) NOT NULL,
	status VARCHAR(255) NOT NULL,
	expired_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL
)