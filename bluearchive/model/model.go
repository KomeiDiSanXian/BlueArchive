// Package model 数据库操作
package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Init(databasePath string) error {
	db, err := Open(databasePath)
	if err != nil {
		return err
	}
	defer db.Close()

	db.AutoMigrate(
		&Character{},
		&Profile{},
		&Property{},
		&Weapon{},
		&WeaponStar{},
		&WeaponSkill{},
		&Equipment{},
		&Adapatation{},
		&CharaSkill{},
		&SkillUsage{},
	)
	return nil
}

func Open(databasePath string) (*gorm.DB, error) {
	return gorm.Open("sqlite3", databasePath)
}
