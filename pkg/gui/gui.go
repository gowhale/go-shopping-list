package gui

import (
	"fmt"
	"go-shopping-list/pkg/recipe"
	"log"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/sync/errgroup"
)

const (
	screenWidth  = 600
	screenHeight = 1200
)

type screen struct {
	p *widget.ProgressBar
	l *widget.Label
}

//go:generate go run github.com/vektra/mockery/cmd/mockery -name ScreenInterface -inpkg --filename screen_mock.go
type ScreenInterface interface {
	UpdateProgessBar(float64)
	UpdateLabel(string)
}

type macWorkflow struct{}

//go:generate go run github.com/vektra/mockery/cmd/mockery -name WorkflowInterface -inpkg --filename workflow_mock.go
type WorkflowInterface interface {
	runReminder(s ScreenInterface, currentIng recipe.Ingredients) error
}

func (s *screen) UpdateProgessBar(percent float64) {
	s.p.SetValue(percent)
	s.p.Refresh()
}

func (s *screen) UpdateLabel(msg string) {
	s.l.SetText(msg)
	s.l.Refresh()
}

func NewApp(recipes []recipe.Recipe) fyne.Window {

	myApp := app.New()
	myWindow := myApp.NewWindow("List Widget")

	label := widget.NewLabel("Click a recipe to add the ingredients...")

	// Progress bar for adding ings
	p := widget.NewProgressBar()

	s := &screen{
		l: label,
		p: p,
	}

	w := macWorkflow{}

	// Recipe list with all recipes
	recipeList := createNewListOfRecipes(s, &recipe.FileInteractionImpl{}, &w, recipes)

	// Create content grid
	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(600, 1150)), recipeList)
	gridTop := container.New(layout.NewGridWrapLayout(fyne.NewSize(600, 50)), label, p)
	masterGrid := container.New(layout.NewVBoxLayout(), gridTop, grid)

	// Set Window and execute
	myWindow.SetFixedSize(true)
	myWindow.Resize(fyne.Size{Width: screenWidth, Height: screenHeight})
	myWindow.SetContent(masterGrid)

	return myWindow
}

func createNewListOfRecipes(s ScreenInterface, f recipe.FileReader, w WorkflowInterface, recipes []recipe.Recipe) *widget.List {
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
			o.(*widget.Button).OnTapped = func() { itemClicked(s, recipes[i], f, w) }
		})
}

func itemClicked(s ScreenInterface, r recipe.Recipe, f recipe.FileReader, w WorkflowInterface) {
	err := AddIngredientsToReminders(r, s, f, w)
	if err != nil {
		log.Printf("error whilst adding ingredients to reminds err=%e", err)
	}
}

func AddIngredientsToReminders(r recipe.Recipe, s ScreenInterface, f recipe.FileReader, w WorkflowInterface) error {
	s.UpdateLabel(fmt.Sprintf("Starting to add ingredients for Recipe: %s", r.Name))

	progress := float64(0.0)
	s.UpdateProgessBar(progress)
	ingAdded := []recipe.Ingredients{}

	g := new(errgroup.Group)
	for _, ing := range r.Ings {
		ing := ing
		g.Go(func() error {
			if err := w.runReminder(s, ing); err != nil {
				return err
			}
			defer func() {
				ingAdded = append(ingAdded, ing)
				progress = float64(len(ingAdded)) / float64(len(r.Ings))
				s.UpdateProgessBar(progress)
				log.Printf("progress=%.2f adding ing='%s'", progress, ing.String())
			}()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	if err := f.IncrementPopularity(r.Name); err != nil {
		return err
	}

	progress = 1
	log.Printf("progress=%.2f", progress)
	s.UpdateProgessBar(progress)
	s.UpdateLabel("Finished. Select another recipe to add more.")
	return nil
}

var execCommand = exec.Command

func (*macWorkflow) runReminder(s ScreenInterface, currentIng recipe.Ingredients) error {
	cmd := execCommand("automator", "-i", fmt.Sprintf(`"%s"`, currentIng.String()), "shopping.workflow")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding the following ingredient=%s err=%w", currentIng.String(), err)
	}
	s.UpdateLabel(fmt.Sprintf("Added Ingredient: %s", currentIng.String()))
	return nil
}
