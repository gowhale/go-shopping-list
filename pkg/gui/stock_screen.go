package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// NewApp returns a fyne.Window
func NewStockScreen(ings []string) fyne.Window {
	myApp := app.New()
	myWindow := myApp.NewWindow("List Widget")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	list, _ := createNewListOfIngs(ings)

	nextB := createNextButton(myWindow)
	previousB := createBackButton()

	listGrid := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, recipeListHeight)), list)

	masterGrid := container.New(layout.NewVBoxLayout(), listGrid, nextB, previousB)

	myWindow.Resize(fyne.Size{Width: screenWidth, Height: screenHeight})
	myWindow.SetContent(masterGrid)

	return myWindow
}

func createBackButton() *widget.Button {
	b := widget.NewButtonWithIcon("Back", fyne.NewMenuItemSeparator().Icon, func() {
		log.Println("Previous button pressed...")
	})
	b.Importance = widget.HighImportance
	return b
}

func createNextButton(w fyne.Window) *widget.Button {
	b := widget.NewButtonWithIcon("Next", fyne.NewMenuItemSeparator().Icon, func() {
		log.Println("Next button pressed...")
		LoadUnitTypeScreen(w)
	})
	b.Importance = widget.HighImportance
	return b
}

func createNewListOfIngs(ingStr []string) (*widget.List, map[string]bool) {
	ingStr = append(ingStr, "Add a new ingredient...")
	selectedIngredient := map[string]bool{}
	for _, r := range ingStr {
		selectedIngredient[r] = false
	}

	var l *widget.List
	l = widget.NewList(
		func() int {
			return len(ingStr)
		},
		func() fyne.CanvasObject {
			hbox := container.NewGridWithColumns(listContainerCols)
			hbox1 := container.NewGridWithColumns(listCols)
			hbox1.Add(widget.NewCheck("", func(bool) {}))
			hbox1.Add(widget.NewLabel("table"))
			hbox.Add(hbox1)
			return hbox
		},
		func(li widget.ListItemID, o fyne.CanvasObject) {
			// Update Checkbox
			listContainer := o.(*fyne.Container).Objects[listBoxIndex].(*fyne.Container)
			recipeCheckBox := listContainer.Objects[checkIndex].(*widget.Check)
			recipeCheckBox.Checked = selectedIngredient[ingStr[li]]
			recipeCheckBox.OnChanged = func(b bool) {
				for _, r := range ingStr {
					selectedIngredient[r] = false
				}
				selectedIngredient[ingStr[li]] = b
				l.Refresh()
			}
			recipeCheckBox.Refresh()

			// Update label
			listContainer.Objects[labelIndex].(*widget.Label).SetText(ingStr[li])
		})
	return l, selectedIngredient
}
