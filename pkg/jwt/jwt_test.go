package jwt

import (
	"testing"
	"time"

	"github.com/foxdex/ftx-site/config"
	"github.com/stretchr/testify/assert"
)

func TestUserClaims(t *testing.T) {
	t.Parallel()

	var validToken string
	t.Run("generator token should successfully", func(t *testing.T) {
		uc := NewUserClaims("11", "22", 1)
		tokenStr, err := uc.Generator()
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenStr)
		validToken = tokenStr
	})

	t.Run("when token string is invalid, parse token should failed", func(t *testing.T) {
		var uc *UserClaims
		claims, err := uc.Parse("invalid token string")
		assert.Nil(t, claims)
		assert.NotNil(t, err)
	})

	t.Run("when token string is valid, parse token should successfully", func(t *testing.T) {
		var uc *UserClaims
		claims, err := uc.Parse(validToken)
		assert.NoError(t, err)
		assert.Equal(t, "11", claims.Email, "checking email")
		assert.Equal(t, "22", claims.Personality, "checking personality")
		assert.Equal(t, uint8(1), claims.KycLevel, "checking kycLevel")
		assert.Equal(t, "", claims.Prize, "checking prize")
		assert.Equal(t, config.GetConfig().Jwt.Issuer, claims.Issuer, "checking issuer")
		assert.InDelta(t, time.Now().Add(time.Hour*time.Duration(24*7)).Unix(), claims.ExpiresAt.Unix(), 10)
	})
}
