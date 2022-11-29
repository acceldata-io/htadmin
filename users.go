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
