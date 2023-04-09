package src

import (
	"fmt"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	WebAuthn *webauthn.WebAuthn
}

// WebAuthnを新規に登録するためのチャレンジを返す。
// ユーザーとセッションはDBに保存し、クッキーにセッションIDを保存している。
func (h *Handler) BeginCreateHandler(c echo.Context) error {
	name := c.QueryParam("name")
	displayName := c.QueryParam("display_name")

	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}

	user := &User{
		ID:          []byte(RandomString(20)),
		Name:        name,
		DisplayName: displayName,
	}

	options, session, err := h.WebAuthn.BeginRegistration(user)
	if err != nil {
		return err
	}

	id := RandomString(10)

	InsertSession(id, session)
	InsertUser(id, user)

	c.SetCookie(&http.Cookie{
		Name:  "session_id",
		Value: id,
	})

	return c.JSON(http.StatusOK, options)
}

// WebAuthのクレデンシャルを検証して成功の場合は、新規にクレデンシャルを登録する。
// CookieからセッションIDを取得し、DBからユーザーとセッションを取得する。
func (h *Handler) CreateHandler(c echo.Context) error {
	id, err := c.Cookie("session_id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "session_id is required")
	}

	response, err := protocol.ParseCredentialCreationResponseBody(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	user, err := GetUser(id.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user not found")
	}
	session, err := GetSession(id.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "session not found")
	}

	credential, err := h.WebAuthn.CreateCredential(user, *session, response)
	if err != nil {
		return err
	}

	user.Credentials = append(user.Credentials, *credential)
	InsertUser(id.Value, user)

	return nil
}

// ログイン用のチャレンジを返す
// ログインページで毎回生成する。
// セッションはDBに保存して、IDをクッキーに保存している。
func (h *Handler) BeginLoginHandler(c echo.Context) error {
	credential, session, err := h.WebAuthn.BeginDiscoverableLogin()
	if err != nil {
		return err
	}

	id := RandomString(10)

	InsertSession(id, session)
	c.SetCookie(&http.Cookie{
		Name:  "login_session_id",
		Value: id,
	})

	return c.JSON(http.StatusOK, credential)
}

func (h *Handler) LoginHandler(c echo.Context) error {
	id, err := c.Cookie("login_session_id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "session_id is required")
	}
	session, err := GetSession(id.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "session not found")
	}

	response, err := protocol.ParseCredentialRequestResponseBody(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	var loggedInUser *User = nil
	handler := func(rawID, userHandle []byte) (user webauthn.User, err error) {
		u, err := GetUserById(userHandle)
		if err != nil {
			fmt.Printf("user not found: %v", err)
			return nil, err
		}
		loggedInUser = u
		return u, nil
	}

	_, err = h.WebAuthn.ValidateDiscoverableLogin(handler, *session, response)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, loggedInUser)
}

func (h *Handler) DebugDatabase(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"users":    Users,
		"sessions": Sessions,
	})
}
