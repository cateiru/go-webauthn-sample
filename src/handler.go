package src

import (
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	WebAuthn *webauthn.WebAuthn
}

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

func (h *Handler) BeginLoginHandler(c echo.Context) error {
	return nil
}

func (h *Handler) LoginHandler(c echo.Context) error {
	return nil
}

func (h *Handler) DebugDatabase(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"users":    Users,
		"sessions": Sessions,
	})
}
