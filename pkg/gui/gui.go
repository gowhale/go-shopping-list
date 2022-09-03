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
	screenHeight      = 1200
	recipeListHeight  = 1050
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

// NewApp returns a fyne.Window
func NewApp(recipes []recipe.Recipe, wf workflowInterface) fyne.Window {
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
	recipeList := createNewListOfRecipes(s, &recipe.FileInteractionImpl{}, wf, recipes)

	// Create content grid
	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, recipeListHeight)), recipeList)
	gridTop := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, progressBarHeight)), label, p)
	masterGrid := container.New(layout.NewVBoxLayout(), gridTop, grid)

	// Set Window and execute
	myWindow.SetFixedSize(true)
	myWindow.Resize(fyne.Size{Width: screenWidth, Height: screenHeight})
	myWindow.SetContent(masterGrid)

	return myWindow
}

func createNewListOfRecipes(s screenInterface, f recipe.FileReader, w workflowInterface, recipes []recipe.Recipe) *widget.List {
	// Recipe list with all recipes
	return widget.NewList(
		func() int {
			return len(recipes)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("template", func() {})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Button).SetText(recipes[i].Name)
			o.(*widget.Button).OnTapped = func() {
				itemClicked(s, recipes[i], f, w)
			}
		})
}

func itemClicked(s screenInterface, r recipe.Recipe, f recipe.FileReader, w workflowInterface) {
	err := addIngredientsToReminders(r, s, f, w)
	if err != nil {
		log.Printf("error whilst adding ingredients to reminds err=%e", err)
	}
}
