package api

import (
	"net/http"
	"time"

	_ "echo-server/docs"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

type Token struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// @Summary Login endpoint
// @Description Logs in a user and returns a JWT token
// @ID login
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login Request"
// @Success 200 {object} Token "Token"
// @Failure 401 {string} string "Unauthorized"
// @Router /login [post]
func login(c echo.Context) error {

	var loginRequest LoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		return err
	}

	username := loginRequest.Username
	password := loginRequest.Password
	if password != "password" {
		return echo.ErrUnauthorized
	}

	claims := &jwtCustomClaims{
		username,
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("this_is_the_secret_key"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Token{Token: t})
}

type Test struct {
	Message string `json:"message"`
}

// accessible godoc
// @Summary Accessible endpoint
// @Description Returns a welcome message
// @ID accessible
// @Produce json
// @Success 200 {object} Test
// @Router / [get]
func accessible(c echo.Context) error {
	return c.JSON(http.StatusOK, Test{Message: "hello, welcome!"})
}

// @Security Bearer
// accessible godoc
// @Summary Restricted endpoint
// @Description Returns a welcome message
// @ID restricted
// @Produce json
// @Success 200 {object} Test
// @Router /restricted [get]
func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.JSON(http.StatusOK, Test{Message: "Welcome " + name + "!!!!!"})
}

// func main() {
// 	app := echo.New()
// 	app.GET("/swagger/*", echoSwagger.WrapHandler)

// 	app.Use(middleware.Recover())

// 	app.POST("/login", login)

// 	app.GET("/", accessible)

// 	restrict := app.Group("/restricted")

// 	config := echojwt.Config{
// 		NewClaimsFunc: func(c echo.Context) jwt.Claims {
// 			return new(jwtCustomClaims)
// 		},
// 		SigningKey: []byte("this_is_the_secret_key"),
// 	}
// 	restrict.Use(echojwt.WithConfig(config))
// 	restrict.GET("", restricted)

// 	app.Logger.Fatal(app.Start(":1323"))
// }

// func main() {
// 	app := New()
// 	app.Logger.Fatal(app.Start(":1323"))
// }

//@securityDefinitions.ApiKey Bearer
//@in header
//@name Authorization

// @title Swagger Example API
// @version 1.0
// @description This is a sample Swagger API for the Echo framework.
// @host localhost:1323
// @BasePath /
func New() *echo.Echo {
	app := echo.New()
	app.GET("/swagger/*", echoSwagger.WrapHandler)

	app.Use(middleware.Recover())

	app.POST("/login", login)

	app.GET("/", accessible)

	restrict := app.Group("/restricted")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("this_is_the_secret_key"),
	}
	restrict.Use(echojwt.WithConfig(config))
	restrict.GET("", restricted)

	return app
}
