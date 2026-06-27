package reservation

import (
	"spotsync/internal/domain/zone"
	"time"
)

type Reservation struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	UserID       uint             `gorm:"not null" json:"user_id"`
	ZoneID       uint             `gorm:"not null" json:"zone_id"`
	Zone         zone.ParkingZone `gorm:"foreignKey:ZoneID" json:"zone"`
	LicensePlate string           `gorm:"type:varchar(15);not null" json:"license_plate"`
	Status       string           `gorm:"type:varchar(15);not null;default:'active'" json:"status"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

func (Reservation) TableName() string {
	return "reservations"
}
