package middleware

import "context"

type Logger struct {
}

func (l *Logger) Validate(c context.Context, args ...interface{}) {
}

func (l *Logger) Handler(c context.Context, args ...interface{}) {
	// TODO
}