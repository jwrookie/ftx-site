package dao

import (
	"github.com/foxdex/ftx-site/pkg/utils"
	"gorm.io/gorm"
)

type LuckyModel struct {
	LuckyId     uint64 `gorm:"primaryKey"`
	Email       string
	KycLevel    string
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
	GetByEmail(db *gorm.DB, email string) (*LuckyModel, error)
	Create(db *gorm.DB, model *LuckyModel) error
	Count(db *gorm.DB) (int64, error)
}

type LuckyHandler struct{}

func (l *LuckyHandler) GetByEmail(db *gorm.DB, email string) (*LuckyModel, error) {
	var (
		model LuckyModel
		err   error
	)
	if err = db.Table("lucky").Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (l *LuckyHandler) Create(db *gorm.DB, model *LuckyModel) error {
	var err error

	if model.LuckyId, err = utils.GenSnowflakeID(); err != nil {
		return err
	}

	now := utils.UnixMilliNow()
	model.CreatedAt = now
	model.UpdatedAt = now

	if err = db.Table("lucky").Create(&model).Error; err != nil {
		return err
	}

	return nil
}

func (l *LuckyHandler) Count(db *gorm.DB) (int64, error) {
	var (
		count int64
		err   error
	)

	if err = db.Table("lucky").Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
