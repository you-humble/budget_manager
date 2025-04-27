package user

import "log/slog"

type Repository interface {
	Save(u User) (int64, error)
	FindByLogin(login string) (User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Save(u User) (User, error) {
	id, err := s.repo.Save(u)
	if err != nil {
		slog.Error("failed to save user", slog.String("error", err.Error()))
		return User{}, err
	}

	u.ID = id
	return u, nil
}

func (s *service) FindByLogin(login string) (User, error) {
	u, err := s.repo.FindByLogin(login)
	if err != nil {
		slog.Error("failed to find user", slog.String("error", err.Error()))
		return User{}, err
	}
	return u, nil
}
