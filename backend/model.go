package main

import (
	"github.com/google/uuid"
	"time"
)

type Cluesheet struct {
	Id         uuid.UUID
	Name       *string
	Origin_id  *uuid.UUID /* originating cluesheet */
	Created_by *string    /* ipa unique id */
	Created_at *time.Time /* no timezone by default */
	Edited_by  *string    /* ipa unique id */
	Edited_at  *time.Time /* no timezone */
	Visibility *string    /* TODO: int key? or defined enum? I think SQL enums are difficult to work with in migrations */
	Owners     *[]string
	Groups     *[]string
}
