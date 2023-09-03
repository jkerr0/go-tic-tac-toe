package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func GetSession(c echo.Context) *sessions.Session {
	sess, err := session.Get("session", c)
	if err != nil {
		c.Logger().Error("cannot get session", err)
		c.String(http.StatusInternalServerError, "Cannot get session")
		return nil
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	return sess
}

func SaveSession(sess *sessions.Session, c echo.Context) {
	err := sess.Save(c.Request(), c.Response())
	if err != nil {
		c.Logger().Error("cannot save session", err)
		c.String(http.StatusInternalServerError, "Cannot save session")
	}
}
