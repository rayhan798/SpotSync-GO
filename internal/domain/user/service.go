package user

import (
	"fmt"
	"spotsync/internal/auth"
	"spotsync/internal/domain/user/dto"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

type Service interface {
	CreateUser(req dto.CreateRequest) (*dto.RegisterRes, error)
	LoginUser(req dto.LoginRequest) (*dto.LoginRes, error)
}

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) Service {
	return &service{repo, jwtService}
}

func (s *service) CreateUser(req dto.CreateRequest) (*dto.RegisterRes, error) {
	// ১. যদি রিকোয়েস্টে কোনো রোল না পাঠানো হয়, তবে ডিফল্ট হিসেবে "driver" সেট হবে
	userRole := req.Role
	if userRole == "" {
		userRole = "driver"
	}

	// 🔒 ২. অ্যাডমিন চেকিং লজিক: যদি নতুন ইউজার 'admin' হতে চায়
	if userRole == "admin" {
		// ডাটাবেজে অলরেডি কোনো অ্যাডমিন আছে কিনা তা চেক করার জন্য রেপোজিটরি মেথড কল
		adminExists, err := s.repo.CheckAdminExists()
		if err != nil {
			return nil, fmt.Errorf("failed to check admin existence: %w", err)
		}

		// যদি অলরেডি একজন অ্যাডমিন থাকে, তবে সরাসরি এরর রিটার্ন করবে
		if adminExists {
			return nil, fmt.Errorf("admin already exists") // 👈 এই এররটি হ্যান্ডলারে চলে যাবে
		}
	}

	user := User{
		Name:  req.Name,
		Email: req.Email,
		Role:  userRole,
	}

	err := user.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	err = s.repo.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	response := dto.RegisterRes{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	return &response, nil
}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.LoginRes, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	err = user.CheckPassword(req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token using the jwtService
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	response := dto.LoginRes{
		Token: token,
		User: dto.UserSummary{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return &response, nil
}
