package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type WalletHandler interface {
	CreateWallet(ctx *gin.Context)
	ShowWallet(ctx *gin.Context)
	AddOperation(ctx *gin.Context)
}

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type SessionManager interface {
	AuthMiddleware() gin.HandlerFunc
}

func SetupRouter(wh WalletHandler, uh UserHandler, sm SessionManager) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(sm.AuthMiddleware())

	wallet := router.Group("/wallet")
	wallet.POST("/create", wh.CreateWallet)
	wallet.GET("/show", wh.ShowWallet)
	wallet.POST("/operation/add", wh.AddOperation)

	users := router.Group("/user")
	users.POST("/register", uh.Register)
	users.POST("/login", uh.Login)
	users.DELETE("/logout", uh.Logout)

	slog.Info("the router has been set up")
	return router
}
