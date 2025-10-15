package log

import (
	stdlog "log"
)

// For now, use stdlib log with sane defaults. Can be replaced with structured logger.
func Infof(format string, v ...any) {
	stdlog.Printf("INFO "+format, v...)
}

func Errorf(format string, v ...any) {
	stdlog.Printf("ERROR "+format, v...)
}
