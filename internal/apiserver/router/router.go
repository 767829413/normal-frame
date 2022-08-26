package apiserver

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(g *gin.Engine) {
	installTester(g)
}
