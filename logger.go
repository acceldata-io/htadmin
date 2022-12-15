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
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger(loggerEnabled bool) {
	writerSyncer := getLogWriter(loggerEnabled)
	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writerSyncer, zapcore.InfoLevel)

	logger := zap.New(core)
	sugarLogger = logger.Sugar()
	zapLogger = logger.Sugar().Desugar()
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter(loggerEnabled bool) zapcore.WriteSyncer {
	dummyFile := ioutil.Discard
	logFilePath := logDirPath + "/" + logFileName
	if !loggerEnabled {
		// To disable logging we will send logs to a discard file
		return zapcore.AddSync(dummyFile)
	}
	// Ignored 'os.O_APPEND' to avoid log file growth
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		// If we cannot open/create the log file we will discard logs
		return zapcore.AddSync(dummyFile)
	}
	return zapcore.AddSync(file)
}
