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

type BundlingService interface {
	CreateBundling(ctx context.Context, bundling domain.Bundling) (domain.Bundling, error)
	GetBundling(ctx context.Context, id int64) (domain.Bundling, error)
	AddMenuToBundling(ctx context.Context, bundlingMenu domain.BundlingMenu) (domain.BundlingMenu, error)
	GetMenusInBundling(ctx context.Context, bundlingID int64) ([]domain.Menu, error)
}

// BundlingHandler manages HTTP endpoints for bundling
type BundlingHandler struct {
	BundlingService BundlingService
	Validator       *validator.Validate
}

// NewBundlingHandler initializes bundling routes
func NewBundlingHandler(e *echo.Echo, bundlingService BundlingService) {
	handler := &BundlingHandler{
		BundlingService: bundlingService,
		Validator:       validator.New(),
	}

	e.POST("/bundling", handler.CreateBundling)
	e.GET("/bundling/:id", handler.GetBundling)
	e.POST("/bundling/:id/menu", handler.AddMenuToBundling)
	e.GET("/bundling/:id/menus", handler.GetMenusInBundling)
}

// CreateBundling handles creating a new bundling
func (bh *BundlingHandler) CreateBundling(c echo.Context) error {
	var bundling domain.Bundling
	if err := c.Bind(&bundling); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid input")
	}

	bundling, err := bh.BundlingService.CreateBundling(c.Request().Context(), bundling)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	return json.Response(c, http.StatusCreated, true, "Bundling created successfully", bundling)
}

// GetBundling handles fetching a bundling by ID
func (bh *BundlingHandler) GetBundling(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid bundling ID")
	}

	bundling, err := bh.BundlingService.GetBundling(c.Request().Context(), id)
	if err != nil {
		return json.Response(c, http.StatusNotFound, false, "", "Bundling not found")
	}

	return json.Response(c, http.StatusOK, true, "Bundling fetched successfully", bundling)
}

// AddMenuToBundling adds a menu item to a specific bundling
func (bh *BundlingHandler) AddMenuToBundling(c echo.Context) error {
	bundlingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid bundling ID")
	}

	var bundlingMenu domain.BundlingMenu
	if err := c.Bind(&bundlingMenu); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid input")
	}
	bundlingMenu.BundlingID = bundlingID

	bundlingMenu, err = bh.BundlingService.AddMenuToBundling(c.Request().Context(), bundlingMenu)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	return json.Response(c, http.StatusCreated, true, "Menu added to bundling successfully", bundlingMenu)
}

// GetMenusInBundling retrieves all menu items in a bundling
func (bh *BundlingHandler) GetMenusInBundling(c echo.Context) error {
	bundlingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid bundling ID")
	}

	menus, err := bh.BundlingService.GetMenusInBundling(c.Request().Context(), bundlingID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	return json.Response(c, http.StatusOK, true, "Menus fetched successfully", menus)
}
