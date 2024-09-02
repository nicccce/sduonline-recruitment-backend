package router

import (
	"github.com/gin-gonic/gin"
	"sduonline-recruitment/middleware"
	"sduonline-recruitment/service"
)

func Setup(engine *gin.Engine) {
	// 测试 上线后注释掉
	/*	test := engine.Group("/test")
		{
			//测试panic
			test.GET("/panic", func(c *gin.Context) {
				panic("test panic")
			})
			//初始化数据库
			test.GET("/database_initialization", func(c *gin.Context) {
				aw := app.NewWrapper(c)
				err := model.Database_initialization()
				if err != nil {
					aw.Error(err.Error())
				}
				aw.Success("success!")
			})
		}*/

	// 用户
	user := engine.Group("/user")
	{
		hub := service.UserService{}
		// 测试方法，获取signed jwt
		user.GET("/test_get_jwt", hub.TestGetJWT)
		// 网页通过授权码登录
		user.POST("/web_login", hub.WebLogin)
		// 微信登录
		user.POST("/wechat_login", hub.WeChatLogin)
	}
	user.Use(middleware.JWT(1))
	{
		hub := service.UserService{}
		// 个人信息
		user.GET("/me", hub.Me)
		// 更新个人信息（学号、姓名等基础信息）
		user.POST("/update_info", hub.UpdateInfo)
		// 获取更新用户信息表单
		user.GET("/get_update_info_form", hub.GetUpdateInfoForm)
		// 获取网页登录授权码，授权码有效期五分钟
		user.POST("/generate_web_login", hub.GenerateWebLogin)
	}

	// 事业群-部门
	depSec := engine.Group("/dep_sec")
	{
		hub := service.DepSecService{}
		// 获取各部门介绍
		depSec.GET("/intro", hub.Intro)
	}
	depSec.Use(middleware.JWT(3))
	{
		hub := service.DepSecService{}
		// 编辑部门介绍
		depSec.POST("/edit_section_intro", middleware.SectionPermission(), hub.EditSectionIntro)
		// 列出部门所有（纳新）问题
		depSec.GET("/list_section_questions", middleware.SectionPermission(), hub.ListQuestions)
		// 添加一个部门问题
		depSec.POST("/add_section_question", middleware.SectionPermission(), hub.AddQuestion)
		// 编辑一个部门问题
		depSec.POST("/edit_section_question", middleware.SectionPermission(), hub.EditQuestion)
		// 删除一个部门问题
		depSec.POST("/delete_section_question", middleware.SectionPermission(), hub.DeleteQuestion)
		// 列出部门的所有面试
		depSec.GET("/interview", middleware.SectionPermission(), hub.ListInterviews)
		// 添加一个部门面试
		depSec.POST("/add_interview", middleware.SectionPermission(), hub.AddInterview)
		// 编辑一个部门面试
		depSec.POST("/edit_interview", middleware.SectionPermission(), hub.EditInterview)
		// 删除一个部门面试
		depSec.POST("/delete_interview", middleware.SectionPermission(), hub.DeleteInterview)
	}
	depSec.Use(middleware.JWT(4))
	{
		hub := service.DepSecService{}
		// 编辑分站介绍
		depSec.POST("/edit_department_intro", hub.EditDepartmentIntro)
		// 列出所有通用（纳新）问题
		depSec.GET("/list_universal_questions", hub.ListQuestions)
		// 添加一个通用问题
		depSec.POST("/add_universal_question", hub.AddQuestion)
		// 编辑一个通用问题
		depSec.POST("/edit_universal_question", hub.EditQuestion)
		// 删除一个通用问题
		depSec.POST("/delete_universal_question", hub.DeleteQuestion)
	}

	// 入部申请单
	apl := engine.Group("/apl")
	apl.Use(middleware.JWT(1))
	{
		hub := service.AplService{}
		// 申请者提交申请部门
		apl.POST("/submit_apl", hub.SubmitApl)
		// 申请者获取自己需要回答的问题
		apl.GET("/get_questions", hub.GetQuestions)
		// 申请者提交问题答案
		apl.POST("/submit_answers", hub.SubmitAnswers)
		// 获取用户所有问题与答案。以答案为准，即使apl被删掉，倘若之前回答过某个问题，也会列出
		apl.GET("/get_my_answers", hub.GetUserAnswers)
	}
	apl.Use(middleware.JWT(2))
	{
		hub := service.AplService{}
		// 总监操作
		// 列出某个部门下所有申请单，需要用户有部门权限；如果有dump参数且不为空则导出xlsx
		apl.GET("/list_section_apl", middleware.SectionPermission(), hub.ListSectionApl)
		// 按权限导出所有部门apl
		apl.GET("/export_apls_by_permission", hub.ExportAplsByPermission)
		// 获取某个申请单的所有问题和回答
		apl.GET("/get_apl_answers", middleware.SectionPermission(), hub.GetAplAnswers)
		// 对某申请单评分
		apl.POST("/judge_apl", middleware.SectionPermission(), hub.JudgeApl)
		// 列出用户所有申请单
		apl.GET("/list_user_apls", hub.ListUserApls)
	}
	apl.Use(middleware.JWT(4))
	{
		hub := service.AplService{}
		// 站长操作
		// 获取所有部门申请单
		apl.GET("/list_all_apl", hub.ListAllApl)
	}

	// 人力资源管理
	hr := engine.Group("/hr")
	hr.Use(middleware.JWT(3))
	{
		hub := service.HRService{}
		// 获取所有用户列表，不包括用户部门权限
		hr.GET("/list_all_users", hub.ListAllUsers)
		// 获取某个用户信息，包括用户的部门权限
		hr.GET("/get_user_info", hub.GetUserInfo)
		// 改变用户角色
		hr.POST("/alter_user_role", hub.AlterUserRole)
		// 授予部门权限
		hr.POST("/grant_permission", middleware.SectionPermission(), hub.GrantPermission)
		// 收回部门权限
		hr.POST("/revoke_permission", middleware.SectionPermission(), hub.RevokePermission)
	}
}
