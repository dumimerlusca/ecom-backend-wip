CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE product_status AS ENUM ('draft', 'published', 'deleted');

CREATE TABLE IF NOT EXISTS product (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    subtitle text,
    description text NOT NULL,
    thumbnail_id text,
    status product_status NOT NULL DEFAULT 'draft',
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp
);

CREATE TABLE IF NOT EXISTS product_variant (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    product_id uuid NOT NULL,
    title text NOT NULL,
    sku text,
    barcode bigint,
    material text,
    weight float,
    length float,
    width float,
    height float,
    inventory_quantity int NOT NULL DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp,
    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_product_variant_product_id ON product_variant(product_id);

CREATE TABLE IF NOT EXISTS product_option (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    product_id uuid NOT NULL,
    title text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp,
    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE,
    CONSTRAINT duplicate_option_not_allowed UNIQUE (product_id, title)
);

CREATE INDEX IF NOT EXISTS idx_product_option_product_id ON product_option(product_id);

CREATE TABLE IF NOT EXISTS product_option_value (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    option_id uuid NOT NULL,
    variant_id uuid NOT NULL,
    title text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp,
    FOREIGN KEY (option_id) REFERENCES product_option(id) ON DELETE CASCADE,
    FOREIGN KEY (variant_id) REFERENCES product_variant(id) ON DELETE CASCADE,
    CONSTRAINT duplicate_option_value_not_allowed UNIQUE (option_id, title, variant_id)
);

CREATE INDEX IF NOT EXISTS idx_product_option_variant_id ON product_option_value(variant_id);

CREATE TABLE IF NOT EXISTS currency (
    code VARCHAR(10) PRIMARY KEY NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    symbol_native VARCHAR(20) NOT NULL,
    name text NOT NULL
);

CREATE TABLE IF NOT EXISTS money_amount (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    currency_code VARCHAR(10) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp,
    FOREIGN KEY (currency_code) REFERENCES currency(code) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_variant_money_amount (
    variant_id uuid NOT NULL,
    money_amount_id uuid NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp,
    PRIMARY KEY (variant_id, money_amount_id),
    FOREIGN KEY (variant_id) REFERENCES product_variant(id) ON DELETE CASCADE,
    FOREIGN KEY (money_amount_id) REFERENCES money_amount(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_category (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    parent_id uuid,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp
);

CREATE TABLE IF NOT EXISTS product_category_product (
    category_id uuid NOT NULL,
    product_id uuid NOT NULL,
    PRIMARY KEY (category_id, product_id),
    FOREIGN KEY (category_id) REFERENCES product_category(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS file (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    original_name text NOT NULL,
    mime_type text NOT NULL,
    extension text NOT NULL,
    size int NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS entity_file (
    file_id uuid NOT NULL,
    entity_id text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (file_id, entity_id),
    FOREIGN KEY (file_id) REFERENCES file(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    is_admin bool NOT NULL DEFAULT false,
    activated bool NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS token (
    hash bytea PRIMARY KEY NOT NULL,
    user_id uuid NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp with time zone NOT NULL,
    scope text NOT NULL
);