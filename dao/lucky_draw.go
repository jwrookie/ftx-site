package dao

import (
	"github.com/foxdex/ftx-site/pkg/utils"
	"gorm.io/gorm"
)

type LuckyModel struct {
	LuckyId     uint64 `json:"lucky_id,string" gorm:"primaryKey"`
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

type ILucky interface {
	GetByEmail(db *gorm.DB, email string) (*LuckyModel, error)
	EmailExist(db *gorm.DB, email string) (bool, error)
	Create(db *gorm.DB, model *LuckyModel) error
	Update(db *gorm.DB, email string, updates map[string]interface{}) error
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

func (l *LuckyHandler) EmailExist(db *gorm.DB, email string) (bool, error) {
	var (
		err   error
		count int64
	)
	if err = db.Table("lucky").Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (l *LuckyHandler) Update(db *gorm.DB, email string, updates map[string]interface{}) error {
	var err error
	if err = db.Table("lucky").Where("email = ?", email).UpdateColumns(updates).Error; err != nil {
		return err
	}

	return nil
}
