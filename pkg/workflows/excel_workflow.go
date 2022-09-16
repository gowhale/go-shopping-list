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

	f := excelize.NewFile()

	f.SetCellValue(sheetName, titleColRow, "INGREDIENTS TO BUY")
	f.SetCellValue(sheetName, dateColRow, dateString)

	ingAdded := 0
	row := ingsRowStart
	for _, ing := range ings {
		log.Printf("ingredient=%s status=IN PROGRESS", ing.String())
		cellLocation := fmt.Sprintf("%s%d", ingsColStart, row)
		cellValue := ing.String()
		log.Printf("loc=%s val=%s", cellLocation, cellValue)
		f.SetCellValue(sheetName, cellLocation, cellValue)
		row++
		ingAdded++
		progress := float64(ingAdded) / float64(len(ings))
		s.UpdateProgessBar(progress)
		log.Printf("ingredient=%s status=DONE progress=%.2f", ing.String(), progress)
	}
	listName := fmt.Sprintf("%s/%s.xlsx", listFolder, dateString)
	return f.SaveAs(listName)
}

func (*ExcelWorkflow) RunReminder(_ common.ScreenInterface, currentIng recipe.Ingredient) error {
	return nil
}
