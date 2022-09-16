package workflow

import (
	"go-shopping-list/pkg/recipe"
	"log"
	"math/rand"
	"time"
)

// TerminalFakeWorkflow can be used to just print to termnial
type TerminalFakeWorkflow struct{}

func (*TerminalFakeWorkflow) submitShoppingList(s screenInterface, wf WorkflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error {
	return submitShoppingList(s, wf, fr, recipes, recipeMap)
}

func (*TerminalFakeWorkflow) addIngredientsToReminders(ings []recipe.Ingredient, s screenInterface, w WorkflowInterface) error {
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
