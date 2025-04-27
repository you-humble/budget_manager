package user

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"budget_manager/internal/er"
	"budget_manager/internal/session"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

type Service interface {
	Save(u User) (User, error)
	FindByLogin(login string) (User, error)
}

type SessionManager interface {
	CreateSession(ctx *gin.Context, userID int64) error
	DestroySession(ctx *gin.Context) error
}

type handler struct {
	service Service
	sm      SessionManager
}

func NewHandler(service Service, sm SessionManager) *handler {
	return &handler{service: service, sm: sm}
}

func (h *handler) Register(ctx *gin.Context) {
	cr, err := credentials(ctx)
	if err != nil {
		return
	}

	u, err := h.service.Save(User{
		Login:    cr.Login,
		Password: hashPass(cr.Password, session.RandStringRunes(8)),
	})
	if err != nil {
		er.InternalServerError(ctx, "internal server error")
		return
	}

	if err := h.sm.CreateSession(ctx, u.ID); err != nil {
		er.InternalServerError(ctx, "internal server error")
		return
	}
	ctx.String(http.StatusOK, "message: Success!")
}

func (h *handler) Login(ctx *gin.Context) {
	cr, err := credentials(ctx)
	if err != nil {
		return
	}

	fmt.Println(cr)
	u, err := h.service.FindByLogin(cr.Login)
	if err != nil {
		er.InternalServerError(ctx, "internal server error")
		return
	}

	salt := string(u.Password[:8])
	if !bytes.Equal(hashPass(cr.Password, salt), u.Password) {
		er.BadRequest(ctx, "bad password")
		return
	}

	if err := h.sm.CreateSession(ctx, u.ID); err != nil {
		er.InternalServerError(ctx, "internal server error")
		return
	}
	ctx.String(http.StatusOK, "message: Success!")
}

func (h *handler) Logout(ctx *gin.Context) {
	if err := h.sm.DestroySession(ctx); err != nil {
		er.InternalServerError(ctx, "internal server error")
		return
	}
	ctx.String(http.StatusOK, "message: Success!")
}

func hashPass(plainPass, salt string) []byte {
	hashed := argon2.IDKey([]byte(plainPass), []byte(salt), 1, 64*1024, 4, 32)
	res := make([]byte, len(salt))
	copy(res, salt[:len(salt)])
	return append(res, hashed...)
}

func credentials(ctx *gin.Context) (Credentials, error) {
	cr := Credentials{}

	if err := ctx.BindJSON(&cr); err != nil {
		er.BadRequest(ctx, "bad request")
		return Credentials{}, err
	}

	if cr.Login == "" || cr.Password == "" {
		ctx.String(http.StatusForbidden, "wrong login or password")
		return Credentials{}, errors.New("wrong login or password")
	}

	return cr, nil
}
