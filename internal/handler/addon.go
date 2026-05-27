package handler

import (
	"net/http"
	"pizza-backend/internal/api"
	"pizza-backend/internal/apperror"
	"pizza-backend/internal/response"
	"pizza-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AddonHandler struct {
	service *service.AddonService
}

func NewAddonHandler(s *service.AddonService) *AddonHandler {
	return &AddonHandler{service: s}
}

func (h *AddonHandler) GetAddonListByCategoryID(c *gin.Context) {
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "invalid category id")
		return
	}

	addons, err := h.service.GetAddonListByCategoryID(c, categoryID)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}

	result := make([]api.Addon, 0, len(addons))
	for _, a := range addons {
		result = append(result, api.Addon{
			ID:    a.ID,
			Name:  a.Name,
			Price: a.Price,
		})
	}

	response.RespondSuccess(c, http.StatusOK, result, "")
}
