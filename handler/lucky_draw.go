package handler

import (
	"net/http"

	"github.com/foxdex/ftx-site/pkg/lucky"

	"github.com/foxdex/ftx-site/pkg/consts"

	"github.com/mitchellh/mapstructure"

	"github.com/foxdex/ftx-site/pkg/db"

	"github.com/foxdex/ftx-site/dao"

	"github.com/foxdex/ftx-site/pkg/log"
	"go.uber.org/zap"

	"github.com/foxdex/ftx-site/pkg/jwt"

	"github.com/foxdex/ftx-site/dto"
	"github.com/gin-gonic/gin"
)

var DefaultLuckyDrawHandler = &LuckyDrawHandler{
	luckyDao: &dao.LuckyHandler{},
}

type LuckyDrawHandler struct {
	luckyDao dao.ILucky
}

// CreateToken get qualified for the lottery.
func (h *LuckyDrawHandler) CreateToken(c *gin.Context) {
	var req dto.LuckyCreateTokenReq
	if err := c.ShouldBind(&req); err != nil {
		dto.FailResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := jwt.NewUserClaims(req.Email, req.KycLevel, req.Personality).Generator()
	if err != nil {
		dto.FailResponse(c, http.StatusInternalServerError, err.Error())
		log.Log.Error("create json web token error", zap.Error(err))
		return
	}

	dto.SuccessResponse(c, &dto.LuckyCreateTokenRsp{
		Token: token,
	})
}

// Draw conduct a lottery
func (h *LuckyDrawHandler) Draw(c *gin.Context) {
	var (
		model dao.LuckyModel
		exist bool
		err   error
	)

	claims := c.MustGet(consts.HeaderDRAWTOKEN).(*jwt.UserClaims)
	if exist, err = h.luckyDao.EmailExist(db.Mysql(), claims.Email); err != nil {
		dto.FailResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if exist {
		dto.FailResponse(c, http.StatusBadRequest, "lottery has been conducted")
		return
	}

	prize := lucky.Draw()
	model.Prize = prize

	if err = mapstructure.Decode(claims, &model); err != nil {
		dto.FailResponse(c, http.StatusInternalServerError, err.Error())
		log.Log.Error("decode error", zap.Error(err), zap.Any("input", claims))
		return
	}

	if err = h.luckyDao.Create(db.Mysql(), &model); err != nil {
		dto.FailResponse(c, http.StatusInternalServerError, err.Error())
		log.Log.Error("create lucky", zap.Error(err))
		return
	}

	lucky.IncJackpot()

	dto.SuccessResponse(c, &dto.LuckyDrawRsp{
		Prize: prize,
	})
}

// Award receive your award
func (h *LuckyDrawHandler) Award(c *gin.Context) {
	var (
		req    dto.LuckyAwardReq
		err    error
		claims = c.MustGet(consts.HeaderDRAWTOKEN).(*jwt.UserClaims)
	)
	if err = c.ShouldBind(&req); err != nil {
		dto.FailResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Prize == lucky.Prize40 && req.ClothesSize == "" {
		dto.FailResponse(c, http.StatusBadRequest, "clothes size is required when prize is FTX灰色T恤")
		return
	}

	updates := map[string]interface{}{
		"clothes_size": req.ClothesSize,
		"user_name":    req.UserName,
		"user_phone":   req.UserPhone,
		"address":      req.Address,
	}

	if err = h.luckyDao.Update(db.Mysql(), claims.Email, updates); err != nil {
		dto.FailResponse(c, http.StatusInternalServerError, err.Error())
		log.Log.Error("update lucky", zap.Error(err))
		return
	}

	dto.SuccessResponse(c, nil)
}

// GetResult check if you have won a prize
func (h *LuckyDrawHandler) GetResult(c *gin.Context) {
	var (
		req dto.LuckyGetResultReq
		rsp dto.LuckyGetResultRsp
	)
	if err := c.ShouldBindUri(&req); err != nil {
		dto.FailResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	model, err := h.luckyDao.GetByEmail(db.Mysql(), req.Email)
	if err != nil {
		dto.FailResponse(c, http.StatusBadRequest, err.Error())
		log.Log.Error("get lucky model by email error", zap.Error(err))
		return
	}

	if err = mapstructure.Decode(model, &rsp); err != nil {
		dto.FailResponse(c, http.StatusInternalServerError, err.Error())
		log.Log.Error("decode error", zap.Error(err))
		return
	}

	dto.SuccessResponse(c, &rsp)
}

// GetJackpot get jackpot
func (h *LuckyDrawHandler) GetJackpot(c *gin.Context) {
	dto.SuccessResponse(c, &dto.LuckyGetJackpotRsp{Jackpot: lucky.GetJackpot()})
}
