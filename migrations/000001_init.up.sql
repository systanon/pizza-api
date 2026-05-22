CREATE TABLE categories (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    slug       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE variants (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,      
    value      BIGINT NOT NULL,     
    unit       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE products (
    id          BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    image_url   TEXT NOT NULL DEFAULT '',
    metadata    JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE product_variant_prices (
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id BIGINT NOT NULL REFERENCES variants(id) ON DELETE RESTRICT,
    price      BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    PRIMARY KEY (product_id, variant_id) 
);

CREATE TABLE addons (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    price      BIGINT NOT NULL, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE category_addons (
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    addon_id    BIGINT NOT NULL REFERENCES addons(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (category_id, addon_id)
);