package gui

import (
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

//go:generate go run github.com/vektra/mockery/cmd/mockery -name Test -inpkg --filename test_mock.go
type Test interface {
	addIngredientsToReminders(ings []recipe.Ingredient, s screenInterface, w Test) error
	runReminder(s screenInterface, currentIng recipe.Ingredient) error
	submitShoppingList(s screenInterface, wf Test, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error
}

func createSubmitButton(s screenInterface, wf Test, fr recipe.FileReader, recipes *[]string, recipeMap map[string]recipe.Recipe) *widget.Button {
	return widget.NewButton("Add To Shopping List", func() {
		err := wf.submitShoppingList(s, wf, fr, *recipes, recipeMap)
		if err != nil {
			log.Fatalln(err)
		}
	})
}

// NewApp returns a fyne.Window
func NewApp(recipes []recipe.Recipe, recipeMap map[string]recipe.Recipe, wf Test) fyne.Window {
	myApp := app.New()
	myWindow := myApp.NewWindow("List Widget")

	label := widget.NewLabel("Click a recipe to add the ingredients...")

	// Progress bar for adding ings
	p := widget.NewProgressBar()

	s := &screen{
		l: label,
		p: p,
	}

	fr := &recipe.FileInteractionImpl{}

	// Recipe list with all recipes
	var recipesAsStrings []string
	for _, v := range recipes {
		recipesAsStrings = append(recipesAsStrings, v.Name)
	}
	recipeList := createNewListOfRecipes(recipesAsStrings)

	submit := createSubmitButton(s, wf, fr, &recipeList.Selected, recipeMap)
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

func createNewListOfRecipes(recipesStr []string) *widget.CheckGroup {
	// Recipe list with all recipes
	return widget.NewCheckGroup(recipesStr, nil)
}
