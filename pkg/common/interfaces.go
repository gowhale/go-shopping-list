package common

import "go-shopping-list/pkg/recipe"

//go:generate go run github.com/vektra/mockery/cmd/mockery -name workflowInterface -inpkg --filename workflow_mock.go
type WorkflowInterface interface {
	addIngredientsToReminders(ings []recipe.Ingredient, s WorkflowInterface, w WorkflowInterface) error
	runReminder(s WorkflowInterface, currentIng recipe.Ingredient) error
	submitShoppingList(s WorkflowInterface, wf WorkflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error
}
