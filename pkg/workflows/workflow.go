package workflow

import (
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

	progressBarEmpty = 0.0
	progressBarFull  = 1.0
	recipeFinishLabel = "Finished. Select another recipe to add more."
)

// NewWorkflow will return a mac workflow if workflow file present and running on mac
// Else will return terminal workflow which prints to terminal
func NewWorkflow(f fileChecker, osString string) (WorkflowInterface, error) {
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

type screenInterface interface {
	updateProgessBar(float64)
	updateLabel(string)
}

//go:generate go run github.com/vektra/mockery/cmd/mockery -name workflowInterface -inpkg --filename workflow_mock.go
type WorkflowInterface interface {
	addIngredientsToReminders(ings []recipe.Ingredient, s screenInterface, w WorkflowInterface) error
	runReminder(s screenInterface, currentIng recipe.Ingredient) error
	submitShoppingList(s screenInterface, wf WorkflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error
}

func submitShoppingList(s screenInterface, wf WorkflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error {
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
	return wf.addIngredientsToReminders(ings, s, wf)
}

func ingQueue(ings []recipe.Ingredient, c chan<- recipe.Ingredient) {
	defer close(c)
	for _, ing := range ings {
		c <- ing
	}
}

func ingSend(s screenInterface, w WorkflowInterface, c <-chan recipe.Ingredient, ingAdded *int, totalIngs int) error {
	for ing := range c {
		log.Printf("ingredient=%s status=IN PROGRESS", ing.String())
		if err := w.runReminder(s, ing); err != nil {
			return err
		}
		*ingAdded++
		progress := float64(*ingAdded) / float64(totalIngs)
		s.updateProgessBar(progress)
		log.Printf("ingredient=%s status=DONE progress=%.2f", ing.String(), progress)
	}
	return nil
}

func addIngredientsToReminders(ings []recipe.Ingredient, s screenInterface, w WorkflowInterface) error {
	progress := float64(progressBarEmpty)
	s.updateProgessBar(progress)

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
	s.updateProgessBar(progress)
	s.updateLabel(recipeFinishLabel)
	return nil
}
