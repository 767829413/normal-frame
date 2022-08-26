package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

// gzip 压缩
func Gzip(level int) gin.HandlerFunc {
	return gzip.Gzip(level)
}

// Cors add cors headers.
func Cors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("CloudCluster", "ClusterID")
	return cors.New(config)
}

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}

func Dump() gin.HandlerFunc {
	return gindump.Dump()
}
