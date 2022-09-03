// Package gui is responsible for code relating to user interaction
package gui

import (
	"fmt"
	"go-shopping-list/pkg/recipe"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	exitCodePass         = 0
	exitCodeFail         = 1
	firstArgumentInSlice = 0

	IncrementPopularityString = "IncrementPopularity"
	runReminderString         = "runReminder"
	updateLabelString         = "updateLabel"
	updateProgessBarString    = "updateProgessBar"
)

type guiTest struct {
	suite.Suite
	mockScreen     *mockScreenInterface
	mockFileReader *recipe.MockFileReader
	mockWorkflow   *MockWorkflowInterface
}

func (g *guiTest) SetupTest() {
	g.mockScreen = new(mockScreenInterface)
	g.mockFileReader = new(recipe.MockFileReader)
	g.mockWorkflow = new(MockWorkflowInterface)
}

func TestGuiTest(t *testing.T) {
	suite.Run(t, new(guiTest))
}

func (*guiTest) Test_mockFileInfo() {
	testRecipe := []recipe.Recipe{}
	_ = NewApp(testRecipe)
}

func (g *guiTest) Test_AddIngredientsToReminders_Pass() {
	ing := recipe.Ingredients{
		UnitSize:       "PEAR",
		UnitType:       "BANANA",
		IngredientName: "LYCHE",
	}
	testRecipe := recipe.Recipe{
		Name: "DURIAN",
		Ings: []recipe.Ingredients{ing},
	}
	g.mockScreen.On(updateLabelString, "Starting to add ingredients for Recipe: DURIAN")
	g.mockScreen.On(updateProgessBarString, progressBarEmpty)
	g.mockScreen.On(updateProgessBarString, progressBarFull)
	g.mockScreen.On(updateLabelString, recipeFinishLabel)
	g.mockWorkflow.On(runReminderString, g.mockScreen, ing).Return(nil)
	g.mockFileReader.On(IncrementPopularityString, "DURIAN").Return(nil)
	err := addIngredientsToReminders(testRecipe, g.mockScreen, g.mockFileReader, g.mockWorkflow)
	g.Nil(err)
}

func (g *guiTest) Test_AddIngredientsToReminders_IncrementPopularity_Error() {
	ing := recipe.Ingredients{
		UnitSize:       "PEACH",
		UnitType:       "BLUEBERRY",
		IngredientName: "RASPBERRY",
	}
	testRecipe := recipe.Recipe{
		Name: "MANGO",
		Ings: []recipe.Ingredients{ing},
	}
	g.mockScreen.On(updateLabelString, "Starting to add ingredients for Recipe: MANGO")
	g.mockScreen.On(updateProgessBarString, progressBarEmpty)
	g.mockScreen.On(updateProgessBarString, progressBarFull)
	g.mockScreen.On(updateLabelString, recipeFinishLabel)
	g.mockWorkflow.On(runReminderString, g.mockScreen, ing).Return(nil)
	g.mockFileReader.On(IncrementPopularityString, "MANGO").Return(fmt.Errorf("pop error"))
	err := addIngredientsToReminders(testRecipe, g.mockScreen, g.mockFileReader, g.mockWorkflow)
	g.EqualError(err, "pop error")
}

func (g *guiTest) Test_AddIngredientsToReminders_runReminder_Error() {
	ing := recipe.Ingredients{
		UnitSize:       "ORANGE",
		UnitType:       "BANANA",
		IngredientName: "RASPBERRY",
	}

	testRecipe := recipe.Recipe{
		Name: "APPLE",
		Ings: []recipe.Ingredients{ing},
	}

	g.mockScreen.On(updateLabelString, "Starting to add ingredients for Recipe: APPLE")
	g.mockScreen.On(updateProgessBarString, progressBarEmpty)
	g.mockScreen.On(updateProgessBarString, progressBarFull)
	g.mockScreen.On(updateLabelString, recipeFinishLabel)
	g.mockWorkflow.On(runReminderString, g.mockScreen, ing).Return(fmt.Errorf("reminder error"))
	g.mockFileReader.On(IncrementPopularityString, "APPLE").Return(nil)
	err := addIngredientsToReminders(testRecipe, g.mockScreen, g.mockFileReader, g.mockWorkflow)
	g.EqualError(err, "reminder error")
}

func TestHelperProcess(*testing.T) {
	helper := os.Getenv("GO_WANT_HELPER_PROCESS")
	//pass
	if helper == "1" {
		os.Exit(exitCodePass)
		return
	}
	//fail
	if helper == "2" {
		os.Exit(exitCodeFail)
		return
	}
}

func fakeExecCommandPass(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[firstArgumentInSlice], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func fakeExecCommandFail(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[firstArgumentInSlice], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=2"}
	return cmd
}

func (g *guiTest) Test_runReminder_Pass() {
	execCommand = fakeExecCommandPass
	defer func() {
		execCommand = exec.Command
	}()
	m := macWorkflow{}
	ing := recipe.Ingredients{
		UnitSize:       "WATERMELON",
		UnitType:       "CHERRY",
		IngredientName: "PITAYA",
	}
	g.mockScreen.On(updateLabelString, "Added Ingredient: WATERMELON CHERRY PITAYA")
	err := m.runReminder(g.mockScreen, ing)
	g.Nil(err)
}

func (g *guiTest) Test_runReminder_Error() {
	execCommand = fakeExecCommandFail
	defer func() {
		execCommand = exec.Command
	}()
	m := macWorkflow{}
	ing := recipe.Ingredients{
		UnitSize:       "PEAR",
		UnitType:       "FIG",
		IngredientName: "AVOCADO",
	}
	g.mockScreen.On(updateLabelString, "Added Ingredient: PEAR FIG AVOCADO")
	err := m.runReminder(g.mockScreen, ing)
	g.EqualError(err, "error adding the following ingredient=PEAR FIG AVOCADO err=exit status 1")
}
