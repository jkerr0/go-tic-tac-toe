package handlers

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Context struct {
	Db      *sqlx.DB
	EchoCtx echo.Context
}