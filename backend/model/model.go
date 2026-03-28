package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Cluesheet struct {
	Id         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Origin_id  *uuid.UUID `json:"origin_id"`  /* originating cluesheet */
	Created_by string     `json:"created_by"` /* ipa unique id */
	Created_at time.Time  `json:"created_at"` /* no timezone by default */
	Edited_by  string     `json:"edited_by"`  /* ipa unique id */
	Edited_at  time.Time  `json:"edited_at"`  /* no timezone */
	Visibility string     `json:"visibility"` /* TODO: int key? or defined enum? I think SQL enums are difficult to work with in migrations */
	Owners     []string   `json:"owners"`
	Groups     []string   `json:"groups"`
	Clues      *[]Clue    `json:"clues,omitzero"`
}

// RuleInfo is the resolved rule returned to clients, mapping clue_rules fields to frontend expectations.
type RuleInfo struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

type Clue struct {
	Id          uuid.UUID        `json:"id"`
	Description string           `json:"description"`
	Rule_id     *uuid.UUID       `json:"rule_id"`
	Rule_params *json.RawMessage `json:"rule_params"`
	Origin_id   uuid.UUID        `json:"origin_id"` /* originating cluesheet */
	Created_by  string           `json:"created_by"`
	Created_at  time.Time        `json:"created_at"`
	Edited_by   string           `json:"edited_by"`
	Edited_at   time.Time        `json:"edited_at"`
	Tags        []string         `json:"tags"`
	Children    []*Clue          `json:"children,omitzero"`
	Rule        *RuleInfo        `json:"rule,omitempty"`
	Completions *int             `json:"completions,omitempty"`
}

type UserParticipation struct {
	Cluesheet_id uuid.UUID `json:"cluesheet_id"`
	Ipa_uid      string    `json:"ipa_uid"`
	Hidden       bool      `json:"hidden"`
}

type UserProgress struct {
	Ipa_uid     string    `json:"ipa_uid"`
	Clue_id     uuid.UUID `json:"clue_id"`
	Completions int       `json:"completions"` // TODO this should become a double
}

type ClueRelation struct {
	Id        uuid.UUID `json:"id"`
	Parent_id uuid.UUID `json:"parent_id"`
	Child_id  uuid.UUID `json:"child_id"`
}
