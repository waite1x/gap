package log

import "log"

type DefaultLogger struct {
	opts *Options
}

func (l *DefaultLogger) Log(message string, level Level) {
	if l.opts.Level <= level {
		log.Println(message)
	}
}
