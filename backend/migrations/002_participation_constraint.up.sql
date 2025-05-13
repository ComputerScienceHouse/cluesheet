ALTER TABLE user_participation
DROP COLUMN IF EXISTS id;

ALTER TABLE user_participation
ADD CONSTRAINT user_participation_pk PRIMARY KEY (cluesheet_id, ipa_uid);
