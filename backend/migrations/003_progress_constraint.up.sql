ALTER TABLE user_progress
DROP COLUMN IF EXISTS id;

ALTER TABLE user_progress
Add CONSTRAINT user_progress_pk PRIMARY KEY (clue_id, ipa_uid);
