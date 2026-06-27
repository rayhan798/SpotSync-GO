package zone

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(zone *ParkingZone) error
	FindAll() ([]ParkingZone, map[uint]int, error)
	FindByID(id uint) (*ParkingZone, int, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *repository) FindAll() ([]ParkingZone, map[uint]int, error) {
	var zones []ParkingZone
	if err := r.db.Find(&zones).Error; err != nil {
		return nil, nil, err
	}

	// (total_capacity - active_reservations)
	activeCounts := make(map[uint]int)
	type Result struct {
		ZoneID uint
		Count  int
	}
	var results []Result

	// reservations
	r.db.Table("reservations").Select("zone_id, count(*) as count").Where("status = ?", "active").Group("zone_id").Scan(&results)

	for _, res := range results {
		activeCounts[res.ZoneID] = res.Count
	}

	return zones, activeCounts, nil
}

func (r *repository) FindByID(id uint) (*ParkingZone, int, error) {
	var zone ParkingZone
	if err := r.db.First(&zone, id).Error; err != nil {
		return nil, 0, err
	}

	var activeCount int64
	r.db.Table("reservations").Where("zone_id = ? AND status = ?", id, "active").Count(&activeCount)

	return &zone, int(activeCount), nil
}
