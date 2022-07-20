package dto

import "github.com/foxdex/ftx-site/dao"

type LuckyCreateTokenReq struct {
	Email        string `json:"email" form:"email" binding:"required,email"`
	KycLevel     string `json:"kyc_level" form:"kyc_level" binding:"required,oneof=KYC0 KYC1 KYC2"`
	Personality  string `json:"personality" form:"personality" binding:"required,oneof=IATC EATC IATM EATM IAFC EAFC IAFM EAFM IPTC EPTC IPTM EPTM IPFC EPFC IPFM EPTM"`
	InviterEmail string `json:"inviter_email" form:"inviter_email" binding:"omitempty,email"`
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
	Prize       string `json:"prize" form:"prize" binding:"oneof='FTX三周年礼盒' 'FTX x AMG联合棒球帽' 'FTX x MLB 棒球外套' '交易手续费抵扣券10USD' 'FTX祝福红包' 'FTX清凉防晒衣' 'FTX绒绒袜' 'FTX小龙人暖手充电宝' 'FTX雪花真空杯+小金勺子' 'FTX超萌小耳朵发箍' 'FTX定制纸牌'"`
	ClothesSize string `json:"clothes_size" form:"clothes_size"`
	UserName    string `json:"user_name" form:"user_name" binding:"required"`
	UserPhone   string `json:"user_phone" form:"user_phone" binding:"required"`
	Address     string `json:"address" form:"address" binding:"required"`
	Country     string `json:"country" form:"country" binding:"required"`
	Region      string `json:"region" form:"region" binding:"required"`
	PostalCode  string `json:"postal_code" form:"postal_code" binding:"required"`
}

type LuckyAwardRsp struct {
	dao.LuckyModel `json:",inline" mapstructure:",squash"`
}

type LuckyGetJackpotRsp struct {
	Jackpot uint64 `json:"jackpot,string"`
}

type LuckyDrawRsp struct {
	Prize string `json:"prize"`
}

type GetTicketsReq struct {
	Email string `form:"email" uri:"email" binding:"required,email"`
}

type GetTicketsRsp struct {
	Count uint64 `json:"count,string"`
}
