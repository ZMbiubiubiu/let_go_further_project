package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routers(),
		// ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Create a shutdownError channel. we will use this to receive any errors
	// returned by the graceful Shutdown() function
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		// 拦截中断信号
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		// 优雅关闭，HTTP服务器
		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// 优雅关闭，后台goroutines
		app.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})

		app.wg.Wait()

		// graceful shutdwon finish!
		shutdownError <- nil
	}()

	app.logger.PrintInfo("starting server ", map[string]string{"addr": srv.Addr, "env": app.config.env})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	// At this point we know that the graceful shutdown completed successfully and
	// we log a "stopped server" message
	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
