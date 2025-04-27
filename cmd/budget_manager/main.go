package main

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"budget_manager/internal/database/sqlite"
	"budget_manager/internal/router"
	"budget_manager/internal/session"
	"budget_manager/internal/user"
	"budget_manager/internal/wallet"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	db, err := sqlite.ConnnectSQLite(
		"./assets/sqlite/database.sql",
		"./assets/sqlite/database.db",
	)
	if err != nil {
		slog.Error("init database", slog.String("error", err.Error()))
		return
	}
	slog.Info("connected to the database")
	defer db.Close()

	sm := session.NewSessionManager(
		session.NewRepository(db),
	)

	walletHandler := wallet.NewHandler(
		wallet.NewService(
			wallet.NewRepository(db),
		),
	)

	userHandler := user.NewHandler(
		user.NewService(
			user.NewRepository(db),
		), sm,
	)

	r := router.SetupRouter(walletHandler, userHandler, sm)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		slog.Info("Server is starting on port: " + srv.Addr + " ...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server stopped due error", slog.String("error", err.Error()))
			return
		}
	}()

	<-ctx.Done()
	cancel()
	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("reason", err.Error()))
		return
	}

	slog.Info("Server exited")
}
