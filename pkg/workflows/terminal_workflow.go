package workflow

import (
	"go-shopping-list/pkg/common"
	"go-shopping-list/pkg/recipe"
	"log"
	"math/rand"
	"time"
)

// TerminalFakeWorkflow can be used to just print to termnial
type TerminalFakeWorkflow struct{}

func (*TerminalFakeWorkflow) SubmitShoppingList(s common.ScreenInterface, wf common.WorkflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error {
	return SubmitShoppingList(s, wf, fr, recipes, recipeMap)
}

func (*TerminalFakeWorkflow) AddIngredientsToReminders(ings []recipe.Ingredient, s common.ScreenInterface, w common.WorkflowInterface) error {
	return AddIngredientsToReminders(ings, s, w)
}

func (*TerminalFakeWorkflow) RunReminder(_ common.ScreenInterface, currentIng recipe.Ingredient) error {
	log.Printf("PRETENDING TO ADD INGREDIENT=%s", currentIng.String())
	millisecondsToWait := rand.Intn(maxMilliseconds-minMilliseconds) + minMilliseconds
	time.Sleep(time.Millisecond * time.Duration(millisecondsToWait))
	// The below line creates a bug. I think because race conditions. Maybe I should implement mutex?
	// s.UpdateLabel(fmt.Sprintf("Added Ingredient: %s", currentIng.String()))
	return nil
}
