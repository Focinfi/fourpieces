package fourpieces

import "log"
import "os"

type logger struct {
	debug bool
	*log.Logger
}

func (l *logger) printf(format string, a ...interface{}) {
	if l.debug {
		l.Logger.Printf(format, a...)
	}
}

func (l *logger) println(v ...interface{}) {
	if l.debug {
		l.Logger.Println(v...)
	}
}

func (l *logger) fatalf(format string, a ...interface{}) {
	l.Logger.Fatalf(format, a...)
}

var defaultLogger = &logger{Logger: log.New(os.Stdout, "[FourPieces] ", 0), debug: true}
