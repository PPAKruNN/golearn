package service

import (
	"fmt"

	"github.com/google/uuid"
)

type AuthRepository interface {
	RegisterToken(token string, accountId int)
	DecodeToken(token string) (int, error)
	Reset() error
}

type AuthService struct {
	Repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s AuthService) CreateToken(accountId int) string {
	t := uuid.NewString()
	s.Repo.RegisterToken(t, accountId)

	return t
}

func (s AuthService) DecodeToken(token string) (int, error) {

	accountId, err := s.Repo.DecodeToken(token)
	if err != nil {
		return 0, fmt.Errorf("AccountId returned is not found. AccountId returned: %d", accountId)
	}

	return accountId, nil
}
