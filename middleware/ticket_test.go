package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/foxdex/ftx-site/pkg/jwt"

	"github.com/foxdex/ftx-site/pkg/consts"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTicket(t *testing.T) {
	t.Parallel()

	t.Run("without ticket middleware", func(t *testing.T) {
		r := gin.New()
		r.POST("/", func(context *gin.Context) {})

		req := httptest.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 200)
	})

	t.Run("ticket is invalid", func(t *testing.T) {
		r := gin.New()
		r.Use(Ticket())
		r.POST("/", func(context *gin.Context) {})

		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set(consts.HeaderDRAWTOKEN, "invalid token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 403)
	})

	t.Run("ticket is valid", func(t *testing.T) {
		r := gin.New()
		r.Use(Ticket())
		r.POST("/", func(context *gin.Context) {})

		req := httptest.NewRequest("POST", "/", nil)

		claims := jwt.NewUserClaims("1", "2", "3")
		token, err := claims.Generator()
		assert.NoError(t, err)
		req.Header.Set(consts.HeaderDRAWTOKEN, token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 200)
	})
}
