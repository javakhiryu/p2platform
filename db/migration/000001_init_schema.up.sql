CREATE TABLE "users" (
  "telegram_id" BIGINT PRIMARY KEY,
  "tg_username" VARCHAR NOT NULL,
  "first_name" VARCHAR,
  "last_name" VARCHAR,
  "photo_url" VARCHAR,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

