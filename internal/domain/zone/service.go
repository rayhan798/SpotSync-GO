package zone

import (
	"spotsync/internal/domain/zone/dto"
)

type Service interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.Create(zone); err != nil {
		return nil, err
	}

	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.TotalCapacity,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}, nil
}

func (s *service) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, activeCounts, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []dto.ZoneResponse
	for _, z := range zones {
		active := activeCounts[z.ID]
		res = append(res, dto.ZoneResponse{
			ID:             z.ID,
			Name:           z.Name,
			Type:           z.Type,
			TotalCapacity:  z.TotalCapacity,
			AvailableSpots: z.TotalCapacity - active,
			PricePerHour:   z.PricePerHour,
			CreatedAt:      z.CreatedAt,
		})
	}
	return res, nil
}

func (s *service) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	z, active, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.ZoneResponse{
		ID:             z.ID,
		Name:           z.Name,
		Type:           z.Type,
		TotalCapacity:  z.TotalCapacity,
		AvailableSpots: z.TotalCapacity - active,
		PricePerHour:   z.PricePerHour,
		CreatedAt:      z.CreatedAt,
		UpdatedAt:      z.UpdatedAt,
	}, nil
}
