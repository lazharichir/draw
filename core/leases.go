package core

import (
	"time"
)

type LeaseStatus string

const (
	LeaseStatusPending    LeaseStatus = "pending"
	LeaseStatusActive     LeaseStatus = "active"
	LeaseStatusExpired    LeaseStatus = "expired"
	LeaseStatusTerminated LeaseStatus = "terminated"
)

type Lease struct {
	ID            string
	LeaseholderID int64
	CanvasID      int64
	Area          Area
	Status        LeaseStatus
	Start         time.Time
	End           time.Time
	Price         int64
	Metadata      Metadata
	UpdatedAt     time.Time
	UpdatedBy     int64
	CreatedAt     time.Time
	CreatedBy     int64
}

func (l Lease) IsActiveAt(at time.Time) bool {
	return l.Status == LeaseStatusActive && at.After(l.Start) && at.Before(l.End)
}
