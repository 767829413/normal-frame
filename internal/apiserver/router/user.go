package apiserver

import (
	"github.com/gin-gonic/gin"
)

func installTester(g *gin.Engine) *gin.Engine {
	v1 := g.Group("/v1")
	{
		userv1 := v1.Group("users")
		{
			userv1.POST("")
		}
	}
	return g
}
