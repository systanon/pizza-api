package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"pizza-backend/internal/apperror"
	"pizza-backend/internal/model"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/webhook"
)

type StripeOrderRepository interface {
	GetOrder(ctx context.Context, orderID int64) (*model.OrderDetail, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status model.OrderStatus) error
}

type StripeService struct {
	repo          StripeOrderRepository
	webhookSecret string
}

func NewStripeService(repo StripeOrderRepository, secretKey, webhookSecret string) *StripeService {
	stripe.Key = secretKey
	return &StripeService{repo: repo, webhookSecret: webhookSecret}
}

func (s *StripeService) CreateCheckoutSession(ctx context.Context, orderID int64, successURL, cancelURL string) (string, error) {
	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return "", err
	}

	if order.Status != model.OrderStatusPending {
		return "", apperror.ErrOrderNotPending
	}

	lineItems := buildLineItems(order.Items)

	params := &stripe.CheckoutSessionParams{
		CustomerEmail:      stripe.String(order.Email),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(successURL),
		CancelURL:          stripe.String(cancelURL),
		Metadata: map[string]string{
			"order_id": fmt.Sprintf("%d", orderID),
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return "", err
	}

	return sess.URL, nil
}

func (s *StripeService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	event, err := webhook.ConstructEventWithOptions(payload, signature, s.webhookSecret,
		webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true})
	if err != nil {
		return err
	}

	if event.Type != "checkout.session.completed" {
		return nil
	}

	var sess stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
		return err
	}

	orderID, err := parseOrderID(sess.Metadata["order_id"])
	if err != nil {
		return err
	}

	return s.repo.UpdateOrderStatus(ctx, orderID, model.OrderStatusPaid)
}

func buildLineItems(items []model.OrderItemDetail) []*stripe.CheckoutSessionLineItemParams {
	lineItems := make([]*stripe.CheckoutSessionLineItemParams, 0, len(items))
	for _, it := range items {
		name := it.ProductName
		if it.VariantName != "" {
			name = fmt.Sprintf("%s (%s)", it.ProductName, it.VariantName)
		}
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("pln"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(name),
				},
				UnitAmount: stripe.Int64(it.ItemTotal / int64(it.Quantity)),
			},
			Quantity: stripe.Int64(int64(it.Quantity)),
		})
	}
	return lineItems
}

func parseOrderID(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
