package main

import (
	"context"
	"fmt"
	"golang-boiler-plate/config"
	"golang-boiler-plate/handlers"
	"golang-boiler-plate/middleware"
	"golang-boiler-plate/utils/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

var (
	serviceName           = "golang-boiler-plate"
	serverAddress         = fmt.Sprintf("0.0.0.0:%s", config.Port)
	serverShutdownTimeout = 30 * time.Second
	serverRequestTimeout  = time.Duration(config.RequestTimeout) * time.Second
)

func main() {
	log.SetFlags(0)

	logger.ConfigureLogger(config.LogFormat, config.DebugMode)

	r := chi.NewRouter()

	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CorrelationID)

	userHandler := handlers.NewUserHandler().Routes()

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Mount("/users", userHandler)
	})

	// display all routes in log
	walkFunc := func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		logger.Info("registering [%s] %s", method, route)
		return nil
	}
	if err := chi.Walk(r, walkFunc); err != nil {
		logger.Fatal("walk function error : %v\n", err)
	}

	// run server and setup gracefully shutdown
	quit := make(chan os.Signal, 1)
	srv := &http.Server{
		Addr:         serverAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      r,
	}
	srv.WriteTimeout = serverRequestTimeout
	go func() {
		logger.Info("starting %s on %s", serviceName, serverAddress)
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal("could not serve on %s: %v", serverAddress, err)
		}
	}()
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down '%s', waiting for ongoing requests to complete... [force shutdown in %s]", serviceName, serverShutdownTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("some ongoing requests may be forced to close due to %v", err)
	}
	logger.Info("shutdown '%s' successfully", serviceName)
}
