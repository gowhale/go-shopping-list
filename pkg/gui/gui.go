package gui

import (
	"go-shopping-list/pkg/recipe"
	"go-shopping-list/pkg/workflows"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

func NewApp(recipes []recipe.Recipe) fyne.Window {

	myApp := app.New()
	myWindow := myApp.NewWindow("List Widget")

	label := widget.NewLabel("Click a recipe to add the ingredients...")

	// Progress bar for adding ings
	p := widget.NewProgressBar()

	// Recipe list with all recipes
	recipeList := widget.NewList(
		func() int {
			return len(recipes)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("template", func() {})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Button).SetText(recipes[i].Name)
			o.(*widget.Button).OnTapped = func() { itemClicked(recipes[i], p, label) }
		})

	// Create content grid
	grid := container.NewAdaptiveGrid(1, recipeList)
	gridTop := container.NewAdaptiveGrid(1, label, p)
	masterGrid := container.NewVBox(gridTop, grid)

	// Set Window and execute
	myWindow.SetFixedSize(true)
	myWindow.Resize(fyne.Size{Width: 600, Height: 1200})
	myWindow.SetContent(masterGrid)

	return myWindow
}

func itemClicked(r recipe.Recipe, p *widget.ProgressBar, l *widget.Label) {
	err := workflows.AddIngredientsToReminders(r, p, l)
	if err != nil {
		log.Printf("error whilst adding ingredients to reminds err=%e", err)
	}
}
