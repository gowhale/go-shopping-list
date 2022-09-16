package workflow

import (
	"fmt"
	"go-shopping-list/pkg/common"
	"go-shopping-list/pkg/recipe"
	"log"
	"time"

	excelize "github.com/xuri/excelize/v2"
)

const (
	titleColRow  = "B2"
	titleVal     = "INGREDIENTS TO BUY"
	dateColRow   = "B3"
	ingsRowStart = 5
	ingsColStart = "B"
	listFolder   = "excel-lists"
	sheetName    = "Sheet1"
)

// ExcelWorkflow will create an excel sheet with ingredients
type ExcelWorkflow struct{}

func (*ExcelWorkflow) SubmitShoppingList(s common.ScreenInterface, wf common.WorkflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error {
	return SubmitShoppingList(s, wf, fr, recipes, recipeMap)
}

func (*ExcelWorkflow) AddIngredientsToReminders(ings []recipe.Ingredient, s common.ScreenInterface, w common.WorkflowInterface) error {
	year, month, day := time.Now().Date()
	dateString := fmt.Sprintf("%d-%d-%d", year, month, day)
	return createExcelSheet(s, &excelImpl{}, ings, dateString)
}

func (*ExcelWorkflow) RunReminder(_ common.ScreenInterface, currentIng recipe.Ingredient) error {
	return nil
}

//go:generate go run github.com/vektra/mockery/cmd/mockery -name excel -inpkg --filename excel_mock.go
type excel interface {
	NewFile() *excelize.File
	SetCellValue(f *excelize.File, sheet string, axis string, value interface{}) error
	SaveAs(f *excelize.File, name string, opt ...excelize.Options) error
}

type excelImpl struct{}

func (*excelImpl) NewFile() *excelize.File {
	return excelize.NewFile()
}

func (*excelImpl) SetCellValue(f *excelize.File, sheet string, axis string, value interface{}) error {
	return f.SetCellValue(sheet, axis, value)
}

func (*excelImpl) SaveAs(f *excelize.File, name string, opt ...excelize.Options) error {
	return f.SaveAs(name, opt...)
}

func createExcelSheet(s common.ScreenInterface, e excel, ings []recipe.Ingredient, dateString string) error {

	f := e.NewFile()

	if err := e.SetCellValue(f, sheetName, titleColRow, titleVal); err != nil {
		return err
	}
	if err := e.SetCellValue(f, sheetName, dateColRow, dateString); err != nil {
		return err
	}

	ingAdded := 0
	row := ingsRowStart
	for _, ing := range ings {
		log.Printf("ingredient=%s status=IN PROGRESS", ing.String())
		cellLocation := fmt.Sprintf("%s%d", ingsColStart, row)
		cellValue := ing.String()
		log.Printf("loc=%s val=%s", cellLocation, cellValue)
		if err := e.SetCellValue(f, sheetName, cellLocation, cellValue); err != nil {
			return err
		}
		row++
		ingAdded++
		progress := float64(ingAdded) / float64(len(ings))
		s.UpdateProgessBar(progress)
		log.Printf("ingredient=%s status=DONE progress=%.2f", ing.String(), progress)
	}
	listName := fmt.Sprintf("%s/%s.xlsx", listFolder, dateString)
	return e.SaveAs(f, listName)
}
