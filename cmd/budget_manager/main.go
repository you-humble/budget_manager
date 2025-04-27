package main

import (
	"log/slog"

	"budget_manager/internal/database/sqlite"
	"budget_manager/internal/router"
	"budget_manager/internal/session"
	"budget_manager/internal/user"
	"budget_manager/internal/wallet"
)

func main() {
	db, err := sqlite.ConnnectSQLite()
	if err != nil {
		slog.Error("init database", slog.String("error", err.Error()))
		return
	}
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

	r.Run(":8080")
}
