package util

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/tealeg/xlsx/v3"
)

type AplXlsxRow struct {
	SectionName string
	RealName    string
	StudentID   string
	Qq          string
	Phone       string
	Status      string
	Score       int
	ScoreExtra  int
	Note        string
}

var AplXlsxHeader = []string{"部门", "姓名", "学号", "QQ", "手机", "状态", "一轮评分", "二轮评分", "备注"}

func ExportAplList(aplList []AplXlsxRow) ([]byte, error) {
	workbook := xlsx.NewFile()
	if len(aplList) == 0 {
		return nil, errors.New("没有表格数据")
	}
	currentSectionName := aplList[0].SectionName
	sheet, err := workbook.AddSheet(currentSectionName)
	if err != nil {
		return nil, err
	}
	defer sheet.Close()
	headerRow := sheet.AddRow()
	for _, item := range AplXlsxHeader {
		cell := headerRow.AddCell()
		cell.SetString(item)
	}
	for _, apl := range aplList {
		if apl.SectionName != currentSectionName {
			sh, err := workbook.AddSheet(apl.SectionName)
			if err != nil {
				return nil, err
			}
			currentSectionName = apl.SectionName
			sheet = sh
		}
		row := sheet.AddRow()
		row.AddCell().SetString(apl.SectionName)
		row.AddCell().SetString(apl.RealName)
		row.AddCell().SetString(apl.StudentID)
		row.AddCell().SetString(apl.Qq)
		row.AddCell().SetString(apl.Phone)
		row.AddCell().SetString(apl.Status)
		row.AddCell().SetInt(apl.Score)
		row.AddCell().SetInt(apl.ScoreExtra)
		row.AddCell().SetString(apl.Note)
	}
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err = workbook.Write(writer)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
