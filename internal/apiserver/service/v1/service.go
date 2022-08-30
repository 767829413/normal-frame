package v1

import (
	"github.com/767829413/normal-frame/internal/pkg/store"
)

type Service interface {
	Users() UserSrv
}

type service struct {
	store store.Factory
}

// NewService returns Service interface.
func NewService(st store.Factory) Service {
	return &service{
		store: st,
	}
}

func (s *service) Users() UserSrv {
	return newUsers(s)
}
