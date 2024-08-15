package model

import (
	"gorm.io/gorm"
	"sduonline-recruitment/pkg/util"
	"strconv"
)

type DepSecModel struct {
	AbstractModel
}
type Department struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Intro string `json:"intro"`
}
type Section struct {
	ID           int    `json:"id"`
	DepartmentID int    `json:"department_id"`
	Name         string `json:"name"`
	SectionInfo
}
type SectionInfo struct {
	Intro      string `json:"intro" form:"intro"`
	IntroExtra string `json:"intro_extra" form:"intro_extra"`
	Image      string `json:"image" form:"image"`
}
type depSecDTO struct {
	Section
}

func (receiver DepSecModel) FindAll() []depSecDTO {
	var depSec []depSecDTO
	err := receiver.Tx.Raw("select id,0 department_id,name,intro,null intro_extra from departments " +
		"union " +
		"select id,department_id,name,intro,intro_extra from sections").Find(&depSec).Error
	util.ForwardOrPanic(err)
	return depSec
}
func (receiver DepSecModel) ExistSectionsByID(sectionID []int) bool {
	inVar := "("
	for ix, id := range sectionID {
		inVar += strconv.Itoa(id)
		if ix == len(sectionID)-1 {
			inVar += ")"
		} else {
			inVar += ","
		}
	}
	var count int
	err := receiver.Tx.Raw("select count(distinct id) from sections where id in " + inVar).Find(&count).Error
	util.ForwardOrPanic(err)
	if count != len(sectionID) {
		return false
	} else {
		return true
	}
}
func (receiver DepSecModel) FindSectionByID(id int) *Section {
	var section Section
	err := receiver.Tx.Take(&section, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	util.ForwardOrPanic(err)
	return &section
}
func (receiver DepSecModel) UpdateSectionByID(section *Section) {
	err := receiver.Tx.Save(&section).Error
	util.ForwardOrPanic(err)
}
func (receiver DepSecModel) FindDepartmentByID(id int) *Department {
	var department Department
	err := receiver.Tx.Take(&department, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	util.ForwardOrPanic(err)
	return &department
}
func (receiver DepSecModel) UpdateDepartment(department *Department) {
	err := receiver.Tx.Save(&department).Error
	util.ForwardOrPanic(err)
}
