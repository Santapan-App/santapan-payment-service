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

//go:generate mockery --name ArticleService
type CategoryService interface {
	GetByID(ctx context.Context, id int64) (domain.Category, error)
	Fetch(ctx context.Context, cursor string, num int64) ([]domain.Category, string, error)
}

// ArticleHandler  represent the httphandler for article
type CategoryHandler struct {
	CategoryService CategoryService
	Validator       *validator.Validate
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewCategoryHandler(e *echo.Echo, categoryService CategoryService) {
	validator := validator.New()
	// Register the custom validation function
	validator.RegisterValidation("date", validateDate)

	handler := &CategoryHandler{
		CategoryService: categoryService,
		Validator:       validator,
	}

	e.GET("/categories", handler.Fetch)
	e.GET("/categories/:id", handler.GetByID)
}

// GetByID will get article by given id
func (ah *CategoryHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	categories, err := ah.CategoryService.GetByID(c.Request().Context(), int64(id))
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}
	return json.Response(c, http.StatusOK, true, "Successfully Get Article!", categories)
}

func (ah *CategoryHandler) Fetch(c echo.Context) error {
	cursor := c.QueryParam("cursor")
	num := c.QueryParam("num")

	parseNum, err := strconv.Atoi(num)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid Num")
	}

	categories, nextCursor, err := ah.CategoryService.Fetch(c.Request().Context(), cursor, int64(parseNum))

	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	responseData := map[string]interface{}{
		"categories": categories,
		"nextCursor": nextCursor,
	}

	return json.Response(c, http.StatusOK, true, "Successfully Get Articles!", responseData)
}
