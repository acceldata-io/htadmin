package main

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func setupRouter(credsData map[string]string, enableReqLogger bool) *gin.Engine {
	r := gin.Default()

	if enableReqLogger {
		//
		r.Use(ginzap.Ginzap(zapLogger, time.RFC3339, true))

		// Logs all panic to error log
		//   - stack means whether output the stack info.
		r.Use(ginzap.RecoveryWithZap(zapLogger, true))
	}

	// Adds basic auth middleware to the path
	authorized := r.Group("/", gin.BasicAuth(credsData))

	// Create user POST Req
	authorized.POST("uac/create", createUser)
	// Delete user POST Req
	authorized.POST("uac/delete", deleteUser)

	return r
}
