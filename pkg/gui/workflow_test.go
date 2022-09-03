// Package gui is responsible for code relating to user interaction
package gui

import (
	"fmt"
	"go-shopping-list/pkg/recipe"
	fs "io/fs"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	exitCodePass         = 0
	exitCodeFail         = 1
	firstArgumentInSlice = 0

	incrementPopularityString = "IncrementPopularity"
	runReminderString         = "runReminder"
	updateLabelString         = "updateLabel"
	updateProgessBarString    = "updateProgessBar"
)

type workflowTest struct {
	suite.Suite
	mockScreen      *mockScreenInterface
	mockFileReader  *recipe.MockFileReader
	mockWorkflow    *mockWorkflowInterface
	mockFileChecker *mockFileChecker
}

func (g *workflowTest) SetupTest() {
	g.mockScreen = new(mockScreenInterface)
	g.mockFileReader = new(recipe.MockFileReader)
	g.mockWorkflow = new(mockWorkflowInterface)
	g.mockFileChecker = new(mockFileChecker)
}

func Test_workflowTest(t *testing.T) {
	suite.Run(t, new(workflowTest))
}

func (*workflowTest) Test_mockFileInfo() {
	testRecipe := []recipe.Recipe{}
	_ = NewApp(testRecipe, &TerminalFakeWorkflow{})
}

func (g *workflowTest) Test_AddIngredientsToReminders_Pass() {
	ing := recipe.Ingredient{
		UnitSize:       "PEAR",
		UnitType:       "BANANA",
		IngredientName: "LYCHE",
	}
	testRecipe := recipe.Recipe{
		Name: "DURIAN",
		Ings: []recipe.Ingredient{ing},
	}
	g.mockScreen.On(updateLabelString, "Starting to add ingredients for Recipe: DURIAN")
	g.mockScreen.On(updateProgessBarString, progressBarEmpty)
	g.mockScreen.On(updateProgessBarString, progressBarFull)
	g.mockScreen.On(updateLabelString, recipeFinishLabel)
	g.mockWorkflow.On(runReminderString, g.mockScreen, ing).Return(nil)
	g.mockFileReader.On(incrementPopularityString, "DURIAN").Return(nil)
	err := addIngredientsToReminders(testRecipe, g.mockScreen, g.mockFileReader, g.mockWorkflow)
	g.Nil(err)
}

func (g *workflowTest) Test_AddIngredientsToReminders_IncrementPopularity_Error() {
	ing := recipe.Ingredient{
		UnitSize:       "PEACH",
		UnitType:       "BLUEBERRY",
		IngredientName: "RASPBERRY",
	}
	testRecipe := recipe.Recipe{
		Name: "MANGO",
		Ings: []recipe.Ingredient{ing},
	}
	g.mockScreen.On(updateLabelString, "Starting to add ingredients for Recipe: MANGO")
	g.mockScreen.On(updateProgessBarString, progressBarEmpty)
	g.mockScreen.On(updateProgessBarString, progressBarFull)
	g.mockScreen.On(updateLabelString, recipeFinishLabel)
	g.mockWorkflow.On(runReminderString, g.mockScreen, ing).Return(nil)
	g.mockFileReader.On(incrementPopularityString, "MANGO").Return(fmt.Errorf("pop error"))
	err := addIngredientsToReminders(testRecipe, g.mockScreen, g.mockFileReader, g.mockWorkflow)
	g.EqualError(err, "pop error")
}

