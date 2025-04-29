DROP INDEX IF EXISTS "idx_spaces_creator";
DROP INDEX IF EXISTS "idx_spaces_name";
DROP INDEX IF EXISTS "idx_spaces_created_at";
DROP INDEX IF EXISTS "idx_spaces_creator_created";
DROP INDEX IF EXISTS "idx_spaces_name_trgm";

DROP INDEX IF EXISTS "idx_space_members_space_id";
DROP INDEX IF EXISTS "idx_space_members_user_id";
DROP INDEX IF EXISTS "idx_space_members_space_joined";
DROP INDEX IF EXISTS "idx_space_members_username";

DROP TABLE IF EXISTS "space_members";

DROP TABLE IF EXISTS "spaces";

DROP EXTENSION IF EXISTS "pg_trgm";
