package gui

import (
	"go-shopping-list/pkg/common"
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

func (s *screen) UpdateProgessBar(percent float64) {
	s.p.SetValue(percent)
	s.p.Refresh()
}

func (s *screen) UpdateLabel(msg string) {
	s.l.SetText(msg)
	s.l.Refresh()
}

func createSubmitButton(s common.ScreenInterface, wf common.WorkflowInterface, fr recipe.FileReader, recipes *[]string, recipeMap map[string]recipe.Recipe) *widget.Button {
	return widget.NewButton("Add To Shopping List", func() {
		err := wf.SubmitShoppingList(s, wf, fr, *recipes, recipeMap)
		if err != nil {
			log.Println(err)
		}
	})
}

// NewApp returns a fyne.Window
func NewApp(recipes []recipe.Recipe, recipeMap map[string]recipe.Recipe, wf common.WorkflowInterface) fyne.Window {
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
