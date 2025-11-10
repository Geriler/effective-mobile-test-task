package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	ServiceName string    `json:"service_name"`
	Price       int32     `json:"price"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty"`
}

type Filters struct {
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	UserID      uuid.UUID `json:"user_id,omitempty"`
	ServiceName string    `json:"service_name,omitempty"`
}
