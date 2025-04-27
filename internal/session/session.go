package session

import (
	"errors"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	sessionKey    string = "user"
	sessionCookie string = "session_id"
)

var (
	ErrNoAuth   = errors.New("No session found")
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type SessionRepository interface {
	FindByID(id string) (*Session, error)
	Save(sessID string, userID int64) error
	Delete(sessID string) error
}

type Session struct {
	ID     string `json:"id" db:"id"`
	UserID int64  `json:"user_id" db:"user_id"`
}

type sessionManager struct {
	repo SessionRepository
}

func NewSessionManager(repo SessionRepository) *sessionManager {
	return &sessionManager{repo: repo}
}

func (sm *sessionManager) CreateSession(ctx *gin.Context, userID int64) error {
	sessID := RandStringRunes(32)
	if err := sm.repo.Save(sessID, userID); err != nil {
		slog.Error("failed to create session", slog.String("error", err.Error()))
		return err
	}

	ctx.SetCookie(
		sessionCookie,
		sessID,
		90*24*60,
		"/", "localhost", false, true,
	)

	return nil
}

func (sm *sessionManager) DestroySession(ctx *gin.Context) error {
	sess, err := sm.SessionFromContext(ctx)
	if err == nil {
		if err := sm.repo.Delete(sess.ID); err != nil {
			return err
		}
	}

	ctx.SetCookie(
		sessionCookie,
		"",
		-1,
		"/", "localhost", false, true,
	)

	return nil
}

func (sm *sessionManager) SessionFromContext(ctx *gin.Context) (*Session, error) {
	sess, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil, ErrNoAuth
	}

	return sess, nil
}

func (sm *sessionManager) CheckSession(ctx *gin.Context) (*Session, error) {
	sessID, err := ctx.Cookie(sessionCookie)
	if err == http.ErrNoCookie {
		slog.Warn("no cookie")
		return nil, err
	}

	sess, err := sm.repo.FindByID(sessID)
	if err != nil {
		slog.Error("failed to find session", slog.String("error", err.Error()))
		return nil, err
	}

	return sess, nil
}

var (
	noAuthUrls = map[string]struct{}{
		"/user/login":    struct{}{},
		"/user/register": struct{}{},
	}
)

func (sm *sessionManager) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ok := noAuthUrls[ctx.Request.URL.Path]; ok {
			ctx.Next()
			return
		}

		sess, err := sm.CheckSession(ctx)
		if err != nil {
			ctx.String(http.StatusUnauthorized, "no auth")
			ctx.Abort()
			return
		}

		ctx.Set(sessionKey, sess)
		ctx.Next()
	}
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
