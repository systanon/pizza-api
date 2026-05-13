package handler

import (
	"net/http"
	"pizza-backend/internal/apperror"
	"pizza-backend/internal/response"

	"github.com/gin-gonic/gin"
)

func NotFoundHandler(c *gin.Context) {
	response.RespondError(c, http.StatusBadRequest, apperror.CodeRouteNotFound, "route not found")
}
