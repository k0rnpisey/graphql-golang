package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"not null; searchable"`

	CanBeTopping    bool       `gorm:"default:false"`
	ProductToppings []*Product `gorm:"many2many:ProductTopping;"`
}

func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
	err = UniqNameWithinScope(tx, p, p.Name, p.ID)
	if err != nil {
		return err
	}
	return nil
}

func UniqNameWithinScope(tx *gorm.DB, model interface{}, name string, id uint) error {
	var existingCount int64
	result := tx.Debug().Model(model).Where("name = ? AND id != ?", name, id).Count(&existingCount)
	if result.Error != nil {
		return result.Error
	}
	if existingCount > 0 {
		return fmt.Errorf("name '%s' must be unique among active records", name)
	}
	return nil
}