func (g *workflowTest) Test_AddIngredientsToReminders_runReminder_Error() {
	ing := recipe.Ingredient{
		UnitSize:       "ORANGE",
		UnitType:       "BANANA",
		IngredientName: "RASPBERRY",
	}

	testRecipe := recipe.Recipe{
		Name: "APPLE",
		Ings: []recipe.Ingredient{ing},
	}

	g.mockScreen.On(updateLabelString, "Starting to add ingredients for Recipe: APPLE")
	g.mockScreen.On(updateProgessBarString, progressBarEmpty)
	g.mockScreen.On(updateProgessBarString, progressBarFull)
	g.mockScreen.On(updateLabelString, recipeFinishLabel)
	g.mockWorkflow.On(runReminderString, g.mockScreen, ing).Return(fmt.Errorf("reminder error"))
	g.mockFileReader.On(incrementPopularityString, "APPLE").Return(nil)
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

func (g *workflowTest) Test_macWorkflow_runReminder_Pass() {
	execCommand = fakeExecCommandPass
	defer func() {
		execCommand = exec.Command
	}()
	m := macWorkflow{}
	ing := recipe.Ingredient{
		UnitSize:       "WATERMELON",
		UnitType:       "CHERRY",
		IngredientName: "PITAYA",
	}
	g.mockScreen.On(updateLabelString, "Added Ingredient: WATERMELON CHERRY PITAYA")
	err := m.runReminder(g.mockScreen, ing)
	g.Nil(err)
}

func (g *workflowTest) Test_macWorkflow_runReminder_Error() {
	execCommand = fakeExecCommandFail
	defer func() {
		execCommand = exec.Command
	}()
	m := macWorkflow{}
	ing := recipe.Ingredient{
		UnitSize:       "PEAR",
		UnitType:       "FIG",
		IngredientName: "AVOCADO",
	}
	g.mockScreen.On(updateLabelString, "Added Ingredient: PEAR FIG AVOCADO")
	err := m.runReminder(g.mockScreen, ing)
	g.EqualError(err, "error adding the following ingredient=PEAR FIG AVOCADO err=exit status 1")
}

func (g *workflowTest) Test_terminal_runReminder_Pass() {
	m := TerminalFakeWorkflow{}
	ing := recipe.Ingredient{
		UnitSize:       "WATERMELON",
		UnitType:       "CHERRY",
		IngredientName: "PITAYA",
	}
	g.mockScreen.On(updateLabelString, "Added Ingredient: WATERMELON CHERRY PITAYA")
	err := m.runReminder(g.mockScreen, ing)
	g.Nil(err)
}

func (g *workflowTest) Test_NewWorkflow_macWorkflow_Pass() {
	g.mockFileChecker.On("checkWorkflowExists").Return(true, nil)

	wf, err := NewWorkflow(g.mockFileChecker, "darwin")
	g.Nil(err)
	g.Equal(&macWorkflow{}, wf)
}

func (g *workflowTest) Test_NewWorkflow_checkWorkflowExists_Error() {
	g.mockFileChecker.On("checkWorkflowExists").Return(false, fmt.Errorf("file check error"))

	wf, err := NewWorkflow(g.mockFileChecker, "darwin")
	g.EqualError(err, "file check error")
	g.Nil(wf)
}

func (g *workflowTest) Test_NewWorkflow_termWorkflow_Pass() {
	g.mockFileChecker.On("checkWorkflowExists").Return(false, nil)

	wf, err := NewWorkflow(g.mockFileChecker, "darwin")
	g.Nil(err)
	g.Equal(&TerminalFakeWorkflow{}, wf)
}

func (g *workflowTest) Test_NewWorkflow_termWorkflow_workflowPresent_Pass() {
	g.mockFileChecker.On("checkWorkflowExists").Return(true, nil)

	wf, err := NewWorkflow(g.mockFileChecker, "windows")
	g.Nil(err)
	g.Equal(&TerminalFakeWorkflow{}, wf)
}

func (g *workflowTest) Test_checkWorkflowExistsImpl_Present_Pass() {
	g.mockFileChecker.On("stat", workflowName).Return(nil, nil)

	present, err := checkWorkflowExistsImpl(g.mockFileChecker)
	g.Nil(err)
	g.Equal(true, present)
}

func (g *workflowTest) Test_checkWorkflowExistsImpl_stat_Error() {
	statError := "stat error"
	g.mockFileChecker.On("stat", workflowName).Return(nil, fmt.Errorf(statError))
	g.mockFileChecker.On("isNotExist", fmt.Errorf(statError)).Return(false)

	present, err := checkWorkflowExistsImpl(g.mockFileChecker)
	g.EqualError(err, statError)
	g.Equal(false, present)
}

func (g *workflowTest) Test_checkWorkflowExistsImpl_NotPresent_Pass() {
	g.mockFileChecker.On("stat", workflowName).Return(nil, fs.ErrNotExist)
	g.mockFileChecker.On("isNotExist", fs.ErrNotExist).Return(true)

	present, err := checkWorkflowExistsImpl(g.mockFileChecker)
	g.Nil(err)
	g.Equal(false, present)
}
