package logpack

import (
	"log"
	"os"
)

type LogPack struct {
	Info  *log.Logger
	Err   *log.Logger
	Fatal *log.Logger
}

func NewLogger() *LogPack {

	return &LogPack{
		Info:  log.New(os.Stdout, "INFO\t", log.LstdFlags),
		Err:   log.New(os.Stderr, "ERROR\t", log.Lshortfile|log.LstdFlags),
		Fatal: log.New(os.Stderr, "FATAL\t", log.Lshortfile|log.LstdFlags),
	}
}
