package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCors(t *testing.T) {
	t.Run("without cors middleware", func(t *testing.T) {
		r := gin.New()
		r.GET("/", func(context *gin.Context) {})

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 200)
		assert.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "")
		assert.Equal(t, w.Header().Get("Access-Control-Allow-Credentials"), "")
		assert.Equal(t, w.Header().Get("Access-Control-Allow-Headers"), "")
		assert.Equal(t, w.Header().Get("Access-Control-Allow-Methods"), "")
	})

	t.Run("with cors middleware", func(t *testing.T) {
		r := gin.New()
		r.Use(Cors())
		r.GET("/", func(context *gin.Context) {})

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 200)

		assert.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "")
		assert.Equal(t, w.Header().Get("Access-Control-Allow-Credentials"), "true")
		assert.Equal(t, w.Header().Get("Access-Control-Allow-Headers"), "CSRF-TOKEN")
		assert.Equal(t, w.Header().Get("Access-Control-Allow-Methods"), "POST, GET, OPTIONS,PUT,DELETE")
		assert.Equal(t, w.Header().Get("Access-Control-Expose-Headers"), "CSRF-TOKEN")
	})
}
