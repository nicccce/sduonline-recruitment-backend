package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"html"
	"sduonline-recruitment/middleware"
	"sduonline-recruitment/model"
	"sduonline-recruitment/pkg/app"
	"sduonline-recruitment/pkg/util"
)

type AplService struct {
}

func (receiver AplService) SubmitApl(c *gin.Context) {
	aw := app.NewWrapper(c)
	uc := util.ExtractUserClaims(c)
	user, _ := userModel.FindUserByID(uc.UserID)
	if user.RealName == "" {
		aw.Error("请先完善个人信息")
		return
	}
	type Req struct {
		SectionID []int `form:"section_id[]" binding:"required,gte=1,lte=2"`
	}
	var req Req
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	if !depSecModel.ExistSectionsByID(req.SectionID) {
		aw.Error("部门ID不存在或重复")
		return
	}
	applications := aplModel.UpdateUserApl(uc.UserID, req.SectionID)
	applicationsVO, err := receiver.fillInScoresToApplications(&applications)
	if err != nil {
		aw.Error(err.Error())
		return
	}
	aw.Success(applicationsVO)
}
func (receiver AplService) GetQuestions(c *gin.Context) {
	aw := app.NewWrapper(c)
	uc := util.ExtractUserClaims(c)
	if len(aplModel.FindAplsByUserID(uc.UserID)) == 0 {
		aw.Error("请先报名")
		return
	}
	qns := qnsAnsModel.FindUserQuestions(uc.UserID)
	aw.Success(qns)
}
func (receiver AplService) SubmitAnswers(c *gin.Context) {
	aw := app.NewWrapper(c)
	uc := util.ExtractUserClaims(c)
	if len(aplModel.FindAplsByUserID(uc.UserID)) == 0 {
		aw.Error("请先报名")
		return
	}
	qns := qnsAnsModel.FindUserQuestions(uc.UserID)
	type AnswerReq struct {
		QuestionID int    `json:"question_id" binding:"required"`
		Answer     string `json:"answer" binding:"required"`
	}
	var ansReq []AnswerReq
	if err := c.ShouldBindJSON(&ansReq); err != nil {
		aw.Error(err.Error())
		return
	}
	var answeredQuestionIDs []int
	for _, item := range ansReq {
		if !util.IntInSlice(answeredQuestionIDs, item.QuestionID) {
			answeredQuestionIDs = append(answeredQuestionIDs, item.QuestionID)
		} else {
			aw.Error("回答问题不能重复")
			return
		}
	}
	for _, item := range qns {
		if !util.IntInSlice(answeredQuestionIDs, item.ID) {
			aw.Error("请回答所有问题")
			return
		}
	}
	if len(answeredQuestionIDs) != len(qns) {
		aw.Error("实际回答问题数量和需回答问题数量不相等")
		return
	}
	var answers []model.Answer
	for _, item := range ansReq {
		answers = append(answers, model.Answer{
			UserID:     uc.UserID,
			QuestionID: item.QuestionID,
			Text:       html.EscapeString(item.Answer),
		})
	}
	tx := model.DB.Begin()
	defer func(tx *gorm.DB) {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}(tx)
	tmpQnsAnsModel := model.QnsAnsModel{AbstractModel: model.AbstractModel{Tx: tx}}
	tmpQnsAnsModel.DeleteAnswersByUserID(uc.UserID)
	tmpQnsAnsModel.CreateAnswers(answers)
	tx.Commit()
	aw.Success(answers)
}
func (receiver AplService) ListAllApl(c *gin.Context) {
	aw := app.NewWrapper(c)
	apl := aplModel.FindAllApls()
	if err := receiver.fillInScoresToAplListVOs(&apl); err != nil {
		aw.Error(err.Error())
		return
	}
	if c.Query("dump") != "" {
		receiver.exportAplList(aw, apl)
	} else {
		aw.Success(apl)
	}
}
func (receiver AplService) ListSectionApl(c *gin.Context) {
	aw := app.NewWrapper(c)
	sectionID := util.ExtractSectionID(c)
	apl := aplModel.FindAplsBySectionID(sectionID)
	if err := receiver.fillInScoresToAplListVOs(&apl); err != nil {
		aw.Error(err.Error())
		return
	}
	if c.Query("dump") != "" {
		receiver.exportAplList(aw, apl)
	} else {
		aw.Success(apl)
	}
}
func (receiver AplService) exportAplList(aw *app.Wrapper, aplList []model.AplListVO) {
	var dto []AplXlsxRow
	sectionID := 0
	var interviews []model.Interview
	for _, item := range aplList {
		statusStr := ""
		switch item.Status {
		case 0:
			statusStr = "未面"
		case 1:
			statusStr = "通过"
		case 2:
			statusStr = "拒绝"
		}
		apl := AplXlsxRow{
			SectionName: item.SectionName,
			RealName:    item.RealName,
			StudentID:   item.StudentID,
			Qq:          item.Qq,
			Phone:       item.Phone,
			Status:      statusStr,
			Note:        item.Note,
			Scores:      item.Score,
		}
		if sectionID != item.SectionID {
			sectionID = item.SectionID
			interviews = interviewModel.FindInterviewsBySectionID(&sectionID)
		}
		var err error
		apl.Score, err = receiver.calculateWeightedScore(&item.Score, &interviews)
		if err != nil {
			aw.Error(err.Error())
			return
		}
		dto = append(dto, apl)

	}
	byt, err := ExportAplList(dto, &interviews)
	if err != nil {
		aw.Error(err.Error())
		return
	}
	aw.Ctx.Header("Content-Description", "File Transfer")
	aw.Ctx.Header("Content-Transfer-Encoding", "binary")
	aw.Ctx.Header("Content-Disposition", "attachment; filename=dump.xlsx")
	aw.Ctx.Header("Content-Type", "application/octet-stream")
	aw.Ctx.Data(200, "application/octet-stream", byt)
}
func (receiver AplService) calculateWeightedScore(scores *[]int, interviews *[]model.Interview) (int, error) {
	if len(*scores) != len(*interviews) {
		return 0, errors.New("inner error")
	}
	scoreSum := 0
	weightSum := 0
	for i, score := range *scores {
		if score >= 0 {
			scoreSum += int(int64(score) * ((*interviews)[i].Weight.Int64))
			weightSum += int((*interviews)[i].Weight.Int64)
		}
	}
	if weightSum == 0 {
		return 0, nil
	}
	weightedScore := scoreSum / weightSum
	return weightedScore, nil
}
func (receiver AplService) ExportAplsByPermission(c *gin.Context) {
	aw := app.NewWrapper(c)
	uc := util.ExtractUserClaims(c)
	_, perms := userModel.FindUserByID(uc.UserID)
	var secIDs []int
	for _, perm := range perms {
		secIDs = append(secIDs, perm.SectionID)
	}
	aplList := aplModel.FindAplsBySectionIDs(secIDs)
	err := receiver.fillInScoresToAplListVOs(&aplList)
	if err != nil {
		aw.Error(err.Error())
		return
	}
	receiver.exportAplList(aw, aplList)
}
func (receiver AplService) GetAplAnswers(c *gin.Context) {
	aw := app.NewWrapper(c)
	type GetAplAnsReq struct {
		middleware.SectionReq
		AplID int `form:"apl_id" binding:"required"`
	}
	var req GetAplAnsReq
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	apl := aplModel.FindAplByID(req.AplID)
	if apl == nil {
		aw.Error("apl_id不存在")
		return
	}
	aplAnsVO := aplModel.FindAplAnswersBySectionIDAplID(req.SectionID, req.AplID)
	user, _ := userModel.FindUserByID(apl.UserID)
	aw.Success(gin.H{
		"apl":     apl,
		"answers": aplAnsVO,
		"user":    user})
}
func (receiver AplService) JudgeApl(c *gin.Context) {
	aw := app.NewWrapper(c)
	type JudgeReq struct {
		model.ApplicationJudge
		middleware.SectionReq
		AplID  int   `form:"apl_id" binding:"required"`
		Scores []int `form:"score" binding:"required"`
	}
	var req JudgeReq
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	userID, err := aplModel.GetUserIDByAplID(req.AplID)
	if err != nil {
		aw.Error(err.Error())
		return
	}
	aplModel.UpdateAplJudgeByIDSectionID(&req.ApplicationJudge, req.AplID, req.SectionID)
	if err := scoreModel.UpdateScores(req.Scores, &req.SectionID, &userID); err != nil {
		aw.Error(err.Error())
		return
	}
	aplVO := model.ApplicationJudgeVO{
		Status: req.Status,
		Score:  req.Scores,
		Note:   req.Note,
	}
	aw.Success(aplVO)
}
func (receiver AplService) GetUserAnswers(c *gin.Context) {
	aw := app.NewWrapper(c)
	uc := util.ExtractUserClaims(c)
	vo := qnsAnsModel.FindQuestionsAnswersByUserID(uc.UserID)
	aw.Success(vo)
}
func (receiver AplService) ListUserApls(c *gin.Context) {
	aw := app.NewWrapper(c)
	type ListUserAplReq struct {
		UserID int `form:"user_id" binding:"required"`
	}
	var req ListUserAplReq
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	apls := aplModel.FindAplsByUserID(req.UserID)
	if err := receiver.fillInScoresToAplListVOs(&apls); err != nil {
		aw.Error(err.Error())
		return
	}
	aw.Success(apls)
}
func (receiver AplService) fillInScoresToApplications(applications *[]model.Application) ([]model.ApplicationVO, error) {
	applicationsVO := make([]model.ApplicationVO, 0)
	for _, application := range *applications {
		applicationVO := model.ApplicationVO{}
		err := receiver.fillInSingleScore(&application, &applicationVO)
		if err != nil {
			return nil, err
		}
		applicationsVO = append(applicationsVO, applicationVO)
	}
	return applicationsVO, nil
}
func (receiver AplService) fillInScoresToAplListVOs(aplList *[]model.AplListVO) error {
	for i, apl := range *aplList {
		scores, err := scoreModel.FindScoresBySectionIDUserID(&apl.SectionID, &apl.UserID)
		if err != nil {
			return err
		}
		(*aplList)[i].Score = scores
	}
	return nil
}
func (receiver AplService) fillInSingleScore(application *model.Application, applicationVO *model.ApplicationVO) error {
	applicationVO.ID = application.ID
	applicationVO.UserID = application.UserID
	applicationVO.SectionID = application.SectionID
	applicationVO.CreatedAt = application.CreatedAt
	applicationVO.UpdatedAt = application.UpdatedAt
	applicationVO.ApplicationJudgeVO.Status = application.Status
	applicationVO.ApplicationJudgeVO.Note = application.Note

	scores, err := scoreModel.FindScoresBySectionIDUserID(&application.SectionID, &application.UserID)
	if err != nil {
		return err
	}

	applicationVO.ApplicationJudgeVO.Score = scores

	return nil
}
