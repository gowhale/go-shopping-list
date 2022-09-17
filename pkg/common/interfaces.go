// Package common contains code used between pkgs
package common

import "go-shopping-list/pkg/recipe"

// WorkflowInterface provides common interface used by all workflows
//
//go:generate go run github.com/vektra/mockery/cmd/mockery -name WorkflowInterface -inpkg --filename workflow_mock.go
type WorkflowInterface interface {
	AddIngredientsToReminders(ings []recipe.Ingredient, s ScreenInterface, w WorkflowInterface) error
	RunReminder(s ScreenInterface, currentIng recipe.Ingredient) error
	SubmitShoppingList(s ScreenInterface, wf WorkflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error
}

// ScreenInterface provides common interface
// used gui and workflow pkgs
//
//go:generate go run github.com/vektra/mockery/cmd/mockery -name ScreenInterface -inpkg --filename screen_mock.go
type ScreenInterface interface {
	UpdateProgessBar(float64)
	UpdateLabel(string)
}
