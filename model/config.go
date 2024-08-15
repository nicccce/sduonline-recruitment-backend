package model

import (
	"gorm.io/gorm"
	"sduonline-recruitment/pkg/util"
)

type ConfigModel struct {
	AbstractModel
}
type Config struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c Config) TableName() string {
	return "config"
}

func (c ConfigModel) UpdateConfig(name string, value string) {
	tx := c.Tx.Begin()
	err := tx.Exec("delete from config where name=?", name).Error
	util.ForwardOrPanic(err)
	config := Config{
		Name:  name,
		Value: value,
	}
	err = tx.Create(&config).Error
	util.ForwardOrPanic(err)
	tx.Commit()
}
func (c ConfigModel) FindConfigByName(name string) *Config {
	var config Config
	err := c.Tx.Take(&config, "name=?", name).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	util.ForwardOrPanic(err)
	return &config
}
