package validation

import (
	"github.com/767829413/normal-frame/internal/apiserver/model"
)

func Create(user model.User) (err error) {
	err = IsValidPassword(user.Password)
	return
}
