// Acceldata Inc. and its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// 	Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
