CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    firstname text NOT NULL,
    lastname text NOT NULL,
    email citext UNIQUE NOT NULL,
    phone text UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    token text NOT NULL,
    refresh_token text NOT NULL,
    version integer NOT NULL DEFAULT 1
);