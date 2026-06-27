package reservation

import (
	"errors"
	"spotsync/internal/domain/reservation/dto"
)

type Service interface {
	MakeReservation(userID uint, req dto.ReserveRequest) (*Reservation, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	CancelReservation(userID uint, userRole string, resID uint) error
	GetAllReservations() ([]Reservation, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) MakeReservation(userID uint, req dto.ReserveRequest) (*Reservation, error) {
	return s.repo.CreateAtomic(userID, req.ZoneID, req.LicensePlate)
}

func (s *service) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	list, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var res []dto.MyReservationResponse
	for _, item := range list {
		res = append(res, dto.MyReservationResponse{
			ID:           item.ID,
			LicensePlate: item.LicensePlate,
			Status:       item.Status,
			Zone: dto.ZoneSummary{
				ID:   item.Zone.ID,
				Name: item.Zone.Name,
				Type: item.Zone.Type,
			},
			CreatedAt: item.CreatedAt,
		})
	}
	return res, nil
}

func (s *service) CancelReservation(userID uint, userRole string, resID uint) error {
	res, err := s.repo.FindByID(resID)
	if err != nil {
		return errors.New("reservation not found")
	}

	// Only the user who made the reservation or an admin can cancel it
	if userRole != "admin" && res.UserID != userID {
		return errors.New("unauthorized to cancel this reservation")
	}

	if res.Status != "active" {
		return errors.New("reservation is already cancelled or completed")
	}

	return s.repo.UpdateStatus(resID, "cancelled")
}

func (s *service) GetAllReservations() ([]Reservation, error) {
	return s.repo.FindAllWithPreload()
}
