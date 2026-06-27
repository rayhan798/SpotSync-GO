package zone

import (
	"time"
)

type ParkingZone struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(150);not null" json:"name"`
	Type          string    `gorm:"type:varchar(50);not null" json:"type"` // general, ev_charging, covered
	TotalCapacity int       `gorm:"not null" json:"total_capacity"`
	PricePerHour  float64   `gorm:"type:decimal(10,2);not null" json:"price_per_hour"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GORM Table Name Specification
func (ParkingZone) TableName() string {
	return "parking_zones"
}
