package controller

import (
	"time"

	"git.dev.siap.id/kukuhkkh/app-music/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/auth/request"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/auth/service"
	"git.dev.siap.id/kukuhkkh/app-music/utils/config"
	"git.dev.siap.id/kukuhkkh/app-music/utils/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type authController struct {
	authService service.AuthService
	cfg         *config.Config
}

type AuthController interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Me(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

func NewAuthController(authService service.AuthService, cfg *config.Config) AuthController {
	return &authController{
		authService: authService,
		cfg:         cfg,
	}
}

// do login
// @Summary      Do login
// @Description  API for do login
// @Tags         Authentication
// @Security     Bearer
// @Body 	     request.LoginRequest
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      422  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/login [post]
func (_i *authController) Login(c *fiber.Ctx) error {
	req := new(request.LoginRequest)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	res, err := _i.authService.Login(*req)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     _i.cfg.Cookie.Name,
		Value:    res.Token,
		Expires:  time.Unix(res.ExpiresAt, 0),
		HTTPOnly: _i.cfg.Cookie.HTTPOnly,
		Secure:   _i.cfg.Cookie.Secure,
		SameSite: _i.cfg.Cookie.SameSite,
	})

	return response.Resp(c, response.Response{
		Data:     res,
		Messages: response.Messages{"Login success"},
		Code:     fiber.StatusOK,
	})
}

// Register
// @Summary      Register
// @Description  API for register
// @Tags         Authentication
// @Body 	     request.RegisterRequest
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      422  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/register [post]
func (_i *authController) Register(c *fiber.Ctx) error {
	req := new(request.RegisterRequest)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	res, err := _i.authService.Register(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data:     res,
		Messages: response.Messages{"Register success"},
		Code:     fiber.StatusOK,
	})
}

// Me
// @Summary      Me
// @Description  API for get me
// @Tags         Authentication
// @Security     Bearer
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/auth/me [get]
func (_i *authController) Me(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(*middleware.JWTClaims)

	res, err := _i.authService.Me(claims.UserID)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data:     res,
		Messages: response.Messages{"Get me success"},
		Code:     fiber.StatusOK,
	})
}

// Logout
// @Summary      Logout
// @Description  API for logout
// @Tags         Authentication
// @Success      200  {object}  response.Response
// @Router       /api/v1/auth/logout [post]
func (_i *authController) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     _i.cfg.Cookie.Name,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: _i.cfg.Cookie.HTTPOnly,
		Secure:   _i.cfg.Cookie.Secure,
		SameSite: _i.cfg.Cookie.SameSite,
	})

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Logout success"},
		Code:     fiber.StatusOK,
	})
}
