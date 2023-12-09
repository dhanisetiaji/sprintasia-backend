package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"gotham/config"
	"gotham/models"
	"gotham/requests"
	"gotham/services"
	"gotham/viewModels"
)

type AuthController struct {
	AuthService services.IAuthService
}

// Login godoc
// @Summary
// @Description
// @Tags Auth
// @Accept  json
// @Produce json
// @Param reqBody body LoginRequest true "Login" "email,password,platform"
// @Success 200 {object} viewModels.HTTPSuccessResponse{data=viewModels.Login}
// @Failure 422 {object} viewModels.HTTPErrorResponse{}
// @Failure 400 {object} viewModels.Message{}
// @Failure 500 {object} viewModels.Message{}
// @Router /v1/login [post]
func (a AuthController) Login(c echo.Context) (err error) {
	// Request Bind And Validation
	request := new(requests.LoginRequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, &request.Body); err != nil {
		return err
	}
	v := request.Validate()
	if v != nil {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(v))
	}
	fmt.Println(request.Body.Email, request)

	var user models.User
	user, err = a.AuthService.GetUserByEmail(request.Body.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
				"message": "email or password is incorrect",
			}))
		} else {
			return echo.ErrInternalServerError
		}
	}

	var verify bool
	verify, err = a.AuthService.Check(request.Body.Email, request.Body.Password)
	if !verify {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"message": "email or password is incorrect",
		}))
	}

	accessTokenExp := time.Now().Add(time.Hour * 720).Unix()

	claims := &config.JwtCustomClaims{
		AuthID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var accessToken string
	accessToken, err = token.SignedString([]byte(config.Conf.SecretKey))
	if err != nil {
		return
	}

	// Response
	return c.JSON(http.StatusOK, viewModels.SuccessResponse(viewModels.Login{
		AccessToken:    accessToken,
		AccessTokenExp: accessTokenExp,
		User:           user,
	}))
}

// Register godoc
// @Summary
// @Description
// @Tags Auth
// @Accept  json
// @Produce json
// @Param reqBody body models.RegisterUser true "Register"
// @Success 200 {object} viewModels.HTTPSuccessResponse{}
// @Failure 422 {object} viewModels.HTTPErrorResponse{}
// @Failure 400 {object} viewModels.Message{}
// @Failure 500 {object} viewModels.Message{}
// @Router /v1/register [post]
func (a AuthController) Register(c echo.Context) (err error) {
	// Request Bind And Validation
	var request = models.RegisterUser{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &request); err != nil {
		return err
	}
	user, _ := a.AuthService.GetUserByEmail(request.Email)
	if user != (models.User{}) {
		return c.JSON(http.StatusBadRequest, viewModels.ValidationResponse(map[string]string{
			"message": "email already exists",
		}))
	}

	err = a.AuthService.Register(request)
	if err != nil {
		return echo.ErrBadRequest
	}

	// Response
	return c.JSON(http.StatusOK, viewModels.SuccessResponse(map[string]string{
		"message": "user created successfully",
	}))
}
