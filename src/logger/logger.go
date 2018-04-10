package logger

import (
	"io"
	"log"
)

// Init - Initialize log handlers with output location
func Init(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) (*log.Logger, *log.Logger, *log.Logger) {

	Info := log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning := log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error := log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return Info, Warning, Error
}
