package sfu

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05 -070000"
	IdLen      = 10
)

func generateSessionId() string {
	longId := fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.New().String())))
	return longId[:IdLen]
}

func generateParticipantId() string {
	return generateSessionId()
}

func generateCreationDateTime() string {
	return time.Now().Format(timeFormat)
}
