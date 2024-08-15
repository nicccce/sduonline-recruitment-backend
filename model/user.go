package model

import (
	"gorm.io/gorm"
	"math/rand"
	"sduonline-recruitment/pkg/util"
	"strconv"
	"time"
)

type User struct {
	ID         int    `json:"id"`
	RoleID     int    `json:"role_id"`
	OpenID     string `json:"-"`
	SessionKey string `json:"-"`
	UserInfo
}
type UserInfo struct {
	RealName  string `json:"real_name" form:"real_name" binding:"required"`
	StudentID string `json:"student_id" form:"student_id" binding:"required"`
	Faculty   string `json:"faculty" form:"faculty" binding:"required"`
	Qq        string `json:"qq" form:"qq" binding:"required"`
	Phone     string `json:"phone" form:"phone" binding:"required"`
	Wechat    string `json:"wechat" form:"wechat"`
	School    string `json:"school" form:"school" binding:"required"`
	WxBypassUserInfo
}
type WxBypassUserInfo struct {
	Intro string `json:"intro" form:"intro"`
}
type SectionPermission struct {
	UserID    int `json:"user_id" gorm:"primaryKey"`
	SectionID int `json:"section_id" gorm:"primaryKey"`
}
type WebLogin struct {
	UserID    int `gorm:"primaryKey"`
	Code      string
	CreatedAt time.Time
}

func (w WebLogin) TableName() string {
	return "web_login"
}

type UserModel struct {
	AbstractModel
}

func (receiver UserModel) FindUserByID(id int) (*User, []SectionPermission) {
	var user User
	err := receiver.Tx.Take(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	util.ForwardOrPanic(err)
	var sp []SectionPermission
	err = receiver.Tx.Where("user_id=?", user.ID).Find(&sp).Error
	util.ForwardOrPanic(err)
	return &user, sp
}
func (receiver UserModel) UpdateUser(user *User) {
	err := receiver.Tx.Save(user).Error
	util.ForwardOrPanic(err)
}
func (receiver UserModel) ExistSectionPermission(userID int, sectionID int) bool {
	var secPerm SectionPermission
	err := receiver.Tx.Where("user_id=? and section_id=?", userID, sectionID).Take(&secPerm).Error
	if err == gorm.ErrRecordNotFound {
		return false
	}
	util.ForwardOrPanic(err)
	return true
}
func (receiver UserModel) CreateWebLogin(userID int) *WebLogin {
	tx := receiver.Tx.Begin()
	err := tx.Exec("delete from web_login where user_id=?", userID).Error
	util.ForwardOrRollback(err, tx)
	max := 999999
	min := 100000
	webLogin := WebLogin{
		UserID: userID,
		Code:   strconv.Itoa(rand.Intn(max-min) + min),
	}
	err = tx.Save(&webLogin).Error
	util.ForwardOrRollback(err, tx)
	tx.Commit()
	return &webLogin
}
func (receiver UserModel) ValidateWebLogin(userID int, code string) bool {
	var webLogin WebLogin
	err := receiver.Tx.Where("user_id=? and code=? and created_at>?", userID, code,
		time.Now().Add(-5*time.Minute)).Take(&webLogin).Error
	if err == gorm.ErrRecordNotFound {
		return false
	}
	util.ForwardOrPanic(err)
	err = receiver.Tx.Exec("delete from web_login where user_id=?", userID).Error
	util.ForwardOrPanic(err)
	return true
}
func (receiver UserModel) FindUserByOpenID(openID string) *User {
	var user User
	err := receiver.Tx.Where("open_id=?", openID).Take(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	util.ForwardOrPanic(err)
	return &user
}
func (receiver UserModel) CreateUser(user *User) {
	err := receiver.Tx.Create(user).Error
	util.ForwardOrPanic(err)
}
func (receiver UserModel) FindAllUsers() []User {
	var users []User
	err := receiver.Tx.Order("id desc").Find(&users).Error
	util.ForwardOrPanic(err)
	return users
}

type SectionPermissionVO struct {
	SectionPermission
	SectionName string `json:"section_name"`
}

func (receiver UserModel) FindAllSectionPermissions() []SectionPermissionVO {
	var list []SectionPermissionVO
	err := receiver.Tx.Raw("select sp.*,s.name section_name from section_permissions sp " +
		"left join sections s on s.id=sp.section_id ").Find(&list).Error
	util.ForwardOrPanic(err)
	return list
}
func (receiver UserModel) CreateSectionPermission(userID int, sectionID int) *SectionPermission {
	sp := SectionPermission{
		UserID:    userID,
		SectionID: sectionID,
	}
	err := receiver.Tx.Save(&sp).Error
	util.ForwardOrPanic(err)
	return &sp
}
func (receiver UserModel) DeleteSectionPermission(userID int, sectionID int) {
	err := receiver.Tx.Exec("delete from section_permissions where user_id=? and section_id=?", userID, sectionID).Error
	util.ForwardOrPanic(err)
}
