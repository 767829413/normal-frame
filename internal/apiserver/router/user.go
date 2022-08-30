package apiserver

import (
	userContr "github.com/767829413/normal-frame/internal/apiserver/controller/v1/user"
	"github.com/767829413/normal-frame/internal/pkg/store"
	"github.com/gin-gonic/gin"
)

func installTester(g *gin.Engine) *gin.Engine {
	storeIns := store.GetMySQLIncOr(nil)
	v1 := g.Group("/v1")
	{
		userv1 := v1.Group("users")
		{
			userController := userContr.NewUserController(storeIns)
			userv1.POST("", userController.Create)
		}
	}
	return g
}
