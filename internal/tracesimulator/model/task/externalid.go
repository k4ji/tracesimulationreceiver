package task

import (
	"fmt"
	"regexp"
)

const externalIDPattern = `^[a-zA-Z0-9_-]+$`

type ExternalID struct {
	string
}

// NewExternalID creates a new ExternalID after validating the input
func NewExternalID(id string) (*ExternalID, error) {
	var validID = regexp.MustCompile(externalIDPattern)
	if !validID.MatchString(id) {
		return nil, fmt.Errorf("invalid external ID: %s", id)
	}
	return &ExternalID{id}, nil
}

// Value returns the string value of the ExternalID
func (e ExternalID) Value() string {
	return e.string
}
