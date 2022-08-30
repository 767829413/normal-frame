package v1

import (
	"context"

	"github.com/767829413/normal-frame/internal/apiserver/model"
	"github.com/767829413/normal-frame/internal/pkg/store"
)

type UserSrv interface {
	Create(ctx context.Context, user *model.User) error
	// Update(ctx context.Context, user *model.User) error
	// Delete(ctx context.Context, username string) error
	// Get(ctx context.Context, username string) (*model.User, error)
	// List(ctx context.Context) (*model.UserList, error)
}

type userService struct {
	store store.Factory
}

var _ UserSrv = (*userService)(nil)

func newUsers(srv *service) *userService {
	return &userService{store: srv.store}
}

func (u *userService) Create(ctx context.Context, user *model.User) error {
	u.store.GetDb().Model(&model.User{}).Create(user)
	return nil
}

// func (u *userService) Get(ctx context.Context, username string) (*model.User, error) {
// 	return nil, nil
// }
// func (u *userService) Update(ctx context.Context, user *model.User) error {
// 	return nil
// }
// func (u *userService) Delete(ctx context.Context, username string) error {
// 	return nil
// }
// func (u *userService) List(ctx context.Context) (*model.UserList, error) {
// 	return nil, nil
// }
