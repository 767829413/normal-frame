package v1

import (
	"time"

	v1 "github.com/767829413/normal-frame/internal/apiserver/controller/v1"
	"github.com/767829413/normal-frame/internal/apiserver/model"
	srvv1 "github.com/767829413/normal-frame/internal/apiserver/service/v1"
	"github.com/767829413/normal-frame/internal/apiserver/validation"
	"github.com/767829413/normal-frame/internal/pkg/store"
	"github.com/gin-gonic/gin"
)

// UserController create a user handler used to handle request for user resource.
type UserController struct {
	srv srvv1.Service
}

func NewUserController(st store.Factory) *UserController {
	return &UserController{
		srv: srvv1.NewService(st),
	}
}

// Create add new user to the storage.
func (u *UserController) Create(c *gin.Context) {

	var r model.User
	resp := &v1.Res{State: 1, Msg: "success"}
	// Binding parameters
	if err := c.ShouldBindJSON(&r); err != nil {
		resp.Msg = err.Error()
		resp.State = -1
		resp.WriteResponse(c)
		return
	}
	// Calibrate model data
	if err := validation.Create(r); err != nil {
		resp.Msg = err.Error()
		resp.State = -1
		resp.WriteResponse(c)
		return
	}

	r.Status = 1
	r.LoginedAt = time.Now()

	// Insert the user to the storage.
	if err := u.srv.Users().Create(c, &r); err != nil {
		resp.Msg = err.Error()
		resp.State = -1
		resp.WriteResponse(c)
		return
	}
	resp.WriteResponse(c)
}
