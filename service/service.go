package service

import "sduonline-recruitment/model"

var userModel model.UserModel
var depSecModel model.DepSecModel
var aplModel model.AplModel
var qnsAnsModel model.QnsAnsModel
var configModel model.ConfigModel
var interviewModel model.InterviewModel
var scoreModel model.ScoreModel

func Setup() {
	userModel = model.UserModel{AbstractModel: model.AbstractModel{Tx: model.DB}}
	depSecModel = model.DepSecModel{AbstractModel: model.AbstractModel{Tx: model.DB}}
	aplModel = model.AplModel{AbstractModel: model.AbstractModel{Tx: model.DB}}
	qnsAnsModel = model.QnsAnsModel{AbstractModel: model.AbstractModel{Tx: model.DB}}
	configModel = model.ConfigModel{AbstractModel: model.AbstractModel{Tx: model.DB}}
	interviewModel = model.InterviewModel{AbstractModel: model.AbstractModel{Tx: model.DB}}
	scoreModel = model.ScoreModel{AbstractModel: model.AbstractModel{Tx: model.DB}}
}
