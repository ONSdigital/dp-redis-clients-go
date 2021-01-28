package dpredis

import (
	"encoding/json"
	"time"
)

const (
	dateTimeFMT = "2006-01-02T15:04:05.000Z"
)

// Session defines the format of a user session object as it is stored in the cache.
type Session struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Start        time.Time `json:"start"`
	LastAccessed time.Time `json:"lastAccessed"`
}

type jsonModel struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Start        string `json:"start"`
	LastAccessed string `json:"last_accessed"`
}

// MarshalJSON is custom JSON marshaller for the Session object ensuring the date fields are marshalled into the correct format
func (s *Session) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonModel{
		ID:           s.ID,
		Email:        s.Email,
		Start:        s.Start.Format(dateTimeFMT),
		LastAccessed: s.LastAccessed.Format(dateTimeFMT),
	})
}
