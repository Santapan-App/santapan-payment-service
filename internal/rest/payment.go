package rest

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	encodingJson "encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"santapan_payment_service/domain"
	"santapan_payment_service/internal/rest/middleware"
	"santapan_payment_service/pkg/json"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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
	ctx := c.Request().Context() // Echo ensures this is never nil

	var paymentBody domain.PaymentBody
	if err := c.Bind(&paymentBody); err != nil {
		logrus.Error("Invalid request body: ", err)
		return json.Response(c, http.StatusUnprocessableEntity, false, "Invalid request", nil)
	}

	if err := h.Validator.Struct(paymentBody); err != nil {
		logrus.Error("Validation error: ", err)
		return json.Response(c, http.StatusBadRequest, false, err.Error(), nil)
	}

	userID, ok := c.Get("userID").(int64)
	if !ok {
		logrus.Error("Unauthorized access: missing userID")
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	// iPaymu credentials
	ipaymuVa := os.Getenv("IPAYMU_VA")
	ipaymuKey := os.Getenv("IPAYMU_KEY")

	logrus.Info("IPAYMU VA: ", ipaymuVa)
	logrus.Info("IPAYMU KEY: ", ipaymuKey)

	if ipaymuVa == "" || ipaymuKey == "" {
		logrus.Error("Missing IPAYMU credentials")
		return json.Response(c, http.StatusInternalServerError, false, "Payment service not configured", nil)
	}

	// Prepare the request data
	url, _ := url.Parse("https://sandbox.ipaymu.com/api/v2/payment")
	referenceID := "TRX" + strconv.FormatInt(time.Now().Unix(), 10)

	requestData := map[string]interface{}{
		"product":     paymentBody.Name,
		"qty":         paymentBody.Qty,
		"price":       paymentBody.Price,
		"imageUrl":    paymentBody.ImageUrl,
		"notifyUrl":   "http://payment.santapan.store/payments/callback",
		"referenceId": referenceID,
	}

	// Marshal the request body
	postBody, err := encodingJson.Marshal(requestData)
	if err != nil {
		logrus.Error("Failed to encode JSON: ", err)
		return json.Response(c, http.StatusInternalServerError, false, "Failed to process payment", nil)
	}

	// Generate the signature
	bodyHash := sha256.Sum256(postBody)
	bodyHashToString := hex.EncodeToString(bodyHash[:])
	stringToSign := "POST:" + ipaymuVa + ":" + strings.ToLower(bodyHashToString) + ":" + ipaymuKey

	hmacHash := hmac.New(sha256.New, []byte(ipaymuKey))
	hmacHash.Write([]byte(stringToSign))
	signature := hex.EncodeToString(hmacHash.Sum(nil))

	// Prepare the HTTP request body
	reqBody := ioutil.NopCloser(strings.NewReader(string(postBody)))

	// Create the HTTP request using struct initialization
	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
			"va":           {ipaymuVa},
			"signature":    {signature},
		},
		Body: reqBody,
	}

	logrus.Info("Request Header: ", req.Header)
	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error("Error processing payment: ", err)
		return json.Response(c, http.StatusInternalServerError, false, "Failed to reach payment gateway", nil)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Failed to read response: ", err)
		return json.Response(c, http.StatusInternalServerError, false, "Failed to process payment response", nil)
	}

	logrus.Info("Response: ", string(body))

	// Parse response
	var ipaymuResponse domain.IPaymuResponse
	if err := encodingJson.Unmarshal(body, &ipaymuResponse); err != nil {
		logrus.Error("Failed to decode response JSON: ", err)
		return json.Response(c, http.StatusInternalServerError, false, "Invalid response from payment gateway", nil)
	}

	// Validate response
	if ipaymuResponse.Status != 200 {
		logrus.Error("Payment failed: ", ipaymuResponse.Message)
		return json.Response(c, http.StatusInternalServerError, false, ipaymuResponse.Message, nil)
	}

	if ipaymuResponse.Data == nil {
		logrus.Error("Invalid response: Missing Data field")
		return json.Response(c, http.StatusInternalServerError, false, "Invalid payment response", nil)
	}

	// Store payment
	payment := &domain.Payment{
		UserID:      userID,
		Amount:      paymentBody.Amount,
		Status:      "pending",
		ReferenceID: referenceID,
		Url:         ipaymuResponse.Data.Url,
		SessionID:   ipaymuResponse.Data.SessionID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.PaymentService.Store(ctx, payment); err != nil {
		logrus.Error("Failed to store payment: ", err)
		return json.Response(c, http.StatusInternalServerError, false, "Payment processing error", nil)
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
