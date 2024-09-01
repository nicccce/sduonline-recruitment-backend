package service

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"
	"sduonline-recruitment/model"
	"sduonline-recruitment/pkg/app"
	"sduonline-recruitment/pkg/util"
	"strconv"
)

type DepSecService struct {
}

func (receiver DepSecService) Intro(c *gin.Context) {
	aw := app.NewWrapper(c)
	arr := depSecModel.FindAll()
	type VO struct {
		ID         int    `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
		Intro      string `json:"intro,omitempty"`
		IntroExtra string `json:"intro_extra"`
		Sections   []VO   `json:"sections,omitempty"`
	}
	depID2ix := map[int]int{}
	var vo []VO
	for _, depSec := range arr {
		if depSec.DepartmentID == 0 {
			depID2ix[depSec.ID] = len(vo)
			vo = append(vo, VO{
				ID:       depSec.ID,
				Name:     depSec.Name,
				Intro:    depSec.Intro,
				Sections: []VO{},
			})
		} else {
			vo[depID2ix[depSec.DepartmentID]].Sections = append(vo[depID2ix[depSec.DepartmentID]].Sections, VO{
				ID:         depSec.ID,
				Name:       depSec.Name,
				Intro:      depSec.Intro,
				IntroExtra: depSec.IntroExtra,
				Sections:   nil,
			})
		}
	}
	aw.Success(vo)
}
func (receiver DepSecService) EditSectionIntro(c *gin.Context) {
	aw := app.NewWrapper(c)
	sectionID := util.ExtractSectionID(c)
	section := depSecModel.FindSectionByID(sectionID)
	if section == nil {
		aw.Error("section不存在")
		return
	}
	var req model.SectionInfo
	if err := aw.Ctx.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	section.SectionInfo = req
	depSecModel.UpdateSectionByID(section)
	aw.Success(section)
}
func (receiver DepSecService) EditDepartmentIntro(c *gin.Context) {
	aw := app.NewWrapper(c)
	type DepIntroReq struct {
		DepartmentID int    `form:"department_id" binding:"required"`
		Intro        string `form:"intro"`
	}
	var req DepIntroReq
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	dep := depSecModel.FindDepartmentByID(req.DepartmentID)
	if dep == nil {
		aw.Error("department不存在")
		return
	}
	dep.Intro = req.Intro
	depSecModel.UpdateDepartment(dep)
	aw.Success(dep)
}
func (receiver DepSecService) ListQuestions(c *gin.Context) {
	aw := app.NewWrapper(c)
	var sectionID *int
	if util.HasSectionID(c) {
		sid := util.ExtractSectionID(c)
		sectionID = &sid
	}
	qns := qnsAnsModel.FindQuestionsByOptionalSectionID(sectionID)
	aw.Success(qns)
}
func (receiver DepSecService) AddQuestion(c *gin.Context) {
	aw := app.NewWrapper(c)
	if c.PostForm("stem") == "" {
		aw.Error("请输入题干")
		return
	}
	var sectionID *int
	if util.HasSectionID(c) {
		sid := util.ExtractSectionID(c)
		sectionID = &sid
	}
	qns := qnsAnsModel.CreateQuestion(c.PostForm("stem"), sectionID)
	aw.Success(qns)
}
func (receiver DepSecService) EditQuestion(c *gin.Context) {
	aw := app.NewWrapper(c)
	type EditQnsReq struct {
		ID   int    `form:"id" binding:"required"`
		Stem string `form:"stem" binding:"required"`
	}
	var req EditQnsReq
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	var sectionID *int
	var sid int
	if util.HasSectionID(c) {
		sid = util.ExtractSectionID(c)
		sectionID = &sid
	}
	qns := qnsAnsModel.FindQuestionByID(req.ID)
	if qns == nil {
		aw.Error("问题不存在")
		return
	}
	if sectionID != nil && !qns.SectionID.Equal(null.IntFrom(int64(sid))) {
		aw.Error("问题所属section_id与传入section_id不一致")
		return
	}
	qns.Stem = req.Stem
	qnsAnsModel.UpdateQuestion(qns)
	aw.Success(qns)
}
func (receiver DepSecService) DeleteQuestion(c *gin.Context) {
	aw := app.NewWrapper(c)
	var sectionID *int
	if util.HasSectionID(c) {
		sid := util.ExtractSectionID(c)
		sectionID = &sid
	}
	qnsID, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		aw.Error(err.Error())
		return
	}
	qnsAnsModel.DeleteQuestionByOptionalSectionID(qnsID, sectionID)
	aw.OK()
}
func (receiver DepSecService) AddInterview(c *gin.Context) {
	aw := app.NewWrapper(c)
	var req model.InterviewDTO
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	var sectionID *int
	sid := util.ExtractSectionID(c)
	sectionID = &sid
	interview := interviewModel.CreateInterview(req, sectionID)
	aw.Success(interview)
}
func (receiver DepSecService) EditInterview(c *gin.Context) {
	aw := app.NewWrapper(c)
	var req model.InterviewDTO
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		aw.Error(err.Error())
		return
	}
	var sectionID *int
	sid := util.ExtractSectionID(c)
	sectionID = &sid
	interview := interviewModel.UpdateInterview(req, id, sectionID)
	aw.Success(interview)
}
func (receiver DepSecService) DeleteInterview(c *gin.Context) {
	aw := app.NewWrapper(c)
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		aw.Error(err.Error())
		return
	}
	var sectionID *int
	sid := util.ExtractSectionID(c)
	sectionID = &sid
	interviewModel.DeleteInterview(id, sectionID)
	aw.OK()
}
func (receiver DepSecService) ListInterviews(c *gin.Context) {
	aw := app.NewWrapper(c)
	var sectionID *int
	sid := util.ExtractSectionID(c)
	sectionID = &sid
	interviews := interviewModel.FindInterviewsBySectionID(sectionID)
	aw.Success(interviews)
}
