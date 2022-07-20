package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/foxdex/ftx-site/config"

	"github.com/foxdex/ftx-site/pkg/lucky"

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
			name:        "when inviter email is invalid, should get 400",
			status:      http.StatusBadRequest,
			errorReason: `Error:Field validation for 'InviterEmail' failed on the 'email' tag`,
			payload: &dto.LuckyCreateTokenReq{
				Email:        "123@gmail.com",
				KycLevel:     "KYC0",
				Personality:  "IATC",
				InviterEmail: "invalid email",
			},
		},
		{
			name:   "when without inviter email, create token should successfully",
			status: http.StatusOK,
			payload: &dto.LuckyCreateTokenReq{
				Email:       "123@gmail.com",
				KycLevel:    "KYC0",
				Personality: "IATC",
			},
		},
		{
			name:   "when params are ok, create token should successfully",
			status: http.StatusOK,
			payload: &dto.LuckyCreateTokenReq{
				Email:        "123@gmail.com",
				KycLevel:     "KYC0",
				Personality:  "IATC",
				InviterEmail: "456@gmail.com",
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
				return
			}

			var (
				uc  *jwt.UserClaims
				rsp dto.LuckyCreateTokenRsp
			)
			data, err := json.Marshal(rw.Data)
			assert.NoError(t, err)
			err = json.Unmarshal(data, &rsp)
			assert.NoError(t, err)

			claims, err := uc.Parse(rsp.Token)
			assert.NoError(t, err)
			assert.Equal(t, tc.payload.Email, claims.Email, "checking email")
			assert.Equal(t, tc.payload.KycLevel, claims.KycLevel, "checking kycLevel")
			assert.Equal(t, tc.payload.Personality, claims.Personality, "checking personality")
			assert.Equal(t, tc.payload.InviterEmail, claims.InviterEmail, "checking inviter email")
			assert.Equal(t, config.GetConfig().Jwt.Issuer, claims.Issuer, "checking issuer")
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
		mockFn      func(*TestCase)
		email       string
	}

	luckyModel := &dao.LuckyModel{
		LuckyId:     1,
		Email:       "123@gmail.com",
		KycLevel:    "KYC2",
		Personality: "IATC",
		Prize:       lucky.Prize30,
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
			email:       "invalid",
		},
		{
			name:        "get lucky model by email error",
			status:      http.StatusBadRequest,
			errorReason: `mock error`,
			email:       "123@gmail.com",
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}
				luckyDao.EXPECT().GetByEmail(gomock.Any(), "123@gmail.com").Return(nil, errors.New("mock error"))
			},
		},
		{
			name:   "success",
			status: http.StatusOK,
			email:  "123@gmail.com",
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
			c.Params = gin.Params{
				{
					Key:   "email",
					Value: tc.email,
				},
			}
			c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/lucky/%s", tc.email), nil)
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
		lucky       *LuckyDrawHandler
		payload     *dto.LuckyAwardReq
		mockFn      func(*TestCase)
	}

	cases := []TestCase{
		{
			name:        "invalid prize",
			status:      http.StatusBadRequest,
			errorReason: `Error:Field validation for 'Prize' failed on the 'oneof' tag`,
			payload: &dto.LuckyAwardReq{
				Prize:     lucky.Prize1000,
				UserName:  "jw",
				UserPhone: "12311112222",
				Address:   "xx省xx市xx区xxxx小区",
			},
		},
		{
			name:        "invalid address",
			status:      http.StatusBadRequest,
			errorReason: `Error:Field validation for 'Address' failed on the 'required' tag`,
			payload: &dto.LuckyAwardReq{
				Prize:     lucky.Prize40,
				UserName:  "jw",
				UserPhone: "12311112222",
			},
		},
		{
			name:        "clothes size is required when prize is FTX x MLB 棒球外套 or FTX清凉防晒衣",
			status:      http.StatusBadRequest,
			errorReason: "clothes size is required when prize is FTX x MLB 棒球外套 or FTX清凉防晒衣",
			payload: &dto.LuckyAwardReq{
				Prize:      lucky.Prize30,
				UserName:   "jw",
				UserPhone:  "12311112222",
				Address:    "xx省xx市xx区xxxx小区",
				Country:    "cn",
				Region:     "sz",
				PostalCode: "000000",
			},
		},
		{
			name:        "db mock error",
			status:      http.StatusInternalServerError,
			errorReason: `mock error`,
			payload: &dto.LuckyAwardReq{
				Prize:      lucky.Prize50,
				UserName:   "jw",
				UserPhone:  "12311112222",
				Address:    "xx省xx市xx区xxxx小区",
				Country:    "cn",
				Region:     "sz",
				PostalCode: "000000",
			},
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}
				luckyDao.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
		},
		{
			name:   "success",
			status: http.StatusOK,
			payload: &dto.LuckyAwardReq{
				Prize:      lucky.Prize50,
				UserName:   "jw",
				UserPhone:  "12311112222",
				Address:    "xx省xx市xx区xxxx小区",
				Country:    "cn",
				Region:     "sz",
				PostalCode: "000000",
			},
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}

				input := map[string]interface{}{
					"user_name":    "jw",
					"user_phone":   "12311112222",
					"address":      "xx省xx市xx区xxxx小区",
					"clothes_size": "",
					"country":      "cn",
					"region":       "sz",
					"postal_code":  "000000",
				}
				luckyDao.EXPECT().Update(gomock.Any(), "123@gmail.com", input).Return(nil)
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

			claims := jwt.NewUserClaims("123@gmail.com", "KYC2", "IATC", "456@gmail.com")
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
			}
		})
	}
}

