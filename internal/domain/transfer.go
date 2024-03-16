package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransferStatus string

const (
	StatusPending   TransferStatus = "PENDING"
	StatusCompleted TransferStatus = "COMPLETED"
	StatusFailed    TransferStatus = "FAILED"
)

type Transfer struct {
	gorm.Model
	ID            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	FromAccountID uuid.UUID      `gorm:"type:uuid"`
	ToAccountID   uuid.UUID      `gorm:"type:uuid"`
	Amount        float64        `gorm:"type:numeric"`
	Status        TransferStatus `gorm:"type:varchar(50)"`
	TransactionID string         `gorm:"type:varchar(100);uniqueIndex"` // For tracking with external systems
}
