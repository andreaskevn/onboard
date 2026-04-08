package config

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLog() error {
	var err error

	Log, err = zap.NewProduction()
	if err != nil {
		return err
	}
	return nil
}
