package gui

import (
	"go-shopping-list/pkg/common"
	"go-shopping-list/pkg/recipe"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
)

type screen struct {
	p *widget.ProgressBar
	l *widget.Label
}

func (s *screen) UpdateProgessBar(percent float64) {
	s.p.SetValue(percent)
	// s.p.Refresh()
}

func (s *screen) UpdateLabel(msg string) {
	s.l.SetText(msg)
	// s.l.Refresh()
}

func createSubmitButton(s common.ScreenInterface, wf common.WorkflowInterface, fr recipe.FileReader, recipes binding.StringList, recipesAsStrings []string, recipeMap map[string]recipe.Recipe) *widget.Button {
	b := widget.NewButtonWithIcon("Add To Shopping List", fyne.NewMenuItemSeparator().Icon, func() {
		// err := wf.SubmitShoppingList(s, wf, fr, *recipes, recipeMap)
		// if err != nil {
		// 	log.Println(err)
		// }
		for i, r := range recipesAsStrings {
			v, err := recipes.GetValue(i)
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("recipe=%s val=%s", r, v)
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
	recipeList := createNewListOfRecipes(recipesAsStrings)

	// log.Fatalln(recipeList.Length())
	submit := createSubmitButton(s, wf, fr, nil, recipesAsStrings, recipeMap)
	gridTop := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, progressBarHeight)), label, p)
	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, recipeListHeight)), recipeList)
	gridBottum := container.New(layout.NewGridWrapLayout(fyne.NewSize(screenWidth, submitBarHeight)), submit)
	masterGrid := container.New(layout.NewVBoxLayout(), gridTop, grid, gridBottum)

	// Set Window and execute
	// myWindow.SetFixedSize(true)
	myWindow.Resize(fyne.Size{Width: screenWidth, Height: screenHeight})
	myWindow.SetContent(masterGrid)

	return myWindow
}

func createNewListOfRecipes(recipesStr []string) *widget.List {

	selectedRecipe := map[string]bool{}
	for _, r := range recipesStr {
		selectedRecipe[r] = false
	}
	// Recipe list with all recipes
	// TODO: Update this
	checkBoxs := binding.NewBoolList()
	// dataBindings := binding.NewBool()

	log.Println(len(recipesStr))
	log.Println(len(recipesStr))
	log.Println(len(recipesStr))
	log.Println(len(recipesStr))
	log.Println(len(recipesStr))

	// // c := widget.NewCheckGroup(recipesStr, nil)
	// var customCheckGroup *widget.Check
	for _ = range recipesStr {
		// currentBind := binding.NewBool()
		checkBoxs.Append(false)
		// customCheckGroup = widget.NewCheck("", func(b bool) {})
	}
	l := widget.NewList(
		func() int {
			return len(recipesStr)
		},
		func() fyne.CanvasObject {
			hbox := container.NewGridWithColumns(1)

			hbox1 := container.NewGridWithColumns(2)
			hbox1.Add(widget.NewCheck("", func(bool) {
			}))
			hbox1.Add(widget.NewLabel("table"))
			// Can't  get  the widget.ListItemID  there !!!
			hbox.Add(hbox1)

			return hbox
			// widget.NewCheck(recipesStr[0], func(b bool) {})
		},
		func(li widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			// for i := 0; i < 3; i++ {
			// 	if i == 3 {
			box1 := box.Objects[0].(*fyne.Container)
			check1 := box1.Objects[0].(*widget.Check)
			check1.Checked = selectedRecipe[recipesStr[li]]
			check1.Refresh()
			check1.OnChanged = func(b bool) {
				selectedRecipe[recipesStr[li]] = b
				log.Printf("recipe=%s val=%t\n %+v", recipesStr[li], b, selectedRecipe)
			}
			// bnd, _ := checkBoxs.
			// check1.Bind(bnd)
			lb1 := box1.Objects[1].(*widget.Label)
			// lb1 := box1.Objects[1].(*widget.Button)
			lb1.SetText(recipesStr[li])
			// 		lb2 := box1.Objects[1].(*widget.Button)
			// 		lb2.SetText("Delete")

			// 	} else {
			// 		lb1 := box.Objects[i].(*widget.Label)
			// 		lb1.SetText(recipesStr[i])
			// }

			// o.(widget.Check).Text = recipesStr[i]

			// log.Println("making check for", recipesStr[i])
			// _, _ = i.(binding.Bool).Get()
			// o.(*widget.Label).Text = "i.AddListener()"
			// o.(*widget.Check).Text = "apple"
			// log.Println(i.(binding.String).Get())
			// o.(*widget.Label).SetText(recipesStr[0])
			// o.(*widget.Check).Bind(txt)
			// o.(*widget.Check).Text = recipesStr[0]
		})
	// checkBoxs.GetItem()
	return l
}
