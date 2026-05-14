package handler

import (
	"errors"
	"net/http"
	"pizza-backend/internal/api"
	"pizza-backend/internal/apperror"
	"pizza-backend/internal/model"
	"pizza-backend/internal/response"
	"pizza-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) GetProductList(c *gin.Context) {
	var query api.GetProductQueryList

	if err := c.ShouldBindQuery(&query); err != nil {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, err.Error())
		return
	}

	params := model.ProductParams{
		CategoryID: query.CategoryID,
		Offset:     query.Offset,
		Limit:      query.Limit,
		Page:       query.Page,
		PerPage:    query.PerPage,
		Q:          query.Q,
		SortOrder:  query.SortOrder,
	}

	products, total, pages, err := h.service.GetProductList(c, params)

	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}

	var result = make([]api.Product, 0, len(products))
	for _, product := range products {

		result = append(result, api.Product{
			ID:          product.ID,
			CategoryID:  product.CategoryID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ImageURL:    product.ImageURL,
			CreatedAt:   product.CreatedAt,
		})
	}
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.Header("X-Total-Pages", strconv.Itoa(pages))
	response.RespondSuccess(c, http.StatusOK, result, "product list retrieved successfully")
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "invalid id")
		return
	}

	product, err := h.service.GetProductByID(c, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.RespondError(c, http.StatusNotFound, apperror.CodeNotFound, "product not found")
			return
		}
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}

	response.RespondSuccess(c, http.StatusOK, product, "")
}
