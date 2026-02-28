/* todo constraints and the like */
CREATE TABLE cluesheet (
    id UUID,
    name text,
    origin_id UUID, /* originating cluesheet */
    created_by text, /* ipa unique id */
    created_at timestamp, /* no timezone by default */
    edited_by text, /* ipa unique id */
    edited_at timestamp, /* no timezone */
    visibility text, /* TODO: int key? or defined enum? I think SQL enums are difficult to work with in migrations */
    owners text[],
    groups text[]
);

CREATE TABLE clue (
    id UUID,
    description text,
    rule_id UUID,
    rule_params jsonb,
    origin_id UUID, /* originating cluesheet */
    created_by text,
    created_at timestamp,
    edited_by text,
    edited_at timestamp,
    tags text[]
);

CREATE TABLE clue_relation (
    id UUID,
    parent_id UUID,
    child_id UUID
);

CREATE TABLE user_progress (
    id UUID,
    ipa_uid text, /* ipa unique id */
    clue_id UUID,
    completions integer
);

CREATE TABLE clue_rules (
    id UUID,
    name text,
    jq_expr text,
    req_params text /* TODO */
);

CREATE TABLE user_participation (
    id UUID,
    cluesheet_id UUID,
    ipa_uid text, /* ipa unique id */
    hidden boolean
);
