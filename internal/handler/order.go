package handler

import (
	"errors"
	"net/http"
	"strconv"

	"pizza-backend/internal/api"
	"pizza-backend/internal/apperror"
	"pizza-backend/internal/model"
	"pizza-backend/internal/response"
	"pizza-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req api.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, err.Error())
		return
	}

	order, err := h.service.CreateOrder(c, req.CartID, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrNotFound):
			response.RespondError(c, http.StatusNotFound, apperror.CodeNotFound, "cart not found")
		case errors.Is(err, apperror.ErrCartEmpty):
			response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, err.Error())
		default:
			response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		}
		return
	}

	response.RespondSuccess(c, http.StatusCreated, orderDetailToAPI(order), "order created")
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "invalid order id")
		return
	}

	order, err := h.service.GetOrder(c, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.RespondError(c, http.StatusNotFound, apperror.CodeNotFound, "order not found")
			return
		}
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}

	response.RespondSuccess(c, http.StatusOK, orderDetailToAPI(order), "")
}

func orderDetailToAPI(o *model.OrderDetail) api.OrderResponse {
	items := make([]api.OrderItem, 0, len(o.Items))
	for _, it := range o.Items {
		addons := make([]api.OrderAddon, 0, len(it.Addons))
		for _, a := range it.Addons {
			addons = append(addons, api.OrderAddon{
				ID:    a.AddonID,
				Name:  a.AddonName,
				Price: a.Price,
			})
		}
		items = append(items, api.OrderItem{
			ID:           it.OrderItem.ID,
			ProductID:    it.ProductID,
			ProductName:  it.ProductName,
			VariantID:    it.VariantID,
			VariantName:  it.VariantName,
			VariantPrice: it.VariantPrice,
			Quantity:     it.Quantity,
			ItemTotal:    it.ItemTotal,
			Addons:       addons,
		})
	}
	return api.OrderResponse{
		ID:     o.Order.ID,
		CartID: o.CartID,
		Email:  o.Email,
		Status: string(o.Status),
		Total:  o.Total,
		Items:  items,
	}
}
