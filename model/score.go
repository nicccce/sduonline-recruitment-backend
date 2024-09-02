package model

import (
	"errors"
	"gorm.io/gorm"
)

type ScoreModel struct {
	AbstractModel
}
type Score struct {
	ID          int `json:"id"`
	UserID      int `json:"user_id"`
	InterviewID int `json:"interview_id"`
	Score       int `json:"score"`
}

func (receiver ScoreModel) UpdateScores(scores []int, sectionID, userID *int) error {
	var interviews []Interview
	err := receiver.Tx.Order("id ASC").Find(&interviews, "section_id=?", *sectionID).Error
	if err != nil {
		return err
	}
	if len(scores) != len(interviews) {
		return errors.New("scores does not match interviews")
	}
	for i := 0; i < len(scores); i++ {
		err := receiver.updateOrInitScore(&(interviews[i].ID), userID, &scores[i])
		if err != nil {
			return err
		}
	}
	return nil
}
func (receiver ScoreModel) updateOrInitScore(interviewID, userID, scoreNum *int) error {
	var score Score
	err := receiver.Tx.First(&score, "user_id=? AND interview_id=?", *userID, *interviewID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		score = Score{
			UserID:      *userID,
			Score:       *scoreNum,
			InterviewID: *interviewID,
		}
		err := receiver.Tx.Create(&score).Error
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	score.Score = *scoreNum
	err = receiver.Tx.Save(&score).Error
	if err != nil {
		return err
	}
	return nil
}
func (receiver ScoreModel) FindScoresBySectionIDUserID(sectionID, userID *int) ([]int, error) {
	var scoreNums []int = make([]int, 0)
	var interviews []Interview
	err := receiver.Tx.Order("id ASC").Find(&interviews, "section_id=?", *sectionID).Error
	if err != nil {
		return scoreNums, err
	}
	for i := 0; i < len(interviews); i++ {
		var score Score
		err := receiver.Tx.First(&score, "user_id=? AND interview_id=?", *userID, interviews[i].ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			scoreNums = append(scoreNums, -1)
			continue
		}
		if err != nil {
			return scoreNums, err
		}
		scoreNums = append(scoreNums, score.Score)
	}
	return scoreNums, nil
}
