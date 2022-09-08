package gui

import (
	"fmt"
	"go-shopping-list/pkg/recipe"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	workflowName = "shopping.workflow"
	macOSName    = "darwin"

	minMilliseconds = 100
	maxMilliseconds = 500
)

// NewWorkflow will return a mac workflow if workflow file present and running on mac
// Else will return terminal workflow which prints to terminal
func NewWorkflow(f fileChecker, osString string) (workflowInterface, error) {
	workflowPresent, err := f.checkWorkflowExists(f)
	if err != nil {
		return nil, err
	}
	if workflowPresent {
		if osString == macOSName {
			log.Println("Using mac workflow to create reminders")
			return &macWorkflow{}, nil
		}
		log.Println("Not running on mac!")
	} else {
		log.Println("No workflow found!")
	}
	log.Println("Printing to terminal to simulate adding of ingredients!")
	return &TerminalFakeWorkflow{}, nil
}

// WorkflowChecker is struct used to check if the shopping.workflow file exists
type WorkflowChecker struct{}

//go:generate go run github.com/vektra/mockery/cmd/mockery -name fileChecker -inpkg --filename file_checker_mock.go
type fileChecker interface {
	checkWorkflowExists(f fileChecker) (bool, error)
	stat(name string) (fs.FileInfo, error)
	isNotExist(err error) bool
}

// checkWorkflowExists tests to see if the mac workflow exists
func (*WorkflowChecker) checkWorkflowExists(f fileChecker) (bool, error) {
	return checkWorkflowExistsImpl(f)
}

func (*WorkflowChecker) stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (*WorkflowChecker) isNotExist(err error) bool {
	return os.IsNotExist(err)
}

func checkWorkflowExistsImpl(f fileChecker) (bool, error) {
	_, err := f.stat(workflowName)
	if err == nil {
		return true, nil
	}
	if f.isNotExist(err) {
		return false, nil
	}
	return false, err
}

//go:generate go run github.com/vektra/mockery/cmd/mockery -name workflowInterface -inpkg --filename workflow_mock.go
type workflowInterface interface {
	addIngredientsToReminders(ings []recipe.Ingredient, s screenInterface, w workflowInterface) error
	runReminder(s screenInterface, currentIng recipe.Ingredient) error
	submitShoppingList(s screenInterface, wf workflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error
}

func submitShoppingList(s screenInterface, wf workflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error {
	log.Println("Currently selected Recipes:")
	recipesSelected := []recipe.Recipe{}
	for _, v := range recipes {
		if r, ok := recipeMap[v]; ok {
			recipesSelected = append(recipesSelected, r)
			if err := fr.IncrementPopularity(fr, r.Name); err != nil {
				return err
			}
		}

	}
	ings, err := recipe.CombineRecipesToIngredients(recipesSelected)
	if err != nil {
		return err
	}
	return wf.addIngredientsToReminders(ings, s, wf)
}

func addIngredientsToReminders(ings []recipe.Ingredient, s screenInterface, w workflowInterface) error {
	progress := float64(progressBarEmpty)
	s.updateProgessBar(progress)
	ingAdded := []recipe.Ingredient{}

	g := new(errgroup.Group)
	for _, ing := range ings {
		ing := ing
		g.Go(func() error {
			if err := w.runReminder(s, ing); err != nil {
				return err
			}
			defer func() {
				ingAdded = append(ingAdded, ing)
				progress = float64(len(ingAdded)) / float64(len(ings))
				s.updateProgessBar(progress)
				log.Printf("progress=%.2f adding ing='%s'", progress, ing.String())
			}()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	progress = progressBarFull
	log.Printf("progress=%.2f", progress)
	s.updateProgessBar(progress)
	s.updateLabel(recipeFinishLabel)
	return nil
}

var execCommand = exec.Command

type macWorkflow struct{}

func (*macWorkflow) submitShoppingList(s screenInterface, wf workflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error {
	return submitShoppingList(s, wf, fr, recipes, recipeMap)
}

func (*macWorkflow) addIngredientsToReminders(ings []recipe.Ingredient, s screenInterface, w workflowInterface) error {
	return addIngredientsToReminders(ings, s, w)
}

func (*macWorkflow) runReminder(s screenInterface, currentIng recipe.Ingredient) error {
	cmd := execCommand("automator", "-i", fmt.Sprintf(`"%s"`, currentIng.String()), "shopping.workflow")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding the following ingredient=%s err=%w", currentIng.String(), err)
	}
	s.updateLabel(fmt.Sprintf("Added Ingredient: %s", currentIng.String()))
	return nil
}

// TerminalFakeWorkflow can be used to just print to termnial
type TerminalFakeWorkflow struct{}

func (*TerminalFakeWorkflow) submitShoppingList(s screenInterface, wf workflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error {
	return submitShoppingList(s, wf, fr, recipes, recipeMap)
}

func (*TerminalFakeWorkflow) addIngredientsToReminders(ings []recipe.Ingredient, s screenInterface, w workflowInterface) error {
	return addIngredientsToReminders(ings, s, w)
}

func (*TerminalFakeWorkflow) runReminder(_ screenInterface, currentIng recipe.Ingredient) error {
	log.Printf("PRETENDING TO ADD INGREDIENT=%s", currentIng.String())
	millisecondsToWait := rand.Intn(maxMilliseconds-minMilliseconds) + minMilliseconds
	time.Sleep(time.Millisecond * time.Duration(millisecondsToWait))
	// The below line creates a bug. I think because race conditions. Maybe I should implement mutex?
	// s.updateLabel(fmt.Sprintf("Added Ingredient: %s", currentIng.String()))
	return nil
}
