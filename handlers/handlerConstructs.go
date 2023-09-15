package handlers

import "github.com/rs/zerolog"

// fields passed to handlers
type HandlerHelper struct {
	l    *zerolog.Logger
	port string
}

func CreateHandlerHelper(logger *zerolog.Logger, portNum *string) HandlerHelper {
	return HandlerHelper{l: logger, port: *portNum}
}
