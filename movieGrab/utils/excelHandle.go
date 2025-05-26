package utils

import (
	"fmt"
	//"github.com/360EntSecGroup-Skylar/excelize"
	excelize "github.com/xuri/excelize/v2"
	"sync"
)

type ExcelHander struct {
	filePath     string
	savePath     string
	excelizeFile *excelize.File
	sheetName    string
	streamWriter *excelize.StreamWriter
	mutex        sync.Mutex
}

func CreateExcelFile(filePath string, savePath string) *ExcelHander {
	eh := ExcelHander{
		filePath: filePath,
		savePath: savePath,
	}
	return &eh
}

func (this *ExcelHander) SetSheetName(sheetName string) {
	this.sheetName = sheetName
}

// 打开excel文件并设置sheet
// param sheetName string
// return err error
func (this *ExcelHander) OpenExcelFile(sheetName string) (err error) {
	this.mutex.Lock()
	this.excelizeFile, err = excelize.OpenFile(this.filePath)
	if err != nil {
		fmt.Println("打开excel文件失败", err)
		return err
	}
	this.sheetName = sheetName

	return err
}

func (this *ExcelHander) CloseExcelFile() {
	this.excelizeFile.Close()
	this.mutex.Unlock()
}

// 以excel流方式读取文件内容
func (this *ExcelHander) ReadOrginContentsToExcel() {
	this.streamWriter, _ = this.excelizeFile.NewStreamWriter(this.sheetName)
	rows, _ := this.excelizeFile.GetRows(this.sheetName)
	cols, _ := this.excelizeFile.GetCols(this.sheetName)
	//将源文件内容先写入excel
	for rowid, row_pre := range rows {
		row_p := make([]interface{}, len(cols)+1)
		for colID_p := 0; colID_p < len(cols); colID_p++ {
			if row_pre == nil {
				row_p[colID_p] = nil
			} else {
				row_p[colID_p] = row_pre[colID_p]
			}
		}
		cell_pre, _ := excelize.CoordinatesToCellName(1, rowid+1)
		if err := this.streamWriter.SetRow(cell_pre, row_p); err != nil {
			fmt.Println(err)
		}
	}
}

// 以excel流方式追加文件内容
func (this *ExcelHander) AppendContentsToExcel(dataMap map[int]map[string]string) {
	rows, _ := this.excelizeFile.GetRows(this.sheetName)
	cols, _ := this.excelizeFile.GetCols(this.sheetName)
	//将新加contents写进流式写入器
	for i, data := range dataMap {
		row := make([]interface{}, len(cols))
		if data == nil {
			continue
		}
		row = this.sortRowContent(data)
		cell, _ := excelize.CoordinatesToCellName(1, len(rows)+i+1) //决定写入的位置
		// for _, r := range row {
		// 	fmt.Println(r.(string))
		// }
		if err := this.streamWriter.SetRow(cell, row); err != nil {
			fmt.Println(err)
		}
	}
	this.streamWriter.Flush()
	this.excelizeFile.SaveAs(this.savePath)
}

// 普通方式将文件内容写入excel
// param dataMap map[string]string
func (this *ExcelHander) WriteContentToExcel(dataMap map[string]string) {
	rows, _ := this.excelizeFile.GetRows(this.sheetName)
	length := len(rows)
	rowContents := this.sortRowContent(dataMap)
	for key, value := range rowContents {
		cell, _ := excelize.CoordinatesToCellName(key+1, length+1)
		err := this.excelizeFile.SetCellValue(this.sheetName, cell, value.(string))
		if err != nil {
			fmt.Println("title set is:", err)
		}
	}
}

func (this *ExcelHander) SaveAsExcel() (err error) {
	if err = this.excelizeFile.SaveAs(this.savePath); err != nil {
		fmt.Println("保存excel文件失败", err)
	}
	return
}

// 获取row值
// param rowData map[string]string
// return rowRes []string
func (this *ExcelHander) sortRowContent(rowData map[string]string) (rowRes []interface{}) {
	if len(rowData) == 0 {
		return nil
	}
	var titles = []string{"标　　题", "译　　名", "片　　名", "年　　代", "产　　地", "类　　别", "语　　言", "上映日期", "IMDb评分", "豆瓣评分", "片　　长", "导　　演", "编　　剧", "演　　员", "简　　介", "电影链接"}
	rowRes = make([]interface{}, len(titles))
	for i := 0; i < len(titles); i++ {
		info, ok := rowData[titles[i]]
		if ok {
			rowRes[i] = info
		} else {
			rowRes[i] = "-"
		}
	}
	return
}

// 返回列名所在的列号，例如：A1
func (this *ExcelHander) getRowNumChar(efile *excelize.File, sheetName string, rowName string) (rowNum string) {
	var flag bool
	allRow, _ := efile.GetRows(sheetName)
	for index, row := range allRow {
		flag = false
		if 0 == index {
			for k, v := range row {
				if v == rowName {
					rowNum = string(rune(65 + k))
					flag = true
					break
				}
			}
			if !flag {
				rowNum = string(rune(81))
			}
			break
		}
	}
	return
}

func (this *ExcelHander) CreateExcelTitle(sheetName string) {
	var titles = []string{"标　　题", "译　　名", "片　　名", "年　　代", "产　　地", "类　　别", "语　　言", "上映日期", "IMDb评分", "豆瓣评分", "片　　长", "导　　演", "编　　剧", "演　　员", "简　　介", "电影链接"}
	f := excelize.NewFile()
	err := f.SetSheetName("Sheet1", sheetName)
	index, err := f.GetSheetIndex("Sheet1")
	if err != nil {
		fmt.Println("修改sheet失败", err)
		return
	}
	f.SetActiveSheet(index)

	for i, t := range titles {
		f.SetCellStr(sheetName, string(65+i)+"1", t)
	}

	if err := f.SaveAs(this.filePath); err != nil {
		fmt.Println("创建电影excel模板失败:", err)
		return
	}
	// for tk, tv := range title {
	// 	err := eFile.SetCellValue(sheetName, string(rune(65+tk+1))+"1", tv)
	// 	if err != nil {
	// 		fmt.Println("title set is:", err)
	// 		break
	// 	}
	// }
}
