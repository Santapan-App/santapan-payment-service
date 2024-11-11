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
type BannerService interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]domain.Banner, string, error)
}

// ArticleHandler  represent the httphandler for article
type BannerHandler struct {
	BannerService BannerService
	Validator     *validator.Validate
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewBannerHandler(e *echo.Echo, bannerService BannerService) {
	validator := validator.New()
	// Register the custom validation function
	validator.RegisterValidation("date", validateDate)

	handler := &BannerHandler{
		BannerService: bannerService,
		Validator:     validator,
	}

	e.GET("/banners", handler.Fetch)
}

func (ah *BannerHandler) Fetch(c echo.Context) error {
	cursor := c.QueryParam("cursor")
	num := c.QueryParam("num")

	parseNum, err := strconv.Atoi(num)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid Num")
	}

	banners, nextCursor, err := ah.BannerService.Fetch(c.Request().Context(), cursor, int64(parseNum))

	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	responseData := map[string]interface{}{
		"banners":    banners,
		"nextCursor": nextCursor,
	}

	return json.Response(c, http.StatusOK, true, "Successfully Get Articles!", responseData)
}
