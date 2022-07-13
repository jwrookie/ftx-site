package dto

import "github.com/foxdex/ftx-site/dao"

type LuckyCreateTokenReq struct {
	Email       string `json:"email" form:"email" binding:"required,email"`
	KycLevel    string `json:"kyc_level" form:"kyc_level" binding:"required,oneof=KYC0 KYC1 KYC2"`
	Personality string `json:"personality" form:"personality" binding:"required,oneof=IATC EATC IATM EATM IAFC EAFC IAFM EAFM IPTC EPTC IPTM EPTM IPFC EPFC IPFM EPTM"`
}

type LuckyCreateTokenRsp struct {
	Token string `json:"token"`
}

type LuckyGetResultReq struct {
	Email string `form:"email" uri:"email" binding:"required,email"`
}

type LuckyGetResultRsp struct {
	dao.LuckyModel `json:",inline" mapstructure:",squash"`
}

type LuckyAwardReq struct {
	ClothesSize string `json:"clothes_size" form:"clothes_size"`
	UserName    string `json:"user_name" form:"user_name" binding:"required"`
	UserPhone   string `json:"user_phone" form:"user_phone" binding:"required"`
	Address     string `json:"address" form:"address" binding:"required"`
}

type LuckyAwardRsp struct {
	dao.LuckyModel `json:",inline" mapstructure:",squash"`
}

type LuckyGetJackpotRsp struct {
	Jackpot uint64 `json:"jackpot,string"`
}
