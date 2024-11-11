package rest

import (
	"context"
	"net/http"
	"santapan/domain"
	"santapan/pkg/json"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

// TransactionService defines the methods for handling transaction-related logic
type TransactionService interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]any, string, error)
	GetCart(ctx context.Context) ([]domain.CartItem, error)
	AddToCart(ctx context.Context, item domain.CartItem) error
	DoTransaction(ctx context.Context) error
}

// TransactionHandler represents the HTTP handler for transactions
type TransactionHandler struct {
	TransactionService TransactionService
	Validator          *validator.Validate
}

// NewTransactionHandler initializes the transaction resources endpoint
func NewTransactionHandler(e *echo.Echo, transactionService TransactionService) {
	validator := validator.New()
	handler := &TransactionHandler{
		TransactionService: transactionService,
		Validator:          validator,
	}

	e.GET("/histories", handler.Fetch)
	e.GET("/cart", handler.GetCart)
	e.POST("/cart", handler.AddToCart)
	e.POST("/transaction", handler.DoTransaction)
}

// Fetch retrieves transaction histories
func (th *TransactionHandler) Fetch(c echo.Context) error {
	cursor := c.QueryParam("cursor")
	num := c.QueryParam("num")

	parseNum, err := strconv.Atoi(num)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid Num")
	}

	transactions, nextCursor, err := th.TransactionService.Fetch(c.Request().Context(), cursor, int64(parseNum))
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	responseData := map[string]interface{}{
		"transactions": transactions,
		"nextCursor":   nextCursor,
	}

	return json.Response(c, http.StatusOK, true, "Successfully retrieved transaction histories!", responseData)
}

// GetCart retrieves the current cart items
func (th *TransactionHandler) GetCart(c echo.Context) error {
	cartItems, err := th.TransactionService.GetCart(c.Request().Context())
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	return json.Response(c, http.StatusOK, true, "Successfully retrieved cart items!", cartItems)
}

// AddToCart adds an item to the cart
func (th *TransactionHandler) AddToCart(c echo.Context) error {
	var item domain.CartItem
	if err := c.Bind(&item); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid item data")
	}

	if err := th.TransactionService.AddToCart(c.Request().Context(), item); err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	return json.Response(c, http.StatusOK, true, "Item successfully added to cart!", nil)
}

// DoTransaction finalizes the cart into a transaction
func (th *TransactionHandler) DoTransaction(c echo.Context) error {
	if err := th.TransactionService.DoTransaction(c.Request().Context()); err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	return json.Response(c, http.StatusOK, true, "Transaction completed successfully!", nil)
}
