package dto

type LuckyDto struct {
	LuckyId     uint64 `json:"lucky_id,string"`
	Email       string `json:"email"`
	KycLevel    uint8  `json:"kyc_level"`
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
