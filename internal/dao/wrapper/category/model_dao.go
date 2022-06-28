package category

import (
	"github.com/glide-im/api/internal/dao/common"
	"github.com/glide-im/api/internal/pkg/db"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type Category struct {
	AppID  int64  `json:"app_id,omitempty"`
	Name   string `json:"title"`
	Weight int64  `json:"weight"`
	Icon   string `json:"icon"`
}

var CategoryDao = &CategoryH{}

type CategoryH struct {
}

func (a *CategoryH) GetModel(app_id int64) *gorm.DB {
	return db.DB.Model(&Category{}).Where("app_id = ?", app_id)
}

type CategoryUser struct {
	AppID      int64  `json:"app_id,omitempty"`
	CategoryId int64  `json:"category_id"`
	UId        string `json:"uid"`
}

var CategoryUserDao = &CategoryUserH{}

type CategoryUserH struct {
}

func (s *CategoryUserH) GetModel(app_id int64) *gorm.DB {
	return db.DB.Model(&CategoryUser{}).Where("app_id = ?", app_id)
}

func (s *CategoryUserH) Updates(uid int64, category_ids []int64) error {
	var _db = db.DB
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		_uid := cast.ToString(uid)
		_db = s.GetModel(1).Where("uid = ?", _uid).Delete(&Category{})
		if err := common.JustError(_db); err != nil {
			return err
		}

		var categories = []CategoryUser{}
		for _, category_id := range category_ids {
			categories = append(categories, CategoryUser{
				AppID:      1,
				CategoryId: category_id,
				UId:        _uid,
			})
		}
		_db = db.DB.CreateInBatches(categories, 100)
		if err := common.JustError(_db); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
