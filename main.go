package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/joshuabl97/sitemap-builder/handlers"
	"github.com/rs/zerolog"
	"golang.org/x/net/html"
)

var portNum = flag.String("port_number", "8080", "The port number the server runs on")
var timeZone = flag.String("timezone", "Etc/Greenwich", "An official TZ identifier")

func main() {
	// instantiate logger
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// setting timezone
	loc, err := time.LoadLocation(*timeZone)
	if err != nil {
		l.Error().Msg("Couldn't determine timezone, using local machine time")
	} else if err == nil {
		time.Local = loc
	}

	// make the logs look pretty
	l = l.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// create a custom logger that wraps the zerolog.Logger we instantiated/customized above
	errorLog := &zerologLogger{l}

	// generating a type HandlerHelper for handler methods
	hh := handlers.CreateHandlerHelper(&l)

	// defining the serve mux (sm)
	sm := chi.NewRouter()

	// define middlewares
	sm.Use(hh.LoggingMiddleware)

	// registering the handlers on the serve mux (sm)
	sm.Get("/healthz", handlers.HealthzHandler)

	// create a new server
	s := http.Server{
		Addr:         ":" + *portNum,           // configure the bind address
		Handler:      sm,                       // set the default handler
		IdleTimeout:  120 * time.Second,        // max duration to wait for the next request when keep-alives are enabled
		ReadTimeout:  5 * time.Second,          // max duration for reading the request
		WriteTimeout: 10 * time.Second,         // max duration before returning the request
		ErrorLog:     log.New(errorLog, "", 0), // set the logger for the server
	}

	// this go function starts the server
	// when the function is done running, that means we need to shutdown the server
	// we can do this by killing the program, but if there are requests being processed
	// we want to give them time to complete
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal().Err(err)
		}
	}()

	// sending kill and interrupt signals to os.Signal channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// does not invoke 'graceful shutdown' unless the signalChannel is closed
	<-sigChan

	l.Info().Msg("Received terminate, graceful shutdown")

	// this timeoutContext allows the server 30 seconds to complete all requests (if any) before shutting down
	timeoutCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = s.Shutdown(timeoutCtx)
	if err != nil {
		l.Fatal().Err(err).Msg("Shutdown exceeded timeout")
		os.Exit(1)
	}
}

// custom logger type that wraps zerolog.Logger
type zerologLogger struct {
	logger zerolog.Logger
}

// implement the io.Writer interface for our custom logger.
func (l *zerologLogger) Write(p []byte) (n int, err error) {
	l.logger.Error().Msg(string(p))
	return len(p), nil
}

func getPageLinks(url string) {
	// fetch the HTML content from the URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	// parse the HTML content
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			// End of the document
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				// Found an <a> tag
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						fmt.Println("Link:", attr.Val)
					}
				}
			}
		}
	}
}
