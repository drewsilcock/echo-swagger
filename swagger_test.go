package echoSwagger

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	_ "github.com/swaggo/gin-swagger/example/basic/docs"
)

func TestWrapHandler(t *testing.T) {
	router := echo.New()

	router.GET("/*", WrapHandler)

	w1 := performRequest("GET", "/index.html", router)
	assert.Equal(t, 200, w1.Code)

	w2 := performRequest("GET", "/", router)
	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, w1.Body.String(), w2.Body.String())

	w3 := performRequest("GET", "/doc.json", router)
	assert.Equal(t, 200, w3.Code)

	w4 := performRequest("GET", "/favicon-16x16.png", router)
	assert.Equal(t, 200, w4.Code)

	w5 := performRequest("GET", "/notfound", router)
	assert.Equal(t, 404, w5.Code)

	w6 := performRequest("GET", "/index.htmlnotfound", router)
	assert.Equal(t, 404, w6.Code)
}

func TestConfig(t *testing.T) {
	router := echo.New()

	swaggerHandler := EchoWrapHandler(func(c *Config) {
		c.URL = "some-other-name.json"
	})
	router.GET("/*", swaggerHandler)

	w1 := performRequest("GET", "/index.html", router)
	assert.Equal(t, 200, w1.Code)

	w2 := performRequest("GET", "/some-other-name.json", router)
	assert.Equal(t, 200, w2.Code)
}

func TestConfigWithOAuth(t *testing.T) {
	router := echo.New()

	swaggerHandler := EchoWrapHandler(func(c *Config) {
		c.OAuth = &OAuthConfig{
			ClientId: "my-client-id",
			Realm:    "my-realm",
			AppName:  "My App Name",
		}
	})
	router.GET("/*", swaggerHandler)

	w1 := performRequest("GET", "/index.html", router)
	assert.Equal(t, 200, w1.Code)
	body1 := w1.Body.String()
	assert.Contains(t, body1, "initOAuth({")
	assert.Contains(t, body1, "clientId: \"my-client-id\"")
	assert.Contains(t, body1, "realm: \"my-realm\"")
	assert.Contains(t, body1, "appName: \"My App Name\"")

	w2 := performRequest("GET", "/doc.json", router)
	assert.Equal(t, 200, w2.Code)

	w3 := performRequest("GET", "/favicon-16x16.png", router)
	assert.Equal(t, 200, w3.Code)

	w4 := performRequest("GET", "/notfound", router)
	assert.Equal(t, 404, w4.Code)
}

func performRequest(method, target string, e *echo.Echo) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()

	e.ServeHTTP(w, r)
	return w
}
