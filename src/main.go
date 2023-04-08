package src

import (
	"math/rand"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

func Main() {
	e := echo.New()

	wConfig := &webauthn.Config{
		RPDisplayName: "Test",                                      // Display Name for your site
		RPID:          "go-webauthn.local",                         // Generally the FQDN for your site
		RPOrigins:     []string{"https://login.go-webauthn.local"}, // The origin URLs allowed for WebAuthn requests
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
