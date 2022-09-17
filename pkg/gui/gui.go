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
	screenHeight      = 500
	progressBarHeight = 50
	submitBarHeight   = 50
	recipeListHeight  = screenHeight - progressBarHeight - submitBarHeight
	progressBarEmpty  = 0.0
	progressBarFull   = 1.0

	listContainerCols = 1
	listCols          = 2
	listBoxIndex      = 0
	checkIndex        = 0
	labelIndex        = 1
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

func createSubmitButton(s common.ScreenInterface, wf common.WorkflowInterface, fr recipe.FileReader, recipes map[string]bool, recipeMap map[string]recipe.Recipe) *widget.Button {
	b := widget.NewButtonWithIcon("Add To Shopping List", fyne.NewMenuItemSeparator().Icon, func() {
		selectedRecipes := []string{}
		for k, v := range recipes {
			if v {
				selectedRecipes = append(selectedRecipes, k)
			}
		}
		err := wf.SubmitShoppingList(s, wf, fr, selectedRecipes, recipeMap)
		if err != nil {
			log.Println(err)
		}
	})
	b.Importance = widget.HighImportance
	return b
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
	selectedRecipe := map[string]bool{}
	recipeList, selectedRecipe := createNewListOfRecipes(selectedRecipe, recipesAsStrings)

	submit := createSubmitButton(s, wf, fr, selectedRecipe, recipeMap)
	gridTop := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, progressBarHeight)), label, p)
	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, recipeListHeight)), recipeList)
	gridBottum := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, submitBarHeight)), submit)
	masterGrid := container.New(layout.NewVBoxLayout(), gridTop, grid, gridBottum)

	myWindow.Resize(fyne.Size{Width: screenWidth, Height: screenHeight})
	myWindow.SetContent(masterGrid)

	return myWindow
}

func createNewListOfRecipes(selectedRecipe map[string]bool, recipesStr []string) (*widget.List, map[string]bool) {
	for _, r := range recipesStr {
		selectedRecipe[r] = false
	}
	l := widget.NewList(
		func() int {
			return len(recipesStr)
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
			recipeCheckBox.Checked = selectedRecipe[recipesStr[li]]
			recipeCheckBox.OnChanged = func(b bool) {
				selectedRecipe[recipesStr[li]] = b
			}
			recipeCheckBox.Refresh()

			// Update label
			listContainer.Objects[labelIndex].(*widget.Label).SetText(recipesStr[li])
		})
	return l, selectedRecipe
}
