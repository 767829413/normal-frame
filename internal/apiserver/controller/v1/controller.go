package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Res struct {
	State int         `json:"state"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
}

func (r *Res) WriteResponse(c *gin.Context) {
	c.JSON(http.StatusOK, r)
}
