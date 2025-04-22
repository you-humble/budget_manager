package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type WalletHandler interface {
	CreateWallet(ctx *gin.Context)
	ShowWallet(ctx *gin.Context)
	AddOperation(ctx *gin.Context)
}

func SetupRouter(h WalletHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/wallet/create", h.CreateWallet)
	router.GET("/wallet/show", h.ShowWallet)
	router.POST("/operation/add", h.AddOperation)

	return router
}

func main() {
	db, err := ConnnectSQLite()
	if err != nil {
		slog.Error("init database", slog.String("error", err.Error()))
		return
	}
	defer db.Close()
	repo := NewWalletRepository(db)
	s := NewWalletService(repo)
	h := NewWalletHandler(s)
	r := SetupRouter(h)

	r.Run(":8080")
}
