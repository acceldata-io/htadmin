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
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// If we cannot open/create the log file we will discard logs
		return zapcore.AddSync(dummyFile)
	}
	return zapcore.AddSync(file)
}
