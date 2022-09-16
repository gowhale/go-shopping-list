// Package workflow is used by all workflow files
package workflow

import (
	"go-shopping-list/pkg/common"
	"go-shopping-list/pkg/recipe"
	"io/fs"
	"log"
	"os"

	"golang.org/x/sync/errgroup"
)

const (
	workflowName = "shopping.workflow"
	macOSName    = "darwin"

	minMilliseconds = 100
	maxMilliseconds = 500

	numOfGoRoutines       = 10
	ingredientsCountStart = 0

	progressBarEmpty  = 0.0
	progressBarFull   = 1.0
	recipeFinishLabel = "Finished. Select another recipe to add more."
)

// NewWorkflow will return a mac workflow if workflow file present and running on mac
// Else will return terminal workflow which prints to terminal
func NewWorkflow(f fileChecker, osString string) (common.WorkflowInterface, error) {
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

// CheckWorkflow is struct used to check if the shopping.workflow file exists
type CheckWorkflow struct{}

//go:generate go run github.com/vektra/mockery/cmd/mockery -name fileChecker -inpkg --filename file_checker_mock.go
type fileChecker interface {
	checkWorkflowExists(f fileChecker) (bool, error)
	stat(name string) (fs.FileInfo, error)
	isNotExist(err error) bool
}

// checkWorkflowExists tests to see if the mac workflow exists
func (*CheckWorkflow) checkWorkflowExists(f fileChecker) (bool, error) {
	return checkWorkflowExistsImpl(f)
}

func (*CheckWorkflow) stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (*CheckWorkflow) isNotExist(err error) bool {
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

// SubmitShoppingList combines recipes and then creates shopping list
func SubmitShoppingList(s common.ScreenInterface, wf common.WorkflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error {
	log.Println("Currently selected Recipes:")
	recipesSelected := []recipe.Recipe{}
	for _, v := range recipes {
		log.Println(v)
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
	return wf.AddIngredientsToReminders(ings, s, wf)
}

func ingQueue(ings []recipe.Ingredient, c chan<- recipe.Ingredient) {
	defer close(c)
	for _, ing := range ings {
		c <- ing
	}
}

func ingSend(s common.ScreenInterface, w common.WorkflowInterface, c <-chan recipe.Ingredient, ingAdded *int, totalIngs int) error {
	for ing := range c {
		log.Printf("ingredient=%s status=IN PROGRESS", ing.String())
		if err := w.RunReminder(s, ing); err != nil {
			return err
		}
		*ingAdded++
		progress := float64(*ingAdded) / float64(totalIngs)
		s.UpdateProgessBar(progress)
		log.Printf("ingredient=%s status=DONE progress=%.2f", ing.String(), progress)
	}
	return nil
}

// AddIngredientsToReminders adds ingredients to reminders app
func AddIngredientsToReminders(ings []recipe.Ingredient, s common.ScreenInterface, w common.WorkflowInterface) error {
	progress := float64(progressBarEmpty)
	s.UpdateProgessBar(progress)

	ingAdded := ingredientsCountStart
	ingWaitingList := make(chan recipe.Ingredient, numOfGoRoutines)
	g := new(errgroup.Group)
	for i := ingredientsCountStart; i < numOfGoRoutines; i++ {
		g.Go(func() error {
			return ingSend(s, w, ingWaitingList, &ingAdded, len(ings))
		})
	}

	ingQueue(ings, ingWaitingList)
	if err := g.Wait(); err != nil {
		return err
	}

	progress = progressBarFull
	log.Printf("progress=%.2f", progress)
	s.UpdateProgessBar(progress)
	s.UpdateLabel(recipeFinishLabel)
	return nil
}
