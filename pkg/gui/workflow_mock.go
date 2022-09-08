// Code generated by mockery v1.0.0. DO NOT EDIT.

package gui

import (
	recipe "go-shopping-list/pkg/recipe"

	mock "github.com/stretchr/testify/mock"
)

// mockWorkflowInterface is an autogenerated mock type for the workflowInterface type
type mockWorkflowInterface struct {
	mock.Mock
}

// addIngredientsToReminders provides a mock function with given fields: ings, s, f, w
func (_m *mockWorkflowInterface) addIngredientsToReminders(ings []recipe.Ingredient, s screenInterface, f recipe.FileReader, w workflowInterface) error {
	ret := _m.Called(ings, s, f, w)

	var r0 error
	if rf, ok := ret.Get(0).(func([]recipe.Ingredient, screenInterface, recipe.FileReader, workflowInterface) error); ok {
		r0 = rf(ings, s, f, w)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// runReminder provides a mock function with given fields: s, currentIng
func (_m *mockWorkflowInterface) runReminder(s screenInterface, currentIng recipe.Ingredient) error {
	ret := _m.Called(s, currentIng)

	var r0 error
	if rf, ok := ret.Get(0).(func(screenInterface, recipe.Ingredient) error); ok {
		r0 = rf(s, currentIng)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// submitShoppingList provides a mock function with given fields: s, wf, fr, recipes, recipeMap
func (_m *mockWorkflowInterface) submitShoppingList(s screenInterface, wf workflowInterface, fr recipe.FileReader, recipes []string, recipeMap map[string]recipe.Recipe) error {
	ret := _m.Called(s, wf, fr, recipes, recipeMap)

	var r0 error
	if rf, ok := ret.Get(0).(func(screenInterface, workflowInterface, recipe.FileReader, []string, map[string]recipe.Recipe) error); ok {
		r0 = rf(s, wf, fr, recipes, recipeMap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
