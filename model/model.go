package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sduonline-recruitment/pkg/conf"
)

var DB *gorm.DB

type AbstractModel struct {
	Tx *gorm.DB
}

func Setup() {
	dbInternal, err := gorm.Open(mysql.Open(conf.Conf.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = dbInternal
}

func Database_initialization() error {
	if err := DB.AutoMigrate(&User{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&Question{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&Answer{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&Application{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&Department{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&Section{}); err != nil {
		return err
	}

	return nil
}
