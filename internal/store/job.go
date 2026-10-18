package store 

import (
	"encoding/json"
	"errors"
	"time"
)
var (
	ErrNoJobs   = errors.New("no jobs available")
	ErrNotFound = errors.New("job not found")
)

type Status string

type Job struct {
	ID             int64
	TenantID       string
	Queue          string
	Kind           string
	Payload        json.RawMessage
	Status         Status
	Priority       int16
	Attempt        int
	MaxAttempts    int
	RunAt          time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	IdempotencyKey *string
	LastError      *string
}