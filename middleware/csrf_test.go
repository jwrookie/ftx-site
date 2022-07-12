package middleware

import (
	"encoding/base64"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/foxdex/ftx-site/config"
	"github.com/foxdex/ftx-site/pkg/consts"
	"github.com/foxdex/ftx-site/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCsrf(t *testing.T) {
	t.Run("without csrf middleware", func(t *testing.T) {
		r := gin.New()
		r.POST("/", func(context *gin.Context) {})

		req := httptest.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 200)
	})

	getCsrfToken := func(ts uint64) string {
		token, _ := utils.RsaEncrypt([]byte(strconv.Itoa(int(ts))), config.GetPublicKey())
		tokenBytes := base64.StdEncoding.EncodeToString(token)
		return tokenBytes
	}

	t.Run("csrf token success", func(t *testing.T) {
		ts := utils.UnixMilliNow()
		r := gin.New()
		r.Use(Csrf())
		r.POST("/", func(context *gin.Context) {})

		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set(consts.HeaderAUTHCSRFTOKEN, getCsrfToken(ts))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 200)
	})

	t.Run("csrf token error", func(t *testing.T) {
		ts := utils.UnixMilliNow()
		r := gin.New()
		r.Use(Csrf())
		r.POST("/", func(context *gin.Context) {})

		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set(consts.HeaderAUTHCSRFTOKEN, getCsrfToken(ts))

		time.Sleep(time.Millisecond * time.Duration(config.GetConfig().Csrf.Interval+1))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 403)
	})

	t.Run("csrf check time success", func(t *testing.T) {
		ts := utils.UnixMilliNow() + 1000
		r := gin.New()
		r.Use(Csrf())
		r.POST("/", func(context *gin.Context) {})

		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set(consts.HeaderAUTHCSRFTOKEN, getCsrfToken(ts))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 200)
	})

	t.Run("csrf check time error", func(t *testing.T) {
		ts := utils.UnixMilliNow() + uint64(config.GetConfig().Csrf.Interval) + 1000
		r := gin.New()
		r.Use(Csrf())
		r.POST("/", func(context *gin.Context) {})

		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set(consts.HeaderAUTHCSRFTOKEN, getCsrfToken(ts))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 403)
	})
}
