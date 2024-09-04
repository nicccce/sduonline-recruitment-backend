package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type config struct {
	Dsn         string `yaml:"dsn"`
	JWTSecret   string `yaml:"jwtSecret"`
	Port        int    `yaml:"port"`
	WxAppID     string `yaml:"wxAppID"`
	WxAppSecret string `yaml:"wxAppSecret"`
	UploadDir   string `yaml:"uploadDir"`
}

var Conf config

func Setup() {
	byt, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(byt, &Conf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", Conf)
}
