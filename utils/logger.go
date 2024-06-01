package utils

import "log"

type Logger interface {
	Debugf(format string, params ...interface{})

	Infof(format string, params ...interface{})

	Warnf(format string, params ...interface{})

	Errorf(format string, params ...interface{})

	Debug(v ...interface{})

	Info(v ...interface{})

	Warn(v ...interface{})

	Error(v ...interface{})
}

type DefaultLogger struct {
}

func (d *DefaultLogger) Debugf(format string, params ...interface{}) {
	log.Printf("Debug:"+format, params...)
}

func (d *DefaultLogger) Infof(format string, params ...interface{}) {
	log.Printf("Info:"+format, params...)
}

func (d *DefaultLogger) Warnf(format string, params ...interface{}) {
	log.Printf("Warn:"+format, params...)
}

func (d *DefaultLogger) Errorf(format string, params ...interface{}) {
	log.Printf("Error:"+format, params...)
}

func (d *DefaultLogger) Debug(v ...interface{}) {
	log.Println(v...)
}

func (d *DefaultLogger) Info(v ...interface{}) {
	log.Println(v...)
}

func (d *DefaultLogger) Warn(v ...interface{}) {
	log.Println(v...)
}

func (d *DefaultLogger) Error(v ...interface{}) {
	log.Println(v...)
}
