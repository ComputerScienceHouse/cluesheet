ALTER TABLE user_progress
ADD COLUMN id UUID;

ALTER TABLE user_progress
DROP CONSTRAINT IF EXISTS user_progress_pk;
