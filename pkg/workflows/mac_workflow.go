package workflow

import (
	"fmt"
	"go-shopping-list/pkg/recipe"
	"os/exec"
)

var execCommand = exec.Command

type macWorkflow struct{}

func (*macWorkflow) submitShoppingList(s screenInterface, wf WorkflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error {
	return submitShoppingList(s, wf, fr, recipes, recipeMap)
}

func (*macWorkflow) addIngredientsToReminders(ings []recipe.Ingredient, s screenInterface, w WorkflowInterface) error {
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
