package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletService interface {
	CreateWallet(wo Wallet) (Wallet, error)
	AddOperation(userID int64, opt Operation) error
	ShowWallet(id int64) (Wallet, error)
}

type walletHandler struct {
	service WalletService
}

func NewWalletHandler(service WalletService) *walletHandler {
	return &walletHandler{service: service}
}

func (h *walletHandler) CreateWallet(ctx *gin.Context) {
	var wlt Wallet
	if err := ctx.BindJSON(&wlt); err != nil {
		badRequest(ctx, errors.New("wrong request"))
		return
	}
	w, err := h.service.CreateWallet(wlt)
	if err != nil {
		badRequest(ctx, errors.New("wrong request"))
		return
	}

	ctx.JSON(http.StatusOK, w)
}

func (h *walletHandler) AddOperation(ctx *gin.Context) {
	var opt OperationOptions
	if err := ctx.BindJSON(&opt); err != nil {
		badRequest(ctx, errors.New("wrong request"))
	}
	if err := h.service.AddOperation(opt.UserID, opt.Operation); err != nil {
		badRequest(ctx, err)
		return
	}

	ctx.String(http.StatusOK, "message: Success!")
}

func (h *walletHandler) ShowWallet(ctx *gin.Context) {
	var id gin.H
	if err := ctx.BindJSON(&id); err != nil {
		badRequest(ctx, errors.New("wrong request"))
	}

	userID, ok := id["user_id"].(float64)
	if !ok {
		badRequest(ctx, errors.New("wrong user ID"))
		return
	}

	w, err := h.service.ShowWallet(int64(userID))
	if err != nil {
		badRequest(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, w)
}

func badRequest(ctx *gin.Context, err error) {
	ctx.String(http.StatusBadRequest, "message: %s", err.Error())
}
