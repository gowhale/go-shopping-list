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
)

func NewWorkflow(f fileChecker, os string) (workflowInterface, error) {
	workflowPresent, err := f.checkWorkflowExists()
	if err != nil {
		return nil, err
	}
	if workflowPresent {
		if os == macOSName {
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

type WorkflowChecker struct{}

//go:generate go run github.com/vektra/mockery/cmd/mockery -name fileChecker -inpkg --filename file_checker_mock.go
type fileChecker interface {
	checkWorkflowExists() (bool, error)
	stat(name string) (fs.FileInfo, error)
	isNotExist(err error) bool
}

// checkWorkflowExists tests to see if the mac workflow exists
func (f *WorkflowChecker) checkWorkflowExists() (bool, error) {
	return checkWorkflowExistsImpl(f)
}

func (f *WorkflowChecker) stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (f *WorkflowChecker) isNotExist(err error) bool {
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
	runReminder(s screenInterface, currentIng recipe.Ingredient) error
}

func addIngredientsToReminders(r recipe.Recipe, s screenInterface, f recipe.FileReader, w workflowInterface) error {
	s.updateLabel(fmt.Sprintf("Starting to add ingredients for Recipe: %s", r.Name))

	progress := float64(progressBarEmpty)
	s.updateProgessBar(progress)
	ingAdded := []recipe.Ingredient{}

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
				s.updateProgessBar(progress)
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

	progress = progressBarFull
	log.Printf("progress=%.2f", progress)
	s.updateProgessBar(progress)
	s.updateLabel(recipeFinishLabel)
	return nil
}

var execCommand = exec.Command

type macWorkflow struct{}

func (*macWorkflow) runReminder(s screenInterface, currentIng recipe.Ingredient) error {
	cmd := execCommand("automator", "-i", fmt.Sprintf(`"%s"`, currentIng.String()), "shopping.workflow")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding the following ingredient=%s err=%w", currentIng.String(), err)
	}
	s.updateLabel(fmt.Sprintf("Added Ingredient: %s", currentIng.String()))
	return nil
}

type TerminalFakeWorkflow struct{}

func (*TerminalFakeWorkflow) runReminder(s screenInterface, currentIng recipe.Ingredient) error {
	log.Printf("PRETENDING TO ADD INGREDIENT=%s", currentIng.String())
	minMilliseconds := 100
	maxMilliseconds := 500
	millisecondsToWait := rand.Intn(maxMilliseconds-minMilliseconds) + minMilliseconds
	time.Sleep(time.Millisecond * time.Duration(millisecondsToWait))
	// The below line creates a bug. I think because race conditions. Maybe I should implement mutex?
	// s.updateLabel(fmt.Sprintf("Added Ingredient: %s", currentIng.String()))
	return nil
}
