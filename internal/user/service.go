package user

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) CreateUser(req CreateUserRequest) (*User, error) {
	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
	}
	err := s.repo.Create(user)
	return user, err
}
