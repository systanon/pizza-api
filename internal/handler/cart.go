package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"pizza-backend/internal/api"
	"pizza-backend/internal/apperror"
	"pizza-backend/internal/model"
	"pizza-backend/internal/response"
	"pizza-backend/internal/service"

	"github.com/gin-gonic/gin"
)

var uuidRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

type CartHandler struct {
	service *service.CartService
}

func NewCartHandler(s *service.CartService) *CartHandler {
	return &CartHandler{service: s}
}

func (h *CartHandler) CreateCart(c *gin.Context) {
	cart, err := h.service.CreateCart(c)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}
	response.RespondSuccess(c, http.StatusCreated, gin.H{"id": cart.ID}, "cart created")
}

func (h *CartHandler) GetCart(c *gin.Context) {
	cartID := c.Param("id")
	if !isValidCartID(cartID) {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "invalid cart id")
		return
	}

	cart, err := h.service.GetCart(c, cartID)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.RespondError(c, http.StatusNotFound, apperror.CodeNotFound, "cart not found")
			return
		}
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}

	response.RespondSuccess(c, http.StatusOK, cartDetailToAPI(cart), "")
}

func (h *CartHandler) AddItem(c *gin.Context) {
	cartID := c.Param("id")
	if !isValidCartID(cartID) {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "invalid cart id")
		return
	}

	var req api.CreateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, err.Error())
		return
	}

	_, err := h.service.AddItem(c, cartID, req.ProductID, req.VariantID, req.Quantity, req.AddonIDs)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.RespondError(c, http.StatusNotFound, apperror.CodeNotFound, "cart not found")
			return
		}
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}

	cart, err := h.service.GetCart(c, cartID)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}
	response.RespondSuccess(c, http.StatusCreated, cartDetailToAPI(cart), "item added")
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	cartID := c.Param("id")
	if !isValidCartID(cartID) {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "invalid cart id")
		return
	}
	itemID, err := strconv.ParseInt(c.Param("itemId"), 10, 64)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "invalid item id")
		return
	}

	var req api.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, err.Error())
		return
	}

	if err := h.service.UpdateItem(c, cartID, itemID, req.Quantity); err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.RespondError(c, http.StatusNotFound, apperror.CodeNotFound, "item not found")
			return
		}
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}

	cart, err := h.service.GetCart(c, cartID)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}
	response.RespondSuccess(c, http.StatusOK, cartDetailToAPI(cart), "item updated")
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	cartID := c.Param("id")
	if !isValidCartID(cartID) {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "invalid cart id")
		return
	}
	itemID, err := strconv.ParseInt(c.Param("itemId"), 10, 64)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "invalid item id")
		return
	}

	if err := h.service.RemoveItem(c, cartID, itemID); err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.RespondError(c, http.StatusNotFound, apperror.CodeNotFound, "item not found")
			return
		}
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}

	cart, err := h.service.GetCart(c, cartID)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}
	response.RespondSuccess(c, http.StatusOK, cartDetailToAPI(cart), "item removed")
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	cartID := c.Param("id")
	if !isValidCartID(cartID) {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "invalid cart id")
		return
	}

	if err := h.service.ClearCart(c, cartID); err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			response.RespondError(c, http.StatusNotFound, apperror.CodeNotFound, "cart not found")
			return
		}
		response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		return
	}

	response.RespondSuccess(c, http.StatusOK, nil, "cart cleared")
}

func isValidCartID(id string) bool {
	return uuidRe.MatchString(id)
}

func cartDetailToAPI(cart *model.CartDetail) api.CartResponse {
	items := make([]api.CartItem, 0, len(cart.Items))
	for _, it := range cart.Items {
		addons := make([]api.CartAddon, 0, len(it.Addons))
		for _, a := range it.Addons {
			addons = append(addons, api.CartAddon{
				ID:    a.AddonID,
				Name:  a.AddonName,
				Price: a.Price,
			})
		}
		items = append(items, api.CartItem{
			ID:           it.CartItem.ID,
			ProductID:    it.CartItem.ProductID,
			ProductName:  it.ProductName,
			ProductImage: it.ProductImage,
			VariantID:    it.CartItem.VariantID,
			VariantName:  it.VariantName,
			VariantUnit:  it.VariantUnit,
			VariantPrice: it.VariantPrice,
			Quantity:     it.CartItem.Quantity,
			Addons:       addons,
		})
	}
	return api.CartResponse{ID: cart.Cart.ID, Items: items}
}
