package model

import (
	"gorm.io/gorm"
	"sduonline-recruitment/pkg/util"
	"sort"
	"strconv"
	"strings"
	"time"
)

type AplModel struct {
	AbstractModel
}
type Application struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	SectionID int `json:"section_id"`
	ApplicationJudge
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type ApplicationJudge struct {
	Status int    `json:"status" form:"status" binding:"gte=0,lte=2"`
	Note   string `json:"note" form:"note"`
}
type ApplicationJudgeVO struct {
	Status int    `json:"status" form:"status" binding:"gte=0,lte=2"`
	Score  []int  `json:"score" form:"score" gorm:"-"`
	Note   string `json:"note" form:"note"`
}
type ApplicationVO struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	SectionID int `json:"section_id"`
	ApplicationJudgeVO
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (receiver AplModel) UpdateUserApl(userID int, sectionID []int) []Application {
	tx := receiver.Tx.Begin()
	err := tx.Exec("delete from applications where user_id=?", userID).Error
	util.ForwardOrRollback(err, tx)
	var applications []Application
	for _, id := range sectionID {
		applications = append(applications, Application{
			UserID:           userID,
			SectionID:        id,
			ApplicationJudge: ApplicationJudge{},
		})
	}
	err = tx.Create(&applications).Error
	util.ForwardOrRollback(err, tx)
	tx.Commit()
	return applications
}
func (receiver AplModel) FindAplsByUserID(userID int) []AplListVO {
	var apl []AplListVO
	err := receiver.Tx.Raw(getAplListSQL("where a.user_id=?"), userID).Find(&apl).Error
	util.ForwardOrPanic(err)
	sort.Slice(apl, func(i, j int) bool {
		return apl[i].ID < apl[j].ID
	})
	return apl
}

type AplListVO struct {
	RealName    string `json:"real_name"`
	StudentID   string `json:"student_id"`
	Qq          string `json:"qq"`
	Phone       string `json:"phone"`
	SectionName string `json:"section_name"`
	ApplicationVO
}

func getAplListSQL(where string) string {
	return "select a.*,u.real_name,u.student_id,u.qq,u.phone,s.name section_name from applications a " +
		"left join users u on u.id=a.user_id " +
		"left join sections s on s.id=a.section_id " +
		where + " " +
		"order by s.id,a.created_at desc"
}
func (receiver AplModel) FindAllApls() []AplListVO {
	var apl []AplListVO
	err := receiver.Tx.Raw(getAplListSQL("")).Find(&apl).Error
	util.ForwardOrPanic(err)
	return apl
}
func (receiver AplModel) FindAplsBySectionID(sectionID int) []AplListVO {
	var apl []AplListVO
	err := receiver.Tx.Raw(getAplListSQL("where a.section_id=?"), sectionID).Find(&apl).Error
	util.ForwardOrPanic(err)
	return apl
}
func (receiver AplModel) FindAplsBySectionIDs(sectionIDs []int) []AplListVO {
	var apl []AplListVO
	var secIDStr []string
	for _, id := range sectionIDs {
		secIDStr = append(secIDStr, strconv.Itoa(id))
	}
	err := receiver.Tx.Raw(getAplListSQL("where a.section_id in (" + strings.Join(secIDStr, ",") + ")")).Find(&apl).Error
	util.ForwardOrPanic(err)
	return apl
}
func (receiver AplModel) FindAplByID(id int) *Application {
	var apl Application
	err := receiver.Tx.Take(&apl, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	util.ForwardOrPanic(err)
	return &apl
}

type aplAnsVO struct {
	Answer
	QuestionStem string `json:"stem"`
	SectionName  string `json:"section_name,omitempty"`
}

func (receiver AplModel) FindAplAnswersBySectionIDAplID(sectionID int, aplID int) []aplAnsVO {
	var vo []aplAnsVO
	err := receiver.Tx.Raw("select ans.*,q.stem question_stem,s.name section_name "+
		"from answers ans "+
		"left join questions q on ans.question_id=q.id "+
		"left join sections s on s.id=q.section_id "+
		"where s.id=? and ans.user_id="+
		"(select user_id from applications where id=?) "+
		"union "+
		"select ans.*,q.stem question_stem,null section_name "+
		"from answers ans "+
		"left join questions q on ans.question_id=q.id "+
		"where q.section_id is null and ans.user_id="+
		"(select user_id from applications where id=?)", sectionID, aplID, aplID).Find(&vo).Error
	util.ForwardOrPanic(err)
	return vo
}
func (receiver AplModel) UpdateAplJudgeByIDSectionID(judge *ApplicationJudge, id int, sectionID int) {
	receiver.Tx.Exec("update applications set status=?,note=?,updated_at=? where id=? and section_id=?",
		judge.Status, judge.Note, time.Now(), id, sectionID)
}
func (receiver AplModel) GetUserIDByAplID(aplID int) (int, error) {
	var apl Application
	err := receiver.Tx.Take(&apl, aplID).Error
	if err != nil {
		return 0, err
	}
	return apl.UserID, nil
}
