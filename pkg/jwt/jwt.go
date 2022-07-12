package jwt

import (
	"errors"
	"time"

	"github.com/foxdex/ftx-site/config"
	"github.com/foxdex/ftx-site/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
)

var (
	InvalidToken = errors.New("invalid token")
)

type UserClaims struct {
	Email       string `json:"email"`
	KycLevel    string `json:"kyc_level"`
	Personality string `json:"personality"`
	Prize       string `json:"prize"`
	jwt.RegisteredClaims
}

func NewUserClaims(email, kycLevel, personality string) *UserClaims {
	return &UserClaims{
		Email:       email,
		KycLevel:    kycLevel,
		Personality: personality,
	}
}

func (uc *UserClaims) Generator() (string, error) {
	jwtConfig := config.GetConfig().Jwt
	uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(24*7)))
	uc.Issuer = jwtConfig.Issuer

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenStr, err := token.SignedString([]byte(jwtConfig.SignKey))
	if err != nil {
		return "", err
	}
	return utils.Base64AESCBCEncrypt(tokenStr)
}

func (uc *UserClaims) Parse(tokenString string) (*UserClaims, error) {
	var err error
	jwtConfig := config.GetConfig().Jwt
	if tokenString, err = utils.Base64AESCBCDecrypt(tokenString); err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.SignKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, InvalidToken
	}

	return claims, nil
}
