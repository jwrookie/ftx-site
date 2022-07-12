package dto

type LuckyDto struct {
	LuckyId     uint64 `json:"lucky_id,string"`
	Email       string `json:"email"`
	KycLevel    string `json:"kyc_level"`
	Personality string `json:"personality"`
	Prize       string `json:"prize"`
	ClothesSize string `json:"clothes_size"`
	UserName    string `json:"user_name"`
	UserPhone   string `json:"user_phone"`
	Address     string `json:"address"`

	CreatedAt uint64 `json:"created_at"`
	UpdatedAt uint64 `json:"updated_at"`
	DeletedAt uint64 `json:"deleted_at"`
}

type LuckyCreateTokenReq struct {
	Email       string `json:"email" form:"email" binding:"required,email"`
	KycLevel    string `json:"kyc_level" form:"kyc_level" binding:"required,oneof=KYC0 KYC1 KYC2"`
	Personality string `json:"personality" form:"personality" binding:"required,oneof=IATC EATC IATM EATM IAFC EAFC IAFM EAFM IPTC EPTC IPTM EPTM IPFC EPFC IPFM EPTM"`
}

type LuckyCreateTokenRsp struct {
	Token string `json:"token"`
}

type LuckyGetResultReq struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}

type LuckyGetResultRsp struct {
	LuckyDto
	// 脱敏处理
	UserName  string `json:"-"`
	UserPhone string `json:"-"`
	Address   string `json:"-"`
}

type LuckyAwardReq struct {
	ClothesSize string `json:"clothes_size" form:"clothes_size"`
	UserName    string `json:"user_name" form:"user_name" binding:"required"`
	UserPhone   string `json:"user_phone" form:"user_phone" binding:"required"`
	Address     string `json:"address" form:"address" binding:"required"`
}

type LuckyAwardRsp struct {
	LuckyDto
}
