package handler

import (
	"net/http"
	"pizza-backend/internal/apperror"
	"pizza-backend/internal/response"
	"pizza-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

func (h *CategoryHandler) GetCategoryList(c *gin.Context) {
	categories, err := h.service.GetCategoryList(c)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}
	response.RespondSuccess(c, http.StatusOK, categories, "")
}
