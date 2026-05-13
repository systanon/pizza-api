CREATE TABLE categories (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    slug       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE products (
    id          BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price       BIGINT NOT NULL,
    image_url   TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);