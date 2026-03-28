package main

import (
	"github.com/google/uuid"
)

type PostCluesheetParams struct {
	Name       string     `json:"name"`
	Origin     *uuid.UUID `json:"origin"`
	Creator    string     `json:"creator"`
	Visibility *string    `json:"visibility"`
	Owners     []string   `json:"owners"`
	Groups     []string   `json:"groups"`
}
