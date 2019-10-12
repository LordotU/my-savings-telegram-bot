package app

import "go.uber.org/zap"

func GetLogger(debug bool) (*zap.Logger, error) {
	var config zap.Config

	if debug {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	return config.Build()
}
