package handler

import (
	"net/http"
	"pizza-backend/internal/api"
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

	result := make([]api.Category, 0, len(categories))

	for _, category := range categories {
		result = append(result, api.Category{
			ID:        category.ID,
			Name:      category.Name,
			Slug:      category.Slug,
			CreatedAt: category.CreatedAt,
		})
	}
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}
	response.RespondSuccess(c, http.StatusOK, result, "")
}
