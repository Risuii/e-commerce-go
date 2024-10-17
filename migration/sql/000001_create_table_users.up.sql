CREATE TABLE users (
	id bigserial primary key,
	email varchar null unique,
	email_verified bool not null default false,
	phone varchar null unique,
	phone_verified bool not null default false,
	password varchar null,
	username varchar null
	created_at timestamptz not null default CURRENT_TIMESTAMP,
	updated_at timestamptz not null default CURRENT_TIMESTAMP,
	deleted_at timestamptz null,
);