package rest

import (
	"context"
	"net/http"
	"santapan_payment_service/domain"
	"santapan_payment_service/internal/rest/middleware"
	"santapan_payment_service/pkg/json"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

// PaymentService represents the payment service interface
type PaymentService interface {
	Store(ctx context.Context, payment *domain.Payment) error
	Update(ctx context.Context, payment *domain.Payment) error
	GetByID(ctx context.Context, id int64) (domain.Payment, error)
	GetByRefID(ctx context.Context, refID string) (domain.Payment, error)
}

// PaymentHandler represents the HTTP handler for payments
type PaymentHandler struct {
	PaymentService PaymentService
	Validator      *validator.Validate
}

// NewPaymentHandler initializes the payment endpoints
func NewPaymentHandler(e *echo.Echo, paymentService PaymentService) {
	validator := validator.New()

	handler := &PaymentHandler{
		PaymentService: paymentService,
		Validator:      validator,
	}

	e.POST("/payment", handler.ProcessPayment, middleware.AuthMiddleware)
	e.GET("/payments/:id", handler.GetPaymentByID, middleware.AuthMiddleware)
	e.POST("/payments/callback", handler.CallbackPayment)
}

// ProcessPayment processes a new payment
func (h *PaymentHandler) ProcessPayment(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var paymentBody domain.PaymentBody
	if err := c.Bind(&paymentBody); err != nil {
		return json.Response(c, http.StatusUnprocessableEntity, false, "Invalid request", nil)
	}

	if err := h.Validator.Struct(paymentBody); err != nil {
		return json.Response(c, http.StatusBadRequest, false, err.Error(), nil)
	}

	userID, ok := c.Get("userID").(int64)
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	// Generate Random ID
	referenceID := "TRX" + strconv.FormatInt(time.Now().Unix(), 10)

	payment := &domain.Payment{
		UserID:      userID,
		Amount:      paymentBody.Amount,
		Status:      "pending",
		ReferenceID: referenceID,
		Url:         "https://google.com/" + referenceID,
		SessionID:   "session-" + referenceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.PaymentService.Store(ctx, payment); err != nil {
		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return json.Response(c, http.StatusCreated, true, "Payment processed successfully", payment)
}

// GetPaymentByID retrieves a specific payment by its ID
func (h *PaymentHandler) GetPaymentByID(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	paymentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid payment ID", nil)
	}

	userID, ok := c.Get("userID").(int64)
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	payment, err := h.PaymentService.GetByID(ctx, paymentID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	if payment.UserID != userID {
		return json.Response(c, http.StatusForbidden, false, "You are not authorized to access this payment", nil)
	}

	return json.Response(c, http.StatusOK, true, "Success", payment)
}

// CallbackPayment handles the payment callback from the payment gateway and update the payment status
func (h *PaymentHandler) CallbackPayment(c echo.Context) error {

	// Parse the
	var iPaymuCallback domain.IPaymuCallback
	if err := c.Bind(&iPaymuCallback); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid request", nil)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// Get the payment details based on the reference ID
	payment, err := h.PaymentService.GetByRefID(ctx, iPaymuCallback.ReferenceID)

	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	// Update the payment status based on the callback status
	if iPaymuCallback.StatusCode == "0" {
		payment.Status = "pending"
	} else if iPaymuCallback.StatusCode == "1" {
		payment.Status = "success"
	} else {
		payment.Status = "expired"
	}

	// Update the payment status
	if err := h.PaymentService.Update(ctx, &payment); err != nil {
		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return json.Response(c, http.StatusOK, true, "Payment status updated successfully", payment)
}
