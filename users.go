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
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type account struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func createUser(c *gin.Context) {
	//
	var user account
	c.BindJSON(&user)

	//
	rand.Seed(time.Now().UnixNano())
	user.Password = randSeq(32)
	ht := &htdata{
		File:    htFile,
		Account: user,
	}

	rCode := http.StatusOK
	rText := "Success\n" + user.Password

	if err := ht.SetPassword(); err != nil {
		rCode = http.StatusConflict
		rText = err.Error()
	}

	c.String(rCode, rText)
}

func deleteUser(c *gin.Context) {
	//
	var user account
	c.BindJSON(&user)

	ht := &htdata{
		File:    htFile,
		Account: user,
	}

	rCode := http.StatusOK
	rText := fmt.Sprintf("User %q deleted successfully", user.Name)

	if err := ht.DeleteUser(); err != nil {
		rCode = http.StatusConflict
		rText = err.Error()
	}

	c.String(rCode, rText)
}
