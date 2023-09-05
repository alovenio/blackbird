package sfu

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

// checkSessionId checks whether a given identifier conforms with
// expected format. An error will be returned if the given id
// is deemed invalid.
func checkSessionId(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return err
	}
	return nil
}

// checkNotBlank checks whether a given name is not blank. An
// error will be returned if the given name is blank.
func checkNotBlank(name string) error {
	if len(strings.TrimSpace(name)) == 0 {
		return fmt.Errorf("name must not be blank")
	}
	return nil
}

// generateSessionId generates and returns a new UUID
func generateSessionId() string {
	return uuid.New().String()
}
