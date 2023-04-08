package src

import (
	"math/rand"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	wConfig := &webauthn.Config{
		RPDisplayName: "Test CAT",                        // Display Name for your site
		RPID:          "localhost",                       // Generally the FQDN for your site
		RPOrigins:     []string{"http://localhost:3000"}, // The origin URLs allowed for WebAuthn requests
	}

	w, err := webauthn.New(wConfig)
	if err != nil {
		panic(err)
	}

	h := &Handler{
		WebAuthn: w,
	}

	e.GET("/begin_create", h.BeginCreateHandler)
	e.POST("/create", h.CreateHandler)
	e.GET("/begin_login", h.BeginLoginHandler)
	e.POST("/login", h.LoginHandler)

	// debug
	e.GET("/debug/database", h.DebugDatabase)

	e.Logger.Fatal(e.Start(":1323"))
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
