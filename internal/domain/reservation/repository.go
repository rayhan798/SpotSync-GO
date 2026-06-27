package reservation

import (
	"errors"
	"spotsync/internal/domain/zone"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrZoneFull = errors.New("parking zone is full")

type Repository interface {
	CreateAtomic(userID uint, zoneID uint, licensePlate string) (*Reservation, error)
	FindByUserID(userID uint) ([]Reservation, error)
	FindByID(id uint) (*Reservation, error)
	UpdateStatus(id uint, status string) error
	FindAllWithPreload() ([]Reservation, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateAtomic(userID uint, zoneID uint, licensePlate string) (*Reservation, error) {
	var newRes Reservation

	// ACID Transaction
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var pZone zone.ParkingZone

		// ১. Row-level locking (FOR UPDATE)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&pZone, zoneID).Error; err != nil {
			return err
		}

		// 2. Active reservations count
		var activeCount int64
		if err := tx.Model(&Reservation{}).Where("zone_id = ? AND status = ?", zoneID, "active").Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. if the active reservations count is greater than or equal to the total capacity, return an error
		if int(activeCount) >= pZone.TotalCapacity {
			return ErrZoneFull
		}

		// 4. new reservation
		newRes = Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}

		if err := tx.Create(&newRes).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &newRes, nil
}

func (r *repository) FindByUserID(userID uint) ([]Reservation, error) {
	var reservations []Reservation
	// GORM Preload
	err := r.db.Preload("Zone").Where("user_id = ?", userID).Order("created_at desc").Find(&reservations).Error
	return reservations, err
}

func (r *repository) FindByID(id uint) (*Reservation, error) {
	var res Reservation
	err := r.db.First(&res, id).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&Reservation{}).Where("id = ?", id).Update("status", status).Error
}

func (r *repository) FindAllWithPreload() ([]Reservation, error) {
	var reservations []Reservation
	err := r.db.Preload("Zone").Order("created_at desc").Find(&reservations).Error
	return reservations, err
}
