CREATE TABLE IF NOT EXISTS products (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    price bigint NOT NULL,
    rating int NOT NULL,
    image text NOT NULL
);