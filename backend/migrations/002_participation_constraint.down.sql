ALTER TABLE user_participation
ADD COLUMN id UUID;

ALTER TABLE user_participation
DROP CONSTRAINT IF EXISTS user_participation_pk;
