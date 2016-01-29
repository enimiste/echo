package middleware

import (
	"bytes"
	"errors"
	"net/http"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/test"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	// Note: Just for the test coverage, not a real test.
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", nil)
	res := test.NewResponseRecorder()
	c := echo.NewContext(req, res, e)

	// Status 2xx
	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	Logger()(h)(c)

	// Status 3xx
	res = test.NewResponseRecorder()
	c = echo.NewContext(req, res, e)
	h = func(c echo.Context) error {
		return c.String(http.StatusTemporaryRedirect, "test")
	}
	Logger()(h)(c)

	// Status 4xx
	res = test.NewResponseRecorder()
	c = echo.NewContext(req, res, e)
	h = func(c echo.Context) error {
		return c.String(http.StatusNotFound, "test")
	}
	Logger()(h)(c)

	// Status 5xx with empty path
	req = test.NewRequest(echo.GET, "", nil)
	res = test.NewResponseRecorder()
	c = echo.NewContext(req, res, e)
	h = func(c echo.Context) error {
		return errors.New("error")
	}
	Logger()(h)(c)
}

func TestLoggerIPAddress(t *testing.T) {
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", nil)
	res := test.NewResponseRecorder()
	c := echo.NewContext(req, res, e)
	buf := new(bytes.Buffer)
	e.Logger().SetOutput(buf)
	ip := "127.0.0.1"
	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	mw := Logger()

	// With X-Real-IP
	req.Header().Add(echo.XRealIP, ip)
	mw(h)(c)
	assert.Contains(t, buf.String(), ip)

	// With X-Forwarded-For
	buf.Reset()
	req.Header().Del(echo.XRealIP)
	req.Header().Add(echo.XForwardedFor, ip)
	mw(h)(c)
	assert.Contains(t, buf.String(), ip)

	// with req.RemoteAddr
	buf.Reset()
	mw(h)(c)
	assert.Contains(t, buf.String(), ip)
}