func TestGetJackpot(t *testing.T) {
	mockDao := mock.NewMockILucky(gomock.NewController(t))
	mockDao.EXPECT().Count(gomock.Any()).Return(int64(0), nil)
	handler := &LuckyDrawHandler{mockDao}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/lucky/jackpot", nil)
	lucky.InitJackpot(mockDao)
	handler.GetJackpot(c)
	assert.Equal(t, 200, w.Code, "checking status code")

	var rw dto.ResponseFormat
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &rw), "unmarshalling response body")

	var rsp dto.LuckyGetJackpotRsp
	data, err := json.Marshal(rw.Data)
	assert.NoError(t, err)
	err = json.Unmarshal(data, &rsp)
	assert.NoError(t, err)
	assert.Equal(t, uint64(5000), rsp.Jackpot)

}

func TestDraw(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		status      int
		errorReason string
		lucky       *LuckyDrawHandler
		mockFn      func(*TestCase)
	}

	cases := []TestCase{
		{
			name:        "check email exist error",
			status:      http.StatusInternalServerError,
			errorReason: `mock error`,
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}
				luckyDao.EXPECT().EmailExist(gomock.Any(), "123@gmail.com").Return(false, errors.New("mock error"))
			},
		},
		{
			name:        "lottery has been conducted",
			status:      http.StatusBadRequest,
			errorReason: `lottery has been conducted`,
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}
				luckyDao.EXPECT().EmailExist(gomock.Any(), "123@gmail.com").Return(true, nil)
			},
		},
		{
			name:        "db mock error",
			status:      http.StatusInternalServerError,
			errorReason: `mock error`,
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}
				input := &dao.LuckyModel{
					Prize:        lucky.Prize1000,
					Email:        "123@gmail.com",
					KycLevel:     "KYC2",
					Personality:  "IATC",
					InviterEmail: "456@gmail.com",
				}
				luckyDao.EXPECT().EmailExist(gomock.Any(), "123@gmail.com").Return(false, nil)
				luckyDao.EXPECT().Create(gomock.Any(), input).Return(errors.New("mock error"))
			},
		},
		{
			name:   "success",
			status: http.StatusOK,
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}
				input := &dao.LuckyModel{
					Prize:        lucky.Prize1000,
					Email:        "123@gmail.com",
					KycLevel:     "KYC2",
					Personality:  "IATC",
					InviterEmail: "456@gmail.com",
				}
				luckyDao.EXPECT().EmailExist(gomock.Any(), "123@gmail.com").Return(false, nil)
				luckyDao.EXPECT().Create(gomock.Any(), input).Return(nil)
				luckyDao.EXPECT().Count(gomock.Any()).Return(int64(0), nil)
				lucky.InitJackpot(luckyDao)
			},
		},
	}

	err := os.Setenv("UNIT_TEST", "true")
	assert.NoError(t, err)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/lucky/draw", nil)
			c.Request.Header.Set("Content-Type", "application/json")

			claims := jwt.NewUserClaims("123@gmail.com", "KYC2", "IATC", "456@gmail.com")
			c.Set(consts.HeaderDRAWTOKEN, claims)

			if tc.mockFn != nil {
				tc.mockFn(&tc)
			}
			tc.lucky.Draw(c)
			assert.Equal(t, tc.status, w.Code, "checking status code")

			var rw dto.ResponseFormat
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &rw), "unmarshalling response body")
			if tc.errorReason != "" {
				assert.Contains(t, rw.Msg, tc.errorReason, "checking error reason")
				return
			}

			var rsp dto.LuckyDrawRsp
			data, err := json.Marshal(rw.Data)
			assert.NoError(t, err)
			err = json.Unmarshal(data, &rsp)
			assert.NoError(t, err)
			assert.Equal(t, lucky.Prize1000, rsp.Prize, "checking prize")
		})
	}
}

