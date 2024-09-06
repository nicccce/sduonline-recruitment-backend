package model

import (
	"gopkg.in/guregu/null.v3"
	"gorm.io/gorm"
	"sduonline-recruitment/pkg/util"
)

type QnsAnsModel struct {
	AbstractModel
}
type Question struct {
	ID        int      `json:"id"`
	SectionID null.Int `json:"section_id"`
	Stem      string   `json:"stem"`
}
type Answer struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	QuestionID int    `json:"question_id"`
	Text       string `json:"answer"`
}
type userQuestionDTO struct {
	Question
	SectionName string `json:"section_name,omitempty"`
}

func (q QnsAnsModel) FindUserQuestions(userID int) []userQuestionDTO {
	var qns []userQuestionDTO
	err := q.Tx.Raw("select *,null section_name from questions where section_id is null "+
		"union "+
		"select q.*,s.name section_name from questions q "+
		"left join sections s on s.id=q.section_id "+
		"where q.section_id in "+
		"(select section_id from applications where user_id=?) ", userID).Find(&qns).Error
	util.ForwardOrPanic(err)
	return qns
}
func (q QnsAnsModel) DeleteAnswersByUserID(userID int) {
	err := q.Tx.Exec("delete from answers where user_id=?", userID).Error
	util.ForwardOrPanic(err)
}
func (q QnsAnsModel) CreateAnswers(answers []Answer) {
	err := q.Tx.Create(&answers).Error
	util.ForwardOrPanic(err)
}
func (q QnsAnsModel) FindQuestionsByOptionalSectionID(sectionID *int) []Question {
	var qns []Question
	var err error
	if sectionID == nil {
		err = q.Tx.Find(&qns, "section_id is null").Error
	} else {
		err = q.Tx.Find(&qns, "section_id=?", sectionID).Error
	}
	util.ForwardOrPanic(err)
	return qns
}
func (q QnsAnsModel) FindQuestionByID(id int) *Question {
	var qns Question
	err := q.Tx.Take(&qns, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	util.ForwardOrPanic(err)
	return &qns
}
func (q QnsAnsModel) CreateQuestion(stem string, sectionID *int) *Question {
	qns := Question{
		SectionID: null.NewInt(0, false),
		Stem:      stem,
	}
	if sectionID != nil {
		sid := *sectionID
		qns.SectionID = null.NewInt(int64(sid), true)
	}
	q.Tx.Create(&qns)
	return &qns
}
func (q QnsAnsModel) UpdateQuestion(question *Question) {
	err := q.Tx.Save(&question).Error
	util.ForwardOrPanic(err)
}
func (q QnsAnsModel) DeleteQuestionByOptionalSectionID(id int, sectionID *int) {
	if sectionID == nil {
		q.Tx.Exec("delete from questions where id=?", id)
	} else {
		q.Tx.Exec("delete from questions where id=? and section_id=?", id, *sectionID)
	}
}

type qnsAnsVO struct {
	QuestionID   int    `json:"question_id"`
	QuestionStem string `json:"question_stem"`
	SectionName  string `json:"section_name,omitempty"`
	AnswerID     int    `json:"answer_id"`
	AnswerText   string `json:"answer"`
}

func (q QnsAnsModel) FindQuestionsAnswersByUserID(userID int) []qnsAnsVO {
	var vo []qnsAnsVO
	err := q.Tx.Raw("select a.question_id,q.stem question_stem,s.name section_name,a.id answer_id,a.text answer_text "+
		"from answers a "+
		"left join questions q on q.id=a.question_id "+
		"left join sections s on q.section_id=s.id "+
		"where a.user_id=? and"+
		"(s.id IN (SELECT section_id FROM applications WHERE user_id = a.user_id) OR s.id IS NULL OR s.id = 0)", userID).Find(&vo).Error
	util.ForwardOrPanic(err)
	return vo
}
