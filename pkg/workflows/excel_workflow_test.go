package workflow

import (
	"fmt"
	"go-shopping-list/pkg/common"
	"go-shopping-list/pkg/fruit"
	"go-shopping-list/pkg/recipe"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	excelize "github.com/xuri/excelize/v2"
)

const (
	newFileFunc      = "newFile"
	setCellValueFunc = "setCellValue"
	saveAsFunc       = "saveAs"
)

type excelTest struct {
	suite.Suite
	mockScreen      *common.MockScreenInterface
	mockFileReader  *recipe.MockFileReader
	mockWorkflow    *common.MockWorkflowInterface
	mockFileChecker *mockFileChecker
	mockExcel       *mockExcel

	ing          recipe.Ingredient
	sampleDate   string
	firstIngCell string
	filePath     string
}

func (g *excelTest) SetupTest() {
	g.mockScreen = new(common.MockScreenInterface)
	g.mockFileReader = new(recipe.MockFileReader)
	g.mockWorkflow = new(common.MockWorkflowInterface)
	g.mockFileChecker = new(mockFileChecker)
	g.mockExcel = new(mockExcel)

	g.ing = recipe.Ingredient{
		UnitSize:       fruit.Watermelon,
		UnitType:       fruit.Cherry,
		IngredientName: fruit.Pomegranate,
	}
	g.sampleDate = "1998-04-26"
	g.firstIngCell = "B5"
	g.filePath = "excel-lists/1998-04-26.xlsx"
}

func Test_exceTest(t *testing.T) {
	suite.Run(t, new(excelTest))
}

func (g *excelTest) Test_terminal_RunReminder_Pass() {
	m := excelWorkflow{}
	log.Println("excel test")

	err := m.RunReminder(g.mockScreen, g.ing)
	g.Nil(err)
}

func (g *excelTest) Test_createExcelSheet_Pass() {
	g.mockExcel.On(newFileFunc).Return(nil)
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, titleColRow, titleVal).Return(nil)
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, dateColRow, g.sampleDate).Return(nil)
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, g.firstIngCell, g.ing.String()).Return(nil)
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockExcel.On(saveAsFunc, (*excelize.File)(nil), g.filePath).Return(nil)
	err := createExcelSheet(g.mockScreen, g.mockExcel, []recipe.Ingredient{g.ing}, g.sampleDate)
	g.Nil(err)
}

func (g *excelTest) Test_createExcelSheet_saveAs_Error() {
	g.mockExcel.On(newFileFunc).Return(nil)
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, titleColRow, titleVal).Return(nil)
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, dateColRow, g.sampleDate).Return(nil)
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, g.firstIngCell, g.ing.String()).Return(nil)
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockExcel.On(saveAsFunc, (*excelize.File)(nil), g.filePath).Return(fmt.Errorf("save error"))
	err := createExcelSheet(g.mockScreen, g.mockExcel, []recipe.Ingredient{g.ing}, g.sampleDate)
	g.EqualError(err, "save error")
}

func (g *excelTest) Test_createExcelSheet_setIng_Error() {
	g.mockExcel.On(newFileFunc).Return(nil)
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, titleColRow, titleVal).Return(nil)
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, dateColRow, g.sampleDate).Return(nil)
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, g.firstIngCell, g.ing.String()).Return(fmt.Errorf("ing add error"))
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockExcel.On(saveAsFunc, (*excelize.File)(nil), g.filePath).Return(nil)
	err := createExcelSheet(g.mockScreen, g.mockExcel, []recipe.Ingredient{g.ing}, g.sampleDate)
	g.EqualError(err, "ing add error")
}

func (g *excelTest) Test_createExcelSheet_SetDate_Error() {
	g.mockExcel.On(newFileFunc).Return(nil)
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, titleColRow, titleVal).Return(nil)
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, dateColRow, g.sampleDate).Return(fmt.Errorf("date error"))
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, g.firstIngCell, g.ing.String()).Return(nil)
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockExcel.On(saveAsFunc, (*excelize.File)(nil), g.filePath).Return(nil)
	err := createExcelSheet(g.mockScreen, g.mockExcel, []recipe.Ingredient{g.ing}, g.sampleDate)
	g.EqualError(err, "date error")
}

func (g *excelTest) Test_createExcelSheet_SetTitle_Error() {
	g.mockExcel.On(newFileFunc).Return(nil)
	g.mockExcel.On(setCellValueFunc, (*excelize.File)(nil), sheetName, titleColRow, titleVal).Return(fmt.Errorf("title error"))
	err := createExcelSheet(g.mockScreen, g.mockExcel, []recipe.Ingredient{g.ing}, g.sampleDate)
	g.EqualError(err, "title error")
}
