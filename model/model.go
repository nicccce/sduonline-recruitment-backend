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
