package gui

import (
	"fmt"
	"go-shopping-list/pkg/recipe"
	"testing"

	"github.com/stretchr/testify/suite"
)

type guiTest struct {
	suite.Suite
	mockScreen     *MockScreenInterface
	mockFileReader *recipe.MockFileReader
	mockWorkflow   *MockWorkflowInterface
}

func (g *guiTest) SetupTest() {
	g.mockScreen = new(MockScreenInterface)
	g.mockFileReader = new(recipe.MockFileReader)
	g.mockWorkflow = new(MockWorkflowInterface)
}

func TestGuiTest(t *testing.T) {
	suite.Run(t, new(guiTest))
}

func (g *guiTest) Test_mockFileInfo() {
	testRecipe := []recipe.Recipe{}
	_ = NewApp(testRecipe)
}

func (g *guiTest) Test_AddIngredientsToReminders_Pass() {
	ing := recipe.Ingredients{
		Unit_size:       "PEAR",
		Unit_type:       "BANANA",
		Ingredient_name: "RASPBERRY",
	}
	testRecipe := recipe.Recipe{
		Name: "APPLE",
		Ings: []recipe.Ingredients{ing},
	}
	g.mockScreen.On("UpdateLabel", "Starting to add ingredients for Recipe: APPLE")
	g.mockScreen.On("UpdateProgessBar", float64(0.0))
	g.mockScreen.On("UpdateProgessBar", float64(1.0))
	g.mockScreen.On("UpdateLabel", "Finished. Select another recipe to add more.")
	g.mockWorkflow.On("runReminder", g.mockScreen, ing).Return(nil)
	g.mockFileReader.On("IncrementPopularity", "APPLE").Return(nil)
	err := AddIngredientsToReminders(testRecipe, g.mockScreen, g.mockFileReader, g.mockWorkflow)
	g.Nil(err)
}

func (g *guiTest) Test_AddIngredientsToReminders_IncrementPopularity_Error() {
	ing := recipe.Ingredients{
		Unit_size:       "PEAR",
		Unit_type:       "BANANA",
		Ingredient_name: "RASPBERRY",
	}
	testRecipe := recipe.Recipe{
		Name: "APPLE",
		Ings: []recipe.Ingredients{ing},
	}
	g.mockScreen.On("UpdateLabel", "Starting to add ingredients for Recipe: APPLE")
	g.mockScreen.On("UpdateProgessBar", float64(0.0))
	g.mockScreen.On("UpdateProgessBar", float64(1.0))
	g.mockScreen.On("UpdateLabel", "Finished. Select another recipe to add more.")
	g.mockWorkflow.On("runReminder", g.mockScreen, ing).Return(nil)
	g.mockFileReader.On("IncrementPopularity", "APPLE").Return(fmt.Errorf("pop error"))
	err := AddIngredientsToReminders(testRecipe, g.mockScreen, g.mockFileReader, g.mockWorkflow)
	g.EqualError(err, "pop error")
}

func (g *guiTest) Test_AddIngredientsToReminders_runReminder_Error() {
	ing := recipe.Ingredients{
		Unit_size:       "PEAR",
		Unit_type:       "BANANA",
		Ingredient_name: "RASPBERRY",
	}

	testRecipe := recipe.Recipe{
		Name: "APPLE",
		Ings: []recipe.Ingredients{ing},
	}

	g.mockScreen.On("UpdateLabel", "Starting to add ingredients for Recipe: APPLE")
	g.mockScreen.On("UpdateProgessBar", float64(0.0))
	g.mockScreen.On("UpdateProgessBar", float64(1.0))
	g.mockScreen.On("UpdateLabel", "Finished. Select another recipe to add more.")
	g.mockWorkflow.On("runReminder", g.mockScreen, ing).Return(fmt.Errorf("reminder error"))
	g.mockFileReader.On("IncrementPopularity", "APPLE").Return(nil)
	err := AddIngredientsToReminders(testRecipe, g.mockScreen, g.mockFileReader, g.mockWorkflow)
	g.EqualError(err, "reminder error")
}
