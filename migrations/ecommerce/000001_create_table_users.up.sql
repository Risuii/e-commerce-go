CREATE TABLE users (
	id bigserial primary key,
	uuid varchar(50) not null,
	email varchar null unique,
	email_verified bool not null default false,
	password text not null,
	username varchar(50) not null,
	last_login varchar(24) null,
	created_at varchar(50) NULL,
	updated_at varchar(50) NULL,
	deleted_at varchar(50) NULL
);

