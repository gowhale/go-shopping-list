package workflow

import (
	"fmt"
	"go-shopping-list/pkg/common"
	"go-shopping-list/pkg/recipe"
	"os/exec"
)

var execCommand = exec.Command

type macWorkflow struct{}

func (*macWorkflow) SubmitShoppingList(s common.ScreenInterface, wf common.WorkflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error {
	return SubmitShoppingList(s, wf, fr, recipes, recipeMap)
}

func (*macWorkflow) AddIngredientsToReminders(ings []recipe.Ingredient, s common.ScreenInterface, w common.WorkflowInterface) error {
	return AddIngredientsToReminders(ings, s, w)
}

func (*macWorkflow) RunReminder(s common.ScreenInterface, currentIng recipe.Ingredient) error {
	cmd := execCommand("automator", "-i", fmt.Sprintf(`"%s"`, currentIng.String()), "shopping.workflow")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding the following ingredient=%s err=%w", currentIng.String(), err)
	}
	s.UpdateLabel(fmt.Sprintf("Added Ingredient: %s", currentIng.String()))
	return nil
}
