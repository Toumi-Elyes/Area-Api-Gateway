// Package logger Logger package handle all the informations linked to the logger
package logger

import "go.uber.org/zap"

var Log *zap.SugaredLogger = nil

// CreateNewLogger Create a new logger for the api-gateway
func CreateNewLogger() error {
	log, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	Log = log.Sugar()
	return nil
}
