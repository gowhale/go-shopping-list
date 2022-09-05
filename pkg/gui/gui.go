package gui

import (
	"fmt"
	"go-shopping-list/pkg/recipe"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	screenWidth       = 600
	screenHeight      = 700
	recipeListHeight  = 650
	progressBarHeight = 50
	progressBarEmpty  = 0.0
	progressBarFull   = 1.0
	recipeFinishLabel = "Finished. Select another recipe to add more."
)

type screen struct {
	p *widget.ProgressBar
	l *widget.Label
}

//go:generate go run github.com/vektra/mockery/cmd/mockery -name screenInterface -inpkg --filename screen_mock.go
type screenInterface interface {
	updateProgessBar(float64)
	updateLabel(string)
}

func (s *screen) updateProgessBar(percent float64) {
	s.p.SetValue(percent)
	s.p.Refresh()
}

func (s *screen) updateLabel(msg string) {
	s.l.SetText(msg)
	s.l.Refresh()
}

func submitShoppingList(s screenInterface, wf workflowInterface, recipes []string, recipeMap map[string]recipe.Recipe) error {
	fmt.Println("Currently selected Recipes:")
	recipesSelected := []recipe.Recipe{}
	for _, v := range recipes {
		if r, ok := recipeMap[v]; ok {
			recipesSelected = append(recipesSelected, r)
			f := recipe.FileInteractionImpl{}
			if err := f.IncrementPopularity(r.Name); err != nil {
				return err
			}
		}

	}
	ings := recipe.CombineRecipesToIngredients(recipesSelected)
	return addIngredientsToReminders(ings, s, &recipe.FileInteractionImpl{}, wf)
}

func createSubmitButton(s screenInterface, wf workflowInterface, recipes []string, recipeMap map[string]recipe.Recipe) *widget.Button {
	return widget.NewButton("Add To Shopping List", func() {
		submitShoppingList(s, wf, recipes, recipeMap)
	})
}

// NewApp returns a fyne.Window
func NewApp(recipes []recipe.Recipe, recipeMap map[string]recipe.Recipe, wf workflowInterface) fyne.Window {
	myApp := app.New()
	myWindow := myApp.NewWindow("List Widget")

	label := widget.NewLabel("Click a recipe to add the ingredients...")

	// Progress bar for adding ings
	p := widget.NewProgressBar()

	s := &screen{
		l: label,
		p: p,
	}

	// Recipe list with all recipes
	var recipesAsStrings []string
	for _, v := range recipes {
		recipesAsStrings = append(recipesAsStrings, v.Name)
	}
	recipeList := createNewListOfRecipes(s, &recipe.FileInteractionImpl{}, wf, recipesAsStrings)

	submit := createSubmitButton(s, wf, recipesAsStrings, recipeMap)
	gridTop := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, progressBarHeight)), label, p)
	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, recipeListHeight)), recipeList)
	gridBottum := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, progressBarHeight)), submit)
	masterGrid := container.New(layout.NewVBoxLayout(), gridTop, grid, gridBottum)

	// Set Window and execute
	myWindow.SetFixedSize(true)
	myWindow.Resize(fyne.Size{Width: screenWidth, Height: screenHeight})
	myWindow.SetContent(masterGrid)

	return myWindow
}

func createCheckBoxs(recipes []recipe.Recipe) []fyne.CanvasObject {
	allChecks := []fyne.CanvasObject{}
	for _, r := range recipes {
		check := widget.NewCheck(r.Name, func(value bool) {
			log.Printf("recipt=%s value=%t", r.Name, value)
		})
		allChecks = append(allChecks, check)
	}
	return allChecks
}

func createNewListOfRecipes(s screenInterface, f recipe.FileReader, w workflowInterface, recipesStr []string) *widget.CheckGroup {
	// Recipe list with all recipes
	return widget.NewCheckGroup(recipesStr, nil)
}