func TestGetTickets(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		status      int
		errorReason string
		lucky       *LuckyDrawHandler
		mockFn      func(*TestCase)
		email       string
		expected    uint64
	}

	cases := []TestCase{
		{
			name:        "when email is invalid, should get 400",
			status:      http.StatusBadRequest,
			errorReason: `Error:Field validation for 'Email' failed on the 'email' tag`,
			email:       "invalid",
		},
		{
			name:   "db mock error",
			status: http.StatusInternalServerError,
			email:  "123@gmail.com",
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}
				luckyDao.EXPECT().CountByEmail(gomock.Any(), "123@gmail.com").Return(int64(0), errors.New("mock error"))
			},
			errorReason: "mock error",
		},
		{
			name:   "invite more than 3 people",
			status: http.StatusOK,
			email:  "123@gmail.com",
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}
				luckyDao.EXPECT().CountByEmail(gomock.Any(), "123@gmail.com").Return(int64(4), nil)
			},
			expected: 4,
		},
		{
			name:   "invite less than 3 people",
			status: http.StatusOK,
			email:  "123@gmail.com",
			mockFn: func(testCase *TestCase) {
				luckyDao := mock.NewMockILucky(gomock.NewController(t))
				testCase.lucky = &LuckyDrawHandler{luckyDao}
				luckyDao.EXPECT().CountByEmail(gomock.Any(), "123@gmail.com").Return(int64(1), nil)
			},
			expected: 2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{
				{
					Key:   "email",
					Value: tc.email,
				},
			}
			c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/lucky/tickets/%s", tc.email), nil)
			if tc.mockFn != nil {
				tc.mockFn(&tc)
			}
			tc.lucky.GetTickets(c)
			assert.Equal(t, tc.status, w.Code, "checking status code")

			var rw dto.ResponseFormat
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &rw), "unmarshalling response body")
			if tc.errorReason != "" {
				assert.Contains(t, rw.Msg, tc.errorReason, "checking error reason")
				return
			}

			var rsp dto.GetTicketsRsp
			data, err := json.Marshal(rw.Data)
			assert.NoError(t, err)
			err = json.Unmarshal(data, &rsp)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, rsp.Count, "checking count")
		})
	}
}
