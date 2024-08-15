package main

import (
	"fmt"
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"sduonline-recruitment/middleware"
	"sduonline-recruitment/model"
	"sduonline-recruitment/pkg/app"
	"sduonline-recruitment/pkg/conf"
	"sduonline-recruitment/router"
	"sduonline-recruitment/service"
	"strconv"
)

func main() {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(nice.Recovery(func(c *gin.Context, err interface{}) {
		aw := app.NewWrapper(c)
		aw.Error("内部错误，请重试：" + fmt.Sprintf("%v", err))
	}))
	engine.Use(middleware.Cors())
	conf.Setup()
	router.Setup(engine)
	model.Setup()
	service.Setup()
	engine.Run(":" + strconv.Itoa(conf.Conf.Port))
}
