package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sduonline-recruitment/model"
	"sduonline-recruitment/pkg/app"
	"sduonline-recruitment/pkg/util"
)

type SectionReq struct {
	SectionID int `form:"section_id" binding:"required"`
}

func SectionPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		aw := app.NewWrapper(c)
		uc := util.ExtractUserClaims(c)
		var req SectionReq
		if err := c.ShouldBind(&req); err != nil {
			aw.Error(err.Error())
			c.Abort()
			return
		}
		depSecModel := model.DepSecModel{AbstractModel: model.AbstractModel{Tx: model.DB}}
		if !depSecModel.ExistSectionsByID([]int{req.SectionID}) {
			aw.Error("section不存在")
			c.Abort()
			return
		}
		if uc.RoleID >= 4 {
			c.Set("sectionID", req.SectionID)
			return
		}
		exist := model.UserModel{AbstractModel: model.AbstractModel{Tx: model.DB}}.ExistSectionPermission(uc.UserID, req.SectionID)
		if !exist {
			aw.Error(fmt.Sprintf("用户 UserID:%v 未拥有 SectionID:%v 的权限", uc.UserID, req.SectionID))
			c.Abort()
			return
		}
		c.Set("sectionID", req.SectionID)
	}
}
