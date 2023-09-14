package handlers

import "github.com/rs/zerolog"

// fields passed to handlers
type HandlerHelper struct {
	l *zerolog.Logger
}

func CreateHandlerHelper(l *zerolog.Logger) HandlerHelper {
	return HandlerHelper{l}
}
