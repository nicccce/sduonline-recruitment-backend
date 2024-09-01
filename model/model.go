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
	if err := DB.AutoMigrate(&Interview{}); err != nil {
		return err
	}
	/*sqlDB, err := DB.DB()
	rows, err := sqlDB.Query("SHOW TABLES")
	if err != nil {
		return err
	}
	defer rows.Close()

	// 删除所有表
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			fmt.Println("Error scanning table name:", err)
			continue
		}

		result := DB.Exec("DROP TABLE IF EXISTS " + tableName)
		if result.Error != nil {
			return result.Error
		}
	}

	if err := DB.AutoMigrate(&User{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&SectionPermission{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&WebLogin{}); err != nil {
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

	if err := DB.AutoMigrate(&Config{}); err != nil {
		return err
	}*/

	return nil
}
