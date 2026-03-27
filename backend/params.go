package main

import (
	"github.com/google/uuid"
)

type PostCluesheetParams struct {
	Name       string
	Origin     *uuid.UUID
	Creator    string
	Visibility *string
	Owners     []string
	Groups     []string
}
