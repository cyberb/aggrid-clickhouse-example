package log

import (
	"fmt"
	"go.uber.org/zap"
)

func Default() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.TimeKey = ""
	//config.EncoderConfig.TimeKey = ""
	config.EncoderConfig.ConsoleSeparator = " "
	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("could not init zap logger: %v", err))
	}

	return logger
}
