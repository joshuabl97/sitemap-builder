package handlers

import "github.com/rs/zerolog"

// fields passed to handlers
type HandlerHelper struct {
	l   *zerolog.Logger
	url string
}

func CreateHandlerHelper(l *zerolog.Logger) HandlerHelper {
	return HandlerHelper{l, "https://www.iana.org/help/example-domains"}
}
