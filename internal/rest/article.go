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
type ArticleService interface {
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	Fetch(ctx context.Context, cursor string, num int64) ([]domain.Article, string, error)
}

// ArticleHandler  represent the httphandler for article
type ArticleHandler struct {
	ArticleService ArticleService
	Validator      *validator.Validate
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewArticleHandler(e *echo.Echo, articleService ArticleService) {
	validator := validator.New()
	// Register the custom validation function
	validator.RegisterValidation("date", validateDate)

	handler := &ArticleHandler{
		ArticleService: articleService,
		Validator:      validator,
	}

	e.GET("/articles", handler.Fetch)
	e.GET("/articles/:id", handler.GetByID)
}

// GetByID will get article by given id
func (ah *ArticleHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	article, err := ah.ArticleService.GetByID(c.Request().Context(), int64(id))
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}
	return json.Response(c, http.StatusOK, true, "Successfully Get Article!", article)
}

func (ah *ArticleHandler) Fetch(c echo.Context) error {
	cursor := c.QueryParam("cursor")
	num := c.QueryParam("num")

	parseNum, err := strconv.Atoi(num)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid Num")
	}

	articles, nextCursor, err := ah.ArticleService.Fetch(c.Request().Context(), cursor, int64(parseNum))

	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	responseData := map[string]interface{}{
		"articles":   articles,
		"nextCursor": nextCursor,
	}

	return json.Response(c, http.StatusOK, true, "Successfully Get Articles!", responseData)
}
