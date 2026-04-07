package model

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	DisplayName string    `json:"displayName" db:"display_name"`
	Description string    `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
	Visible     bool      `json:"visible" db:"visible"`
	CategoryID  uuid.UUID `json:"categoryId" db:"category_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type Category struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Connection struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ApplicationID uuid.UUID `json:"applicationId" db:"application_id"`
	MerchantID    string    `json:"merchantId" db:"merchant_id"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

// Connection statuses.
const (
	ConnectionStatusActive   = "active"
	ConnectionStatusInactive = "inactive"
)

type Review struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ApplicationID uuid.UUID `json:"applicationId" db:"application_id"`
	MerchantID    string    `json:"merchantId" db:"merchant_id"`
	Rating        int       `json:"rating" db:"rating"`
	Body          string    `json:"body" db:"body"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}
