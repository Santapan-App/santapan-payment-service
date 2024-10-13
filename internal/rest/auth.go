package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"tobby/domain"
	"tobby/pkg/json"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

//go:generate mockery --name ArticleService
type TokenService interface {
	Update(ctx context.Context, ar *domain.Token) error
	Store(context.Context, *domain.Token) error
	Delete(ctx context.Context, id int64) error
	GetByUserID(ctx context.Context, id int64) (domain.Token, error)
}

type DeviceService interface {
	Store(context.Context, *domain.Device) error
	Update(ctx context.Context, device *domain.Device) error
	GetByUserID(ctx context.Context, id int64) (domain.Device, error)
}

type UserService interface {
	Store(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

// ArticleHandler  represent the httphandler for article
type AuthHandler struct {
	TokenService TokenService
	UserService  UserService
	Validator    *validator.Validate
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewAuthHandler(e *echo.Echo, tokenService TokenService, userService UserService) {
	validator := validator.New()
	// Register the custom validation function
	validator.RegisterValidation("date", validateDate)

	handler := &AuthHandler{
		TokenService: tokenService,
		UserService:  userService,
		Validator:    validator,
	}
	e.POST("/login", handler.Login)
	e.POST("/register", handler.Register)
}

var jwtSecret = []byte("SANTAPANSECRET") // Replace with your secret key

// Login handles user login and returns tokens
func (th *AuthHandler) Login(c echo.Context) (err error) {
	var loginBody domain.LoginBody
	ctx := c.Request().Context()

	if err = c.Bind(&loginBody); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid request body", nil)
	}

	// Validate the loginBody
	if err = th.Validator.Struct(loginBody); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Validation failed: "+err.Error(), nil)
	}

	user, err := th.UserService.GetByEmail(ctx, loginBody.Email)

	if err != nil || &user == nil {
		return json.Response(c, http.StatusBadRequest, false, "Email or password is incorrect", nil)
	}

	// Password Check Password Bcrypt
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginBody.Password)) != nil {
		return json.Response(c, http.StatusBadRequest, false, "Email or password is incorrect", nil)
	}

	accessToken, err := th.generateToken(user.ID, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to generate access token", nil)
	}

	refreshToken, err := th.generateToken(user.ID, time.Now().Add(time.Hour*24*90).Unix())
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to generate refresh token", nil)
	}

	// Check if a refresh token already exists
	existingToken, err := th.TokenService.GetByUserID(ctx, user.ID)

	if err != nil && err != domain.ErrNotFound {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to get refresh token", nil)
	}

	refreshData := domain.Token{
		RefreshToken: refreshToken,
		UserID:       user.ID,
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
	}

	if existingToken != (domain.Token{}) {
		// Token exists, update it
		refreshData.ID = existingToken.ID
		refreshData.CreatedAt = existingToken.CreatedAt
		if err = th.TokenService.Update(ctx, &refreshData); err != nil {
			logrus.Info(err)
			return json.Response(c, http.StatusInternalServerError, false, "Failed to update refresh token", nil)
		}
	} else {
		// Token does not exist, insert it
		if err = th.TokenService.Store(ctx, &refreshData); err != nil {
			return json.Response(c, http.StatusInternalServerError, false, fmt.Sprintf("Failed to store refresh token: %v", err), nil)
		}
	}

	userData := map[string]interface{}{
		"user":         user,
		"refreshToken": refreshToken,
		"accessToken":  accessToken,
	}

	return json.Response(c, http.StatusOK, true, "Login successfully!", userData)
}

func (th *AuthHandler) Register(c echo.Context) (err error) {
	var registerBody domain.RegisterBody

	ctx := c.Request().Context()

	if err = c.Bind(&registerBody); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid request body", nil)
	}

	checkUserEmail, err := th.UserService.GetByEmail(ctx, registerBody.Email)

	if checkUserEmail != (domain.User{}) {
		return json.Response(c, http.StatusBadRequest, false, "Email has already taken", nil)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(registerBody.Password), bcrypt.DefaultCost)
	// Insert User
	user := domain.User{
		FullName: registerBody.FullName,
		Email:    registerBody.Email,
		Password: string(hashPassword),
	}

	if err = th.UserService.Store(ctx, &user); err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to store user", nil)
	}

	accessToken, err := th.generateToken(user.ID, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to generate access token", nil)
	}

	refreshToken, err := th.generateToken(user.ID, time.Now().Add(time.Hour*24*90).Unix())
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to generate refresh token", nil)
	}

	// Check if a refresh token already exists
	existingToken, err := th.TokenService.GetByUserID(ctx, user.ID)

	if err != nil && err != domain.ErrNotFound {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to get refresh token", nil)
	}

	refreshData := domain.Token{
		RefreshToken: refreshToken,
		UserID:       user.ID,
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
	}

	if existingToken != (domain.Token{}) {
		// Token exists, update it
		refreshData.ID = existingToken.ID
		refreshData.CreatedAt = existingToken.CreatedAt
		if err = th.TokenService.Update(ctx, &refreshData); err != nil {
			return json.Response(c, http.StatusInternalServerError, false, "Failed to update refresh token", nil)
		}
	} else {
		// Token does not exist, insert it
		if err = th.TokenService.Store(ctx, &refreshData); err != nil {
			return json.Response(c, http.StatusInternalServerError, false, fmt.Sprintf("Failed to store refresh token: %v", err), nil)
		}
	}

	userData := map[string]interface{}{
		"user":         user,
		"refreshToken": refreshToken,
		"accessToken":  accessToken,
	}

	return json.Response(c, http.StatusOK, true, "Register successfully!", userData)
}

// Helper function to generate a new access token
func (th *AuthHandler) generateToken(userId int64, exp int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "auth-service",
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp": exp,
	})

	return token.SignedString(jwtSecret)
}

// Custom date validation function
func validateDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	layout := "2006-01-02" // Define the layout corresponding to the date format
	_, err := time.Parse(layout, dateStr)
	return err == nil
}
