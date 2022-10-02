package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func nextUnitSize(w fyne.Window) *widget.Button {
	b := widget.NewButtonWithIcon("Next", fyne.NewMenuItemSeparator().Icon, func() {
		log.Println("Next button pressed...")
		LoadUnitSizeScreen(w)
	})
	b.Importance = widget.HighImportance
	return b
}

func LoadUnitTypeScreen(w fyne.Window) {
	// epic := widget.NewButtonWithIcon("epic", fyne.NewMenuItemSeparator().Icon, func() {
	// 	log.Println("Epic")
	// })

	nextB := nextUnitSize(w)

	list, _ := createNewListOfUnitTypes()
	listGrid := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, recipeListHeight)), list, nextB)

	w.SetContent(listGrid)
}

func createNewListOfUnitTypes() (*widget.List, map[string]bool) {
	unitTypes := []string{"g", "kg"}
	unitTypes = append(unitTypes, "Add a new unit type...")
	selectedIngredient := map[string]bool{}
	for _, r := range unitTypes {
		selectedIngredient[r] = false
	}

	var l *widget.List
	l = widget.NewList(
		func() int {
			return len(unitTypes)
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
			recipeCheckBox.Checked = selectedIngredient[unitTypes[li]]
			recipeCheckBox.OnChanged = func(b bool) {
				for _, r := range unitTypes {
					selectedIngredient[r] = false
				}
				selectedIngredient[unitTypes[li]] = b
				l.Refresh()
			}
			recipeCheckBox.Refresh()

			// Update label
			listContainer.Objects[labelIndex].(*widget.Label).SetText(unitTypes[li])
		})
	return l, selectedIngredient
}
