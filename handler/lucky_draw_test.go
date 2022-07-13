package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/foxdex/ftx-site/pkg/consts"
	"github.com/foxdex/ftx-site/pkg/jwt"

	"github.com/google/go-cmp/cmp"

	"github.com/foxdex/ftx-site/dao"
	"github.com/foxdex/ftx-site/dao/mock"
	"github.com/golang/mock/gomock"

	"github.com/foxdex/ftx-site/dto"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

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
			DefaultLuckyDrawHandler.CreateToken(c)
			assert.Equal(t, tc.status, w.Code, "checking status code")

			var rw dto.ResponseFormat
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &rw), "unmarshalling response body")
			if tc.errorReason != "" {
				assert.Contains(t, rw.Msg, tc.errorReason, "checking error reason")
			}
		})
	}
}

func TestResult(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		status      int
		errorReason string
		lucky       *LuckyDrawHandler
		payload     *dto.LuckyGetResultReq
		mockFn      func(*TestCase)
	}

	luckyModel := &dao.LuckyModel{
		LuckyId:     1,
		Email:       "123@gmail.com",
		KycLevel:    "KYC2",
		Personality: "IATC",
		Prize:       "FTX棒球帽",
		ClothesSize: "",
		UserName:    "JW",
		UserPhone:   "12311112222",
		Address:     "xx省xx市xx区xxxx小区",
		CreatedAt:   1,
		UpdatedAt:   1,
		DeletedAt:   1,
	}

	cases := []TestCase{
		{
			name:        "when email is invalid, should get 400",
			status:      http.StatusBadRequest,
			errorReason: `Field validation for 'Email' failed on the 'email' tag`,
			payload: &dto.LuckyGetResultReq{
				Email: "invalid",
			},
		},
		{
			name:        "get lucky model by email error",
			status:      http.StatusBadRequest,
			errorReason: `mock error`,
			payload: &dto.LuckyGetResultReq{
				Email: "123@gmail.com",
			},
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}
				luckyDao.EXPECT().GetByEmail(gomock.Any(), "123@gmail.com").Return(nil, errors.New("mock error"))
			},
		},
		{
			name:   "success",
			status: http.StatusOK,
			payload: &dto.LuckyGetResultReq{
				Email: "123@gmail.com",
			},
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}
				luckyDao.EXPECT().GetByEmail(gomock.Any(), "123@gmail.com").Return(luckyModel, nil)
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

			if tc.mockFn != nil {
				tc.mockFn(&tc)
			}
			tc.lucky.GetResult(c)
			assert.Equal(t, tc.status, w.Code, "checking status code")

			var rw dto.ResponseFormat
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &rw), "unmarshalling response body")
			if tc.errorReason != "" {
				assert.Contains(t, rw.Msg, tc.errorReason, "checking error reason")
				return
			}

			var newLucky dto.LuckyGetResultRsp
			data, err := json.Marshal(rw.Data)
			assert.NoError(t, err)
			err = json.Unmarshal(data, &newLucky)
			assert.NoError(t, err)
			diff := cmp.Diff(*luckyModel, newLucky.LuckyModel)
			assert.Empty(t, diff, "check lucky model")
		})
	}
}

func TestAward(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		status      int
		errorReason string
		prize       string
		lucky       *LuckyDrawHandler
		payload     *dto.LuckyAwardReq
		mockFn      func(*TestCase)
	}

	luckyModel := &dao.LuckyModel{
		Email:       "123@gmail.com",
		KycLevel:    "KYC2",
		Personality: "IATC",
		Prize:       "FTX棒球帽",
		ClothesSize: "",
		UserName:    "jw",
		UserPhone:   "12311112222",
		Address:     "xx省xx市xx区xxxx小区",
	}

	cases := []TestCase{
		{
			name:        "invalid params",
			status:      http.StatusBadRequest,
			errorReason: `Error:Field validation for 'Address' failed on the 'required' tag`,
			payload: &dto.LuckyAwardReq{
				UserName:  "jw",
				UserPhone: "12311112222",
			},
			prize: "FTX棒球帽",
		},
		{
			name:        "clothes size is required when prize is FTX灰色T恤",
			status:      http.StatusBadRequest,
			errorReason: "clothes size is required when prize is FTX灰色T恤",
			payload: &dto.LuckyAwardReq{
				UserName:  "jw",
				UserPhone: "12311112222",
				Address:   "xx省xx市xx区xxxx小区",
			},
			prize: "FTX灰色T恤",
		},
		{
			name:        "db mock error",
			status:      http.StatusInternalServerError,
			errorReason: `mock error`,
			payload: &dto.LuckyAwardReq{
				UserName:  "jw",
				UserPhone: "12311112222",
				Address:   "xx省xx市xx区xxxx小区",
			},
			prize: "FTX棒球帽",
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}
				luckyDao.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
		},
		{
			name:   "success",
			status: http.StatusOK,
			payload: &dto.LuckyAwardReq{
				UserName:  "jw",
				UserPhone: "12311112222",
				Address:   "xx省xx市xx区xxxx小区",
			},
			prize: "FTX棒球帽",
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}

				input := &dao.LuckyModel{
					UserName:    "jw",
					UserPhone:   "12311112222",
					Address:     "xx省xx市xx区xxxx小区",
					Prize:       "FTX棒球帽",
					Email:       "123@gmail.com",
					KycLevel:    "KYC2",
					Personality: "IATC",
				}
				luckyDao.EXPECT().Create(gomock.Any(), input).Return(nil)
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

			claims := jwt.NewUserClaims("123@gmail.com", "KYC2", "IATC")
			claims.Prize = tc.prize
			c.Set(consts.HeaderDRAWTOKEN, claims)

			if tc.mockFn != nil {
				tc.mockFn(&tc)
			}
			tc.lucky.Award(c)
			assert.Equal(t, tc.status, w.Code, "checking status code")

			var rw dto.ResponseFormat
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &rw), "unmarshalling response body")
			if tc.errorReason != "" {
				assert.Contains(t, rw.Msg, tc.errorReason, "checking error reason")
				return
			}

			var newLucky dto.LuckyAwardRsp
			data, err := json.Marshal(rw.Data)
			assert.NoError(t, err)
			err = json.Unmarshal(data, &newLucky)
			assert.NoError(t, err)
			diff := cmp.Diff(*luckyModel, newLucky.LuckyModel)
			assert.Empty(t, diff, "check lucky model")
		})
	}
}
