package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"pizza-backend/internal/apperror"
	"pizza-backend/internal/response"

	"github.com/gin-gonic/gin"
)

type StripeServiceInterface interface {
	CreateCheckoutSession(ctx context.Context, orderID int64, successURL, cancelURL string) (string, error)
	HandleWebhook(ctx context.Context, payload []byte, signature string) error
}

type StripeHandler struct {
	service    StripeServiceInterface
	successURL string
	cancelURL  string
}

func NewStripeHandler(s StripeServiceInterface, successURL, cancelURL string) *StripeHandler {
	return &StripeHandler{service: s, successURL: successURL, cancelURL: cancelURL}
}

func (h *StripeHandler) CreateCheckout(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "invalid order id")
		return
	}

	url, err := h.service.CreateCheckoutSession(c, id, h.successURL, h.cancelURL)
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrNotFound):
			response.RespondError(c, http.StatusNotFound, apperror.CodeNotFound, "order not found")
		case errors.Is(err, apperror.ErrOrderNotPending):
			response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, err.Error())
		default:
			response.RespondError(c, http.StatusInternalServerError, apperror.CodeInternal, err.Error())
		}
		return
	}

	response.RespondSuccess(c, http.StatusOK, gin.H{"checkout_url": url}, "")
}

func (h *StripeHandler) Webhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "failed to read body")
		return
	}

	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, "missing Stripe-Signature header")
		return
	}

	if err := h.service.HandleWebhook(c, payload, signature); err != nil {
		response.RespondError(c, http.StatusBadRequest, apperror.CodeValidation, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
