CREATE SCHEMA appleshop;

CREATE TABLE appleshop.users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    
    role VARCHAR(50) NOT NULL DEFAULT 'user'
        CHECK (role IN ('user', 'admin')),
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE appleshop.categories (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE appleshop.products (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(255) NOT NULL,
    description TEXT,

    price NUMERIC(10, 2) NOT NULL
        CHECK (price >= 0),

    stock_quantity INTEGER NOT NULL DEFAULT 0
        CHECK (stock_quantity >= 0),

    category_id BIGINT NOT NULL,

    image_url TEXT,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_products_category
        FOREIGN KEY (category_id)
        REFERENCES appleshop.categories(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_products_category_id
    ON appleshop.products(category_id);

CREATE TABLE appleshop.baskets (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_baskets_user
        FOREIGN KEY (user_id)
        REFERENCES appleshop.users(id)
        ON DELETE CASCADE
);

CREATE TABLE appleshop.basket_items (
    id BIGSERIAL PRIMARY KEY,

    basket_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,

    quantity INTEGER NOT NULL
        CHECK (quantity > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_basket_items_basket
        FOREIGN KEY (basket_id)
        REFERENCES appleshop.baskets(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_basket_items_product
        FOREIGN KEY (product_id)
        REFERENCES appleshop.products(id)
        ON DELETE RESTRICT,

    CONSTRAINT uq_basket_product
        UNIQUE (basket_id, product_id)
);

CREATE INDEX idx_basket_items_basket_id
    ON appleshop.basket_items(basket_id);

CREATE INDEX idx_basket_items_product_id
    ON appleshop.basket_items(product_id);

CREATE TABLE appleshop.addresses (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,

    country VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    street VARCHAR(255) NOT NULL,
    house VARCHAR(50) NOT NULL,
    apartment VARCHAR(50),
    postal_code VARCHAR(20),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_addresses_user
        FOREIGN KEY (user_id)
        REFERENCES appleshop.users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_addresses_user_id
    ON appleshop.addresses(user_id);

CREATE TABLE appleshop.orders (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    address_id BIGINT NOT NULL,

    status VARCHAR(50) NOT NULL DEFAULT 'pending'
        CHECK (
            status IN (
                'pending',
                'paid',
                'processing',
                'shipped',
                'delivered',
                'cancelled'
            )
        ),

    total_price NUMERIC(10, 2) NOT NULL
        CHECK (total_price >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_orders_user
        FOREIGN KEY (user_id)
        REFERENCES appleshop.users(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_orders_address
        FOREIGN KEY (address_id)
        REFERENCES appleshop.addresses(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_orders_user_id
    ON appleshop.orders(user_id);

CREATE INDEX idx_orders_address_id
    ON appleshop.orders(address_id);

CREATE INDEX idx_orders_status
    ON appleshop.orders(status);

CREATE TABLE appleshop.order_items (
    id BIGSERIAL PRIMARY KEY,

    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,

    quantity INTEGER NOT NULL
        CHECK (quantity > 0),

    unit_price NUMERIC(10, 2) NOT NULL
        CHECK (unit_price >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id)
        REFERENCES appleshop.orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_order_items_product
        FOREIGN KEY (product_id)
        REFERENCES appleshop.products(id)
        ON DELETE RESTRICT,

    CONSTRAINT uq_order_product
        UNIQUE (order_id, product_id)
);

CREATE INDEX idx_order_items_order_id
    ON appleshop.order_items(order_id);

CREATE INDEX idx_order_items_product_id
    ON appleshop.order_items(product_id);