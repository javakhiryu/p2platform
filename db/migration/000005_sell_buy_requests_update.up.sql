ALTER TABLE sell_requests
ADD COLUMN space_id uuid NOT NULL REFERENCES spaces(space_id) ON DELETE CASCADE;
