CREATE EXTENSION IF NOT EXISTS "pgcrypto"; -- Required for gen_random_uuid()

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "role" varchar NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "products" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "category" varchar NOT NULL,
  "product_name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "brand" varchar,
  "count_in_stock" bigint NOT NULL,
  "price" float NOT NULL,
  "currency" varchar NOT NULL,
  "rating" bigint,
  "is_featured" bool DEFAULT false,
  "user_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "images" (
  "id" bigserial PRIMARY KEY,
  "image_name" varchar NOT NULL,
  "data" bytea NOT NULL,
  "product" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "cartitems" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "cart" uuid  NOT NULL,
  "product" uuid  NOT NULL,
  "quantity" bigint NOT NULL,
  "price" float NOT NULL,
  "currency" varchar NOT NULL,
  "sub_total" float NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "carts" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "total_price" float NOT NULL
);

CREATE TABLE "orders" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_name" varchar NOT NULL,
  "user_id" uuid NOT NULL,
  "total_price" float NOT NULL,
  "delivery_address" varchar NOT NULL,
  "country" varchar NOT NULL,
  "status" varchar NOT NULL DEFAULT 'processing',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "orderitems" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "item_name" varchar NOT NULL,
  "item_sub_total" float NOT NULL,
  "quantity" bigint NOT NULL,
  "item_id" uuid NOT NULL,
  "order_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "payments" (
  "id" varchar PRIMARY KEY,
  "amount" float NOT NULL,
  "currency" varchar NOT NULL,
  "status" varchar NOT NULL DEFAULT 'processing',
  "user_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "products" ("category");

CREATE INDEX ON "products" ("product_name");

COMMENT ON COLUMN "products"."price" IS 'must be positive';

COMMENT ON COLUMN "cartitems"."price" IS 'must be positive';

ALTER TABLE "products" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "images" ADD FOREIGN KEY ("product") REFERENCES "products" ("id");

ALTER TABLE "cartitems" ADD FOREIGN KEY ("cart") REFERENCES "carts" ("id");

ALTER TABLE "cartitems" ADD FOREIGN KEY ("product") REFERENCES "products" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orderitems" ADD FOREIGN KEY ("item_id") REFERENCES "products" ("id");

ALTER TABLE "orderitems" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");