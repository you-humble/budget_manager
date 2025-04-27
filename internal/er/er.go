package er

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InternalServerError(ctx *gin.Context, msg string) {
	ctx.String(http.StatusInternalServerError, "message: %s", msg)
}

func BadRequest(ctx *gin.Context, msg string) {
	ctx.String(http.StatusBadRequest, "message: %s", msg)
}
