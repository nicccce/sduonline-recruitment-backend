package model

import (
	"gopkg.in/guregu/null.v3"
	"sduonline-recruitment/pkg/util"
)

type InterviewModel struct {
	AbstractModel
}

type Interview struct {
	ID        int      `json:"id"`
	SectionID int      `json:"section_id"`
	Name      string   `json:"name"`
	Weight    null.Int `json:"weight"`
}
type InterviewDTO struct {
	Name   string `form:"name" binding:"required"`
	Weight int    `form:"weight" binding:"required"`
}

func (receiver InterviewModel) CreateInterview(interviewDTO InterviewDTO, sectionID *int) Interview {
	interview := Interview{
		SectionID: *sectionID,
		Name:      interviewDTO.Name,
		Weight:    null.NewInt(int64(interviewDTO.Weight), true),
	}
	receiver.Tx.Create(interview)
	return interview
}
func (receiver InterviewModel) UpdateInterview(interviewDTO InterviewDTO, id int, sectionID *int) Interview {
	interview := Interview{
		ID:        id,
		SectionID: *sectionID,
		Name:      interviewDTO.Name,
		Weight:    null.NewInt(int64(interviewDTO.Weight), true),
	}
	err := receiver.Tx.Save(interview).Error
	util.ForwardOrPanic(err)
	return interview
}
func (receiver InterviewModel) DeleteInterview(id int, sectionID *int) {
	err := receiver.Tx.Exec("DELETE FROM interview WHERE id=? AND section_id=?", id, *sectionID).Error
	util.ForwardOrPanic(err)
}
func (receiver InterviewModel) FindInterviewsBySectionID(sectionID *int) []Interview {
	var interviews []Interview
	err := receiver.Tx.Find(&interviews, "section_id=?", *sectionID).Error
	util.ForwardOrPanic(err)
	return interviews
}
