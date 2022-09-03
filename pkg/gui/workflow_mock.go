// Code generated by mockery v1.0.0. DO NOT EDIT.

package gui

import (
	recipe "go-shopping-list/pkg/recipe"

	mock "github.com/stretchr/testify/mock"
)

// MockWorkflowInterface is an autogenerated mock type for the workflowInterface type
type MockWorkflowInterface struct {
	mock.Mock
}

// runReminder provides a mock function with given fields: s, currentIng
func (_m *MockWorkflowInterface) runReminder(s screenInterface, currentIng recipe.Ingredients) error {
	ret := _m.Called(s, currentIng)

	var r0 error
	if rf, ok := ret.Get(0).(func(screenInterface, recipe.Ingredients) error); ok {
		r0 = rf(s, currentIng)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
