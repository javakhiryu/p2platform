CREATE TABLE "spaces" (
    "space_id" uuid UNIQUE PRIMARY KEY,
    "space_name" VARCHAR NOT NULL UNIQUE,
    "hashed_password" VARCHAR NOT NULL,
    "creator_id" BIGINT REFERENCES users(telegram_id),
    "description" VARCHAR NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE TABLE "space_members" (
    "space_id" uuid NOT NULL REFERENCES "spaces"("space_id"),
    "user_id" BIGINT UNIQUE REFERENCES "users"("telegram_id"),
    "username" VARCHAR NOT NULL,
    "joined_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY ("space_id", "user_id")
);

CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE INDEX IF NOT EXISTS "idx_spaces_creator"
  ON "spaces" ("creator_id");

CREATE INDEX IF NOT EXISTS "idx_spaces_name"
  ON "spaces" ("space_name");

CREATE INDEX IF NOT EXISTS "idx_spaces_created_at"
  ON "spaces" ("created_at");

CREATE INDEX IF NOT EXISTS "idx_spaces_creator_created"
  ON "spaces" ("creator_id", "created_at" DESC);

CREATE INDEX IF NOT EXISTS "idx_spaces_name_trgm"
  ON "spaces" USING gin ("space_name" gin_trgm_ops);

CREATE INDEX IF NOT EXISTS "idx_space_members_space_id"
  ON "space_members" ("space_id");

CREATE INDEX IF NOT EXISTS "idx_space_members_user_id"
  ON "space_members" ("user_id");

CREATE INDEX IF NOT EXISTS "idx_space_members_space_joined"
  ON "space_members" ("space_id", "joined_at" DESC);

CREATE INDEX IF NOT EXISTS "idx_space_members_username"
  ON "space_members" ("username");
