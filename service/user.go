package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"html"
	"os"
	"sduonline-recruitment/model"
	"sduonline-recruitment/pkg/app"
	"sduonline-recruitment/pkg/conf"
	"sduonline-recruitment/pkg/util"
	"strconv"
	"strings"
)

type UserService struct {
}

func (receiver UserService) TestGetJWT(c *gin.Context) {
	aw := app.NewWrapper(c)
	type getJWTReq struct {
		UserID    int    `form:"user_id" binding:"required"`
		JWTSecret string `form:"jwt_secret" binding:"required"`
	}
	var req getJWTReq
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	if req.JWTSecret != conf.Conf.JWTSecret {
		aw.Error("jwtSecret不正确")
		return
	}
	user, _ := userModel.FindUserByID(req.UserID)
	if user == nil {
		aw.Error("userID不存在")
		return
	}
	aw.Success(util.GenerateJWT(user.ID, user.RoleID))
}
func (receiver UserService) WeChatLogin(c *gin.Context) {
	aw := app.NewWrapper(c)
	code := c.PostForm("code")
	if code == "" {
		aw.Error("请传递code")
		return
	}
	resp, err := util.GetCode2Session(code)
	if err != nil {
		aw.Error("无法调用微信API：" + err.Error())
		return
	}
	user := userModel.FindUserByOpenID(resp.Openid)
	if user == nil {
		user = &model.User{
			RoleID:     1,
			OpenID:     resp.Openid,
			SessionKey: resp.SessionKey,
		}
		userModel.CreateUser(user)
	} else {
		user.SessionKey = resp.SessionKey
		userModel.UpdateUser(user)
	}
	aw.Success(util.GenerateJWT(user.ID, user.RoleID))
}
func (receiver UserService) Me(c *gin.Context) {
	aw := app.NewWrapper(c)
	uc := util.ExtractUserClaims(c)
	user, permissions := userModel.FindUserByID(uc.UserID)
	if user == nil {
		aw.Error("该接口需要登录")
		return
	}
	apl := aplModel.FindAplsByUserID(user.ID)
	aw.Success(gin.H{"user": user, "permissions": permissions, "apls": apl})
}
func (receiver UserService) UpdateInfo(c *gin.Context) {
	aw := app.NewWrapper(c)
	uc := util.ExtractUserClaims(c)
	bypassOption := configModel.FindConfigByName("wxBypassMode")
	user, _ := userModel.FindUserByID(uc.UserID)
	if bypassOption.Value == "1" {
		var bypassUserInfo model.WxBypassUserInfo
		if err := c.ShouldBind(&bypassUserInfo); err != nil {
			aw.Error(err.Error())
			return
		}
		user.UserInfo.WxBypassUserInfo = model.WxBypassUserInfo{
			Intro: html.EscapeString(bypassUserInfo.Intro),
		}
	} else {
		var userInfo model.UserInfo
		if err := c.ShouldBind(&userInfo); err != nil {
			aw.Error(err.Error())
			return
		}
		user.UserInfo = model.UserInfo{
			RealName:  html.EscapeString(userInfo.RealName),
			StudentID: html.EscapeString(userInfo.StudentID),
			Faculty:   html.EscapeString(userInfo.Faculty),
			Qq:        html.EscapeString(userInfo.Qq),
			Phone:     html.EscapeString(userInfo.Phone),
			Wechat:    html.EscapeString(userInfo.Wechat),
			School:    html.EscapeString(userInfo.School),
			WxBypassUserInfo: model.WxBypassUserInfo{
				Intro: html.EscapeString(userInfo.Intro),
			},
		}
	}
	userModel.UpdateUser(user)
	aw.Success(user)
}
func (receiver UserService) GetUpdateInfoForm(c *gin.Context) {
	aw := app.NewWrapper(c)
	type KV struct {
		Key  string `json:"key"`
		Text string `json:"text"`
	}
	type WxBypassForm struct {
		Normal []KV `json:"normal"`
		Bypass []KV `json:"bypass"`
	}
	var form WxBypassForm
	fileContent, err := os.ReadFile("wxBypassForm.json")
	if err != nil {
		aw.Error(err.Error())
		return
	}
	err = json.Unmarshal(fileContent, &form)
	if err != nil {
		aw.Error(err.Error())
		return
	}
	bypassOption := configModel.FindConfigByName("wxBypassMode")
	if bypassOption.Value == "1" {
		aw.Success(form.Bypass)
	} else {
		aw.Success(form.Normal)
	}
}
func (receiver UserService) GenerateWebLogin(c *gin.Context) {
	aw := app.NewWrapper(c)
	uc := util.ExtractUserClaims(c)
	webLogin := userModel.CreateWebLogin(uc.UserID)
	aw.Success(strconv.Itoa(webLogin.UserID) + ":" + webLogin.Code)
}
func (receiver UserService) WebLogin(c *gin.Context) {
	aw := app.NewWrapper(c)
	login := c.PostForm("login")
	loginArr := strings.Split(login, ":")
	if len(loginArr) != 2 {
		aw.Error("授权码不正确")
		return
	}
	userID, err := strconv.Atoi(loginArr[0])
	if err != nil {
		aw.Error("userID必须为整数")
		return
	}
	user, _ := userModel.FindUserByID(userID)
	if userModel.ValidateWebLogin(userID, loginArr[1]) {
		aw.Success(util.GenerateJWT(userID, user.RoleID))
	} else {
		aw.Error("授权码不存在或已过期")
	}
}
