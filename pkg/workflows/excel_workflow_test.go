package workflow

import (
	"fmt"
	"go-shopping-list/pkg/common"
	"go-shopping-list/pkg/recipe"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	excelize "github.com/xuri/excelize/v2"
)

type excelTest struct {
	suite.Suite
	mockScreen      *common.MockScreenInterface
	mockFileReader  *recipe.MockFileReader
	mockWorkflow    *common.MockWorkflowInterface
	mockFileChecker *mockFileChecker
	mockExcel       *mockExcel
}

func (g *excelTest) SetupTest() {
	g.mockScreen = new(common.MockScreenInterface)
	g.mockFileReader = new(recipe.MockFileReader)
	g.mockWorkflow = new(common.MockWorkflowInterface)
	g.mockFileChecker = new(mockFileChecker)
	g.mockExcel = new(mockExcel)
}

func Test_exceTest(t *testing.T) {
	suite.Run(t, new(excelTest))
}

func (g *excelTest) Test_terminal_RunReminder_Pass() {
	m := ExcelWorkflow{}
	log.Println("excel test")
	ing := recipe.Ingredient{
		UnitSize:       "WATERMELON",
		UnitType:       "CHERRY",
		IngredientName: "PITAYA",
	}

	err := m.RunReminder(g.mockScreen, ing)
	g.Nil(err)
}

func (g *excelTest) Test_createExcelSheet_Pass() {
	ing := recipe.Ingredient{
		UnitSize:       "WATERMELON",
		UnitType:       "CHERRY",
		IngredientName: "PITAYA",
	}
	g.mockExcel.On("NewFile").Return(nil)
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, titleColRow, titleVal).Return(nil)
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, dateColRow, "1998-04-26").Return(nil)
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, "B5", ing.String()).Return(nil)
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockExcel.On("SaveAs", (*excelize.File)(nil), "excel-lists/1998-04-26.xlsx").Return(nil)
	err := createExcelSheet(g.mockScreen, g.mockExcel, []recipe.Ingredient{ing}, "1998-04-26")
	g.Nil(err)
}

func (g *excelTest) Test_createExcelSheet_saveAs_Error() {
	ing := recipe.Ingredient{
		UnitSize:       "WATERMELON",
		UnitType:       "CHERRY",
		IngredientName: "PITAYA",
	}
	g.mockExcel.On("NewFile").Return(nil)
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, titleColRow, titleVal).Return(nil)
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, dateColRow, "1998-04-26").Return(nil)
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, "B5", ing.String()).Return(nil)
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockExcel.On("SaveAs", (*excelize.File)(nil), "excel-lists/1998-04-26.xlsx").Return(fmt.Errorf("save error"))
	err := createExcelSheet(g.mockScreen, g.mockExcel, []recipe.Ingredient{ing}, "1998-04-26")
	g.EqualError(err, "save error")
}

func (g *excelTest) Test_createExcelSheet_setIng_Error() {
	ing := recipe.Ingredient{
		UnitSize:       "WATERMELON",
		UnitType:       "CHERRY",
		IngredientName: "PITAYA",
	}
	g.mockExcel.On("NewFile").Return(nil)
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, titleColRow, titleVal).Return(nil)
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, dateColRow, "1998-04-26").Return(nil)
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, "B5", ing.String()).Return(fmt.Errorf("ing add error"))
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockExcel.On("SaveAs", (*excelize.File)(nil), "excel-lists/1998-04-26.xlsx").Return(nil)
	err := createExcelSheet(g.mockScreen, g.mockExcel, []recipe.Ingredient{ing}, "1998-04-26")
	g.EqualError(err, "ing add error")
}

func (g *excelTest) Test_createExcelSheet_SetDate_Error() {
	ing := recipe.Ingredient{
		UnitSize:       "WATERMELON",
		UnitType:       "CHERRY",
		IngredientName: "PITAYA",
	}
	g.mockExcel.On("NewFile").Return(nil)
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, titleColRow, titleVal).Return(nil)
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, dateColRow, "1998-04-26").Return(fmt.Errorf("date error"))
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, "B5", ing.String()).Return(nil)
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockExcel.On("SaveAs", (*excelize.File)(nil), "excel-lists/1998-04-26.xlsx").Return(nil)
	err := createExcelSheet(g.mockScreen, g.mockExcel, []recipe.Ingredient{ing}, "1998-04-26")
	g.EqualError(err, "date error")
}

func (g *excelTest) Test_createExcelSheet_SetTitle_Error() {
	ing := recipe.Ingredient{
		UnitSize:       "WATERMELON",
		UnitType:       "CHERRY",
		IngredientName: "PITAYA",
	}
	g.mockExcel.On("NewFile").Return(nil)
	g.mockExcel.On("SetCellValue", (*excelize.File)(nil), sheetName, titleColRow, titleVal).Return(fmt.Errorf("title error"))
	err := createExcelSheet(g.mockScreen, g.mockExcel, []recipe.Ingredient{ing}, "1998-04-26")
	g.EqualError(err, "title error")
}
