package service

import (
	"github.com/gin-gonic/gin"
	"sduonline-recruitment/model"
	"sduonline-recruitment/pkg/app"
	"sduonline-recruitment/pkg/util"
	"strconv"
)

type HRService struct {
}

func (receiver HRService) ListAllUsers(c *gin.Context) {
	aw := app.NewWrapper(c)
	users := userModel.FindAllUsers()
	perms := userModel.FindAllSectionPermissions()
	type VO struct {
		model.User
		Permissions []model.SectionPermissionVO `json:"permissions"`
	}
	userID2vo := map[int]VO{}
	var ids []int
	for _, user := range users {
		userID2vo[user.ID] = VO{
			User:        user,
			Permissions: []model.SectionPermissionVO{},
		}
		ids = append(ids, user.ID)
	}
	for _, perm := range perms {
		item := userID2vo[perm.UserID]
		item.Permissions = append(item.Permissions, perm)
		userID2vo[perm.UserID] = item
	}
	var voList []VO
	for _, id := range ids {
		voList = append(voList, userID2vo[id])
	}
	aw.Success(voList)
}
func (receiver HRService) GetUserInfo(c *gin.Context) {
	aw := app.NewWrapper(c)
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		aw.Error("user_id必须为整数")
		return
	}
	user, perms := userModel.FindUserByID(userID)
	aw.Success(gin.H{
		"user":        user,
		"permissions": perms,
	})
}
func (receiver HRService) AlterUserRole(c *gin.Context) {
	aw := app.NewWrapper(c)
	uc := util.ExtractUserClaims(c)
	type AlterRoleReq struct {
		UserID int `form:"user_id" binding:"required"`
		RoleID int `form:"role_id" binding:"required,gte=0,lte=5"`
	}
	var req AlterRoleReq
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	if req.RoleID > uc.RoleID {
		aw.Error("只能将用户角色改为自己同级或更低的")
		return
	}
	user, _ := userModel.FindUserByID(req.UserID)
	if user == nil {
		aw.Error("用户不存在")
		return
	}
	if user.RoleID >= uc.RoleID {
		aw.Error("不能操作与自己同级或更高的用户的角色")
		return
	}
	user.RoleID = req.RoleID
	userModel.UpdateUser(user)
	aw.Success(user)
}
func (receiver HRService) GrantPermission(c *gin.Context) {
	aw := app.NewWrapper(c)
	sectionID := util.ExtractSectionID(c)
	userID, err := strconv.Atoi(c.PostForm("user_id"))
	if err != nil {
		aw.Error("user_id必须为整数")
		return
	}
	user, _ := userModel.FindUserByID(userID)
	if user == nil {
		aw.Error("用户不存在")
		return
	}
	if user.RoleID < 2 {
		aw.Error("请先修改用户角色")
		return
	}
	if userModel.ExistSectionPermission(userID, sectionID) {
		aw.Error("用户已有该权限")
		return
	}
	sp := userModel.CreateSectionPermission(userID, sectionID)
	aw.Success(sp)
}
func (receiver HRService) RevokePermission(c *gin.Context) {
	aw := app.NewWrapper(c)
	sectionID := util.ExtractSectionID(c)
	userID, err := strconv.Atoi(c.PostForm("user_id"))
	if err != nil {
		aw.Error("user_id必须为整数")
		return
	}
	user, _ := userModel.FindUserByID(userID)
	if user == nil {
		aw.Error("用户不存在")
		return
	}
	if !userModel.ExistSectionPermission(userID, sectionID) {
		aw.Error("用户未有该权限")
		return
	}
	userModel.DeleteSectionPermission(userID, sectionID)
	aw.OK()
}
