package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func LoadUnitSizeScreen(w fyne.Window) {
	// epic := widget.NewButtonWithIcon("epic", fyne.NewMenuItemSeparator().Icon, func() {
	// 	log.Println("Epic")
	// })

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	content := container.NewVBox(input, widget.NewButton("Save", func() {
		log.Println("Content was:", input.Text)
	}))

	listGrid := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, recipeListHeight)), content)

	w.SetContent(listGrid)
}
