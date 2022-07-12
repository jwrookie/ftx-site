package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/foxdex/ftx-site/dto"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

var lucky LuckyDrawHandler

func TestCreateToken(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		status      int
		errorReason string
		payload     *dto.LuckyCreateTokenReq
	}

	cases := []TestCase{
		{
			name:        "when email is invalid, should get 400",
			status:      http.StatusBadRequest,
			errorReason: `Field validation for 'Email' failed on the 'email' tag`,
			payload: &dto.LuckyCreateTokenReq{
				Email:       "invalid",
				KycLevel:    "KYC0",
				Personality: "IATC",
			},
		},
		{
			name:        "when kyc level is invalid, should get 400",
			status:      http.StatusBadRequest,
			errorReason: `Error:Field validation for 'KycLevel' failed on the 'oneof' tag`,
			payload: &dto.LuckyCreateTokenReq{
				Email:       "123@gmail.com",
				KycLevel:    "invalid",
				Personality: "IATC",
			},
		},
		{
			name:        "when personality is invalid, should get 400",
			status:      http.StatusBadRequest,
			errorReason: `Error:Field validation for 'Personality' failed on the 'oneof' tag`,
			payload: &dto.LuckyCreateTokenReq{
				Email:       "123@gmail.com",
				KycLevel:    "KYC0",
				Personality: "invalid",
			},
		},
		{
			name:   "when params are ok, create token should successfully",
			status: http.StatusOK,
			payload: &dto.LuckyCreateTokenReq{
				Email:       "123@gmail.com",
				KycLevel:    "KYC0",
				Personality: "IATC",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/lucky/token", nil)
			c.Request.Header.Set("Content-Type", "application/json")
			body, _ := json.Marshal(tc.payload)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			lucky.CreateToken(c)
			assert.Equal(t, tc.status, w.Code, "checking status code")

			var rw dto.ResponseFormat
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &rw), "unmarshalling response body")
			if tc.errorReason != "" {
				assert.Contains(t, rw.Msg, tc.errorReason, "checking error reason")
			}
		})
	}
}
