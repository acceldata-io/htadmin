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
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var (
	addr         = "127.0.0.1:19978"
	sugarLogger  *zap.SugaredLogger
	zapLogger    *zap.Logger
	enableLogger = true
	logDirPath   = "."
	logFileName  = "htadmin.log"
	htFile       = ".htpasswd"
	apiUserFile  = "creds.yaml"
)

var (
	// Version is the tool version. This should be injected during the build/compile time.
	Version = "0.0.0"
	// BuildID is the tool build epoch time in seconds. This should be injected during the build/compile time.
	BuildID = "0"
)

func main() {
	if len(os.Args) == 2 {
		if checkVersionArg(os.Args[1]) {
			fmt.Println("Version: ", Version)
			fmt.Println("BuildID: ", BuildID)
		} else {
			fmt.Println("unknown argument passed, available options are :")
			fmt.Println("\nUsage: htadmin [OPTION]")
			fmt.Println("-v, --version  Prints the htadmin agent version")
			fmt.Println("\nExamples: ")
			fmt.Println("htadmin -v")
			fmt.Println("htadmin --version")
		}
		return
	} else if len(os.Args) > 2 {
		fmt.Println("unknown arguments passed, available options are :")
		fmt.Println("\nUsage: htadmin [OPTION]")
		fmt.Println("-v, --version  Prints the htadmin agent version")
		fmt.Println("\nExamples: ")
		fmt.Println("htadmin -v")
		fmt.Println("htadmin --version")
		return
	}

	// Read API Users
	apiFileContents, err := os.ReadFile(apiUserFile)
	if err != nil {
		fmt.Printf("ERROR: cannot find the api users file %q\n", apiUserFile)
		fmt.Println("ERROR: ", err.Error())
		os.Exit(1)
	}

	apiUserList := apiUsers{}
	if err := yaml.Unmarshal(apiFileContents, &apiUserList); err != nil {
		fmt.Printf("ERROR: cannot read the api users file %q\n", apiUserFile)
		fmt.Println("ERROR: ", err.Error())
		os.Exit(1)
	}

	credsData := map[string]string{}
	for user, pass := range apiUserList.Users {
		user = strings.TrimSpace(user)
		pass = strings.TrimSpace(pass)
		if user != "" && pass != "" {
			credsData[user] = pass
		}
	}
	if len(credsData) <= 0 {
		fmt.Printf("ERROR: users file %q is invalid\n", apiUserFile)
		os.Exit(1)
	}

	// initialize logger
	initLogger(true)
	defer sugarLogger.Sync()

	// The htadmin server
	htAdminServer := &http.Server{
		Addr:    addr,
		Handler: setupRouter(credsData, false),
	}

	// Initializing the htadmin server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := htAdminServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			sugarLogger.Infof("HTAdmin Server Listener: %s\n", err)
			fmt.Printf("INFO: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugarLogger.Infof("Shutting down server...")
	fmt.Println("INFO: Shutting down server ...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := htAdminServer.Shutdown(ctx); err != nil {
		sugarLogger.Errorf("ERROR: Server forced to shutdown: %s\n", err)
		fmt.Printf("ERROR: Server forced to shutdown: %s\n", err)
	}

	sugarLogger.Infof("Server Exited!")
	fmt.Println("INFO: Server exited!")
}

func checkVersionArg(arg string) bool {
	if arg == "version" || arg == "--version" || arg == "-v" {
		return true
	}
	return false
}
