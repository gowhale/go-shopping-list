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
	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(600, 1150)), recipeList)
	gridTop := container.New(layout.NewGridWrapLayout(fyne.NewSize(600, 50)), label, p)
	masterGrid := container.New(layout.NewVBoxLayout(), gridTop, grid)

	// Set Window and execute
	myWindow.SetFixedSize(true)
	myWindow.Resize(fyne.Size{Width: 600, Height: 1200})
	myWindow.SetContent(masterGrid)

	return myWindow
}

func itemClicked(r recipe.Recipe, p *widget.ProgressBar, l *widget.Label) {
	err := recipe.AddIngredientsToReminders(r, p, l)
	if err != nil {
		log.Printf("error whilst adding ingredients to reminds err=%e", err)
	}
}
