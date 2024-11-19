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
  "cart" uuid NOT NULL,
  "product" uuid NOT NULL,
  "quantity" bigint NOT NULL,
  "price" float NOT NULL,
  "sub_total" float NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "carts" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "total_price" float NOT NULL
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "products" ("category");

CREATE INDEX ON "products" ("product_name");

COMMENT ON COLUMN "products"."price" IS 'must be positive';

COMMENT ON COLUMN "cartitems"."price" IS 'must be positive';

ALTER TABLE "products" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "images" ADD FOREIGN KEY ("product") REFERENCES "products" ("id");

ALTER TABLE "cartitems" ADD FOREIGN KEY ("cart") REFERENCES "carts" ("id");

ALTER TABLE "cartitems" ADD FOREIGN KEY ("product") REFERENCES "products" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

