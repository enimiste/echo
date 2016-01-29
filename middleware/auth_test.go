package middleware

import (
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/test"
	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", nil)
	res := test.NewResponseRecorder()
	c := echo.NewContext(req, res, e)
	fn := func(u, p string) bool {
		if u == "joe" && p == "secret" {
			return true
		}
		return false
	}
	ba := BasicAuth(fn)

	// Valid credentials
	auth := Basic + " " + base64.StdEncoding.EncodeToString([]byte("joe:secret"))
	req.Header().Set(echo.Authorization, auth)
	assert.NoError(t, ba(c))

	//---------------------
	// Invalid credentials
	//---------------------

	// Incorrect password
	auth = Basic + " " + base64.StdEncoding.EncodeToString([]byte("joe:password"))
	req.Header().Set(echo.Authorization, auth)
	he := ba(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusUnauthorized, he.Code())
	assert.Equal(t, Basic+" realm=Restricted", res.Header().Get(echo.WWWAuthenticate))

	// Empty Authorization header
	req.Header().Set(echo.Authorization, "")
	he = ba(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusUnauthorized, he.Code())
	assert.Equal(t, Basic+" realm=Restricted", res.Header().Get(echo.WWWAuthenticate))

	// Invalid Authorization header
	auth = base64.StdEncoding.EncodeToString([]byte("invalid"))
	req.Header().Set(echo.Authorization, auth)
	he = ba(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusUnauthorized, he.Code())
	assert.Equal(t, Basic+" realm=Restricted", res.Header().Get(echo.WWWAuthenticate))

	// WebSocket
	c.Request().Header().Set(echo.Upgrade, echo.WebSocket)
	assert.NoError(t, ba(c))
}
