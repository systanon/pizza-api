CREATE TYPE order_status AS ENUM ('pending', 'paid', 'cancelled');

CREATE TABLE orders (
    id         BIGSERIAL PRIMARY KEY,
    cart_id    UUID        NOT NULL REFERENCES carts(id) ON DELETE RESTRICT,
    email      TEXT        NOT NULL,
    status     order_status NOT NULL DEFAULT 'pending',
    total      BIGINT      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Snapshot of cart items at the time of order (prices frozen)
CREATE TABLE order_items (
    id            BIGSERIAL PRIMARY KEY,
    order_id      BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id    BIGINT NOT NULL,
    product_name  TEXT   NOT NULL,
    variant_id    BIGINT,
    variant_name  TEXT,
    variant_price BIGINT NOT NULL,
    quantity      INT    NOT NULL CHECK (quantity > 0),
    item_total    BIGINT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE order_item_addons (
    order_item_id BIGINT NOT NULL REFERENCES order_items(id) ON DELETE CASCADE,
    addon_id      BIGINT NOT NULL,
    addon_name    TEXT   NOT NULL,
    price         BIGINT NOT NULL,
    PRIMARY KEY (order_item_id, addon_id)
);

CREATE TRIGGER trg_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
