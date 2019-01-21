package api

import (
	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Serve(api API, address string) {

	// configure app routes
	handler := routing(api)

	// don't let a panic crash the server.
	handler = handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(handler)

	// log all requests using logrus logger
	handler = handlers.LoggingHandler(
		log.WithField("prefix", "httpd").WriterLevel(log.InfoLevel),
		handler)

	server := &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadHeaderTimeout: 1 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Infof("Start http server on %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}

	log.Info("Server shutdown completed.")
}
