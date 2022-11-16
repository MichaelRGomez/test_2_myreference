// Filename: test2/cmd/api/server.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	//Setting up the HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//shutdown should return it's errors to this channel
	shutdownError := make(chan error)

	//goroutine
	go func() {
		//Creating a quit channel which carries os.Signal values
		quit := make(chan os.Signal, 1)

		//Listening for SIGINT and SIGTERM signals
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		//Blocking until a signal is recieved
		s := <-quit

		//Logging the message
		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		//context with 20-second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		//calling the shutdown function
		shutdownError <- srv.Shutdown(ctx)
	}()

	//Starting the server
	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	//Check if the shutdown process has been initiated
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	//Blocking for notification from shudtown() function
	err = <-shutdownError
	if err != nil {
		return err
	}

	//Graceful shutdown was successful
	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})
	return nil
}
