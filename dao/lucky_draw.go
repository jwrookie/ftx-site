package dao

type LuckyModel struct {
	LuckyId     uint64 `gorm:"primaryKey"`
	Email       string
	KycLevel    uint8
	Personality string
	Prize       string
	ClothesSize string
	UserName    string
	UserPhone   string
	Address     string

	CreatedAt uint64
	UpdatedAt uint64
	DeletedAt uint64
}

type ILucky interface {
	GetByEmail(email string) (*LuckyModel, error)
	Create(model *LuckyModel) error
}

type LuckyHandler struct{}

func (l *LuckyHandler) GetByEmail(email string) (*LuckyModel, error) {
	return nil, nil
}

func (l *LuckyHandler) Create(model *LuckyModel) error {
	return nil
}
