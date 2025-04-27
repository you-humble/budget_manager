package wallet

import (
	"net/http"

	"budget_manager/internal/er"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Save(wo Wallet) (Wallet, error)
	AddOperation(userID int64, opt Operation) error
	ShowWallet(id int64) (Wallet, error)
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateWallet(ctx *gin.Context) {
	var wlt Wallet
	if err := ctx.BindJSON(&wlt); err != nil {
		er.BadRequest(ctx, "bad request")
		return
	}
	w, err := h.service.Save(wlt)
	if err != nil {
		er.BadRequest(ctx, "bad request")
		return
	}

	ctx.JSON(http.StatusOK, w)
}

func (h *handler) AddOperation(ctx *gin.Context) {
	var opt OperationOptions
	if err := ctx.BindJSON(&opt); err != nil {
		er.BadRequest(ctx, "bad request")
	}
	if err := h.service.AddOperation(opt.UserID, opt.Operation); err != nil {
		er.BadRequest(ctx, "bad request")
		return
	}

	ctx.String(http.StatusOK, "message: Success!")
}

func (h *handler) ShowWallet(ctx *gin.Context) {
	var id gin.H
	if err := ctx.BindJSON(&id); err != nil {
		er.BadRequest(ctx, "bad request")
	}

	userID, ok := id["user_id"].(float64)
	if !ok {
		er.BadRequest(ctx, "bad request")
		return
	}

	w, err := h.service.ShowWallet(int64(userID))
	if err != nil {
		er.BadRequest(ctx, "bad request")
		return
	}

	ctx.JSON(http.StatusOK, w)
}
