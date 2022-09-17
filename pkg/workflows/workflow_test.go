package workflow

import (
	"fmt"
	"go-shopping-list/pkg/common"
	"go-shopping-list/pkg/fruit"
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

	incrementPopularityString       = "IncrementPopularity"
	RunReminderString               = "RunReminder"
	UpdateLabelString               = "UpdateLabel"
	UpdateProgessBarString          = "UpdateProgessBar"
	checkWorkflowExistsString       = "checkWorkflowExists"
	statString                      = "stat"
	AddIngredientsToRemindersString = "AddIngredientsToReminders"
	melon                           = "MELON"
	apple                           = "APPLE"
)

type workflowTest struct {
	suite.Suite
	mockScreen      *common.MockScreenInterface
	mockFileReader  *recipe.MockFileReader
	mockWorkflow    *common.MockWorkflowInterface
	mockFileChecker *mockFileChecker
}

func (g *workflowTest) SetupTest() {
	g.mockScreen = new(common.MockScreenInterface)
	g.mockFileReader = new(recipe.MockFileReader)
	g.mockWorkflow = new(common.MockWorkflowInterface)
	g.mockFileChecker = new(mockFileChecker)
}

func Test_workflowTest(t *testing.T) {
	suite.Run(t, new(workflowTest))
}

// func (*workflowTest) Test_mockFileInfo() {
// 	testRecipe := []recipe.Recipe{}
// 	_ = NewApp(testRecipe, nil, &TerminalFakeWorkflow{})
// }

func (g *workflowTest) Test_AddIngredientsToReminders_Pass() {
	ings := []recipe.Ingredient{{
		UnitSize:       "PEAR",
		UnitType:       "BANANA",
		IngredientName: "LYCHE",
	}}
	g.mockScreen.On(UpdateLabelString, "Starting to add ingredients for Recipe: DURIAN")
	g.mockScreen.On(UpdateProgessBarString, progressBarEmpty)
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockScreen.On(UpdateLabelString, recipeFinishLabel)
	g.mockWorkflow.On(RunReminderString, g.mockScreen, ings[firstArgumentInSlice]).Return(nil)
	g.mockFileReader.On(incrementPopularityString, "DURIAN").Return(nil)

	err := AddIngredientsToReminders(ings, g.mockScreen, g.mockWorkflow)
	g.Nil(err)

	t := TerminalFakeWorkflow{}
	err = t.AddIngredientsToReminders(ings, g.mockScreen, g.mockWorkflow)
	g.Nil(err)

	m := macWorkflow{}
	err = m.AddIngredientsToReminders(ings, g.mockScreen, g.mockWorkflow)
	g.Nil(err)
}

func (g *workflowTest) Test_AddIngredientsToReminders_RunReminder_Error() {
	ings := []recipe.Ingredient{{
		UnitSize:       "ORANGE",
		UnitType:       "BANANA",
		IngredientName: "RASPBERRY",
	}}

	g.mockScreen.On(UpdateLabelString, "Starting to add ingredients for Recipe: APPLE")
	g.mockScreen.On(UpdateProgessBarString, progressBarEmpty)
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockScreen.On(UpdateLabelString, recipeFinishLabel)
	g.mockWorkflow.On(RunReminderString, g.mockScreen, ings[firstArgumentInSlice]).Return(fmt.Errorf("reminder error"))
	g.mockFileReader.On(incrementPopularityString, apple).Return(nil)
	err := AddIngredientsToReminders(ings, g.mockScreen, g.mockWorkflow)
	g.EqualError(err, "reminder error")
}

func TestHelperProcess(*testing.T) {
	helper := os.Getenv("GO_WANT_HELPER_PROCESS")
	//pass
	if helper == "1" {
		os.Exit(exitCodePass) //nolint
		return
	}
	//fail
	if helper == "2" {
		os.Exit(exitCodeFail) //nolint
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

func (g *workflowTest) Test_macWorkflow_RunReminder_Pass() {
	execCommand = fakeExecCommandPass
	defer func() {
		execCommand = exec.Command
	}()
	m := macWorkflow{}
	ing := recipe.Ingredient{
		UnitSize:       fruit.Watermelon,
		UnitType:       "CHERRY",
		IngredientName: "PITAYA",
	}
	g.mockScreen.On(UpdateLabelString, "Added Ingredient: Watermelon CHERRY PITAYA")
	err := m.RunReminder(g.mockScreen, ing)
	g.Nil(err)
}

func (g *workflowTest) Test_macWorkflow_RunReminder_Error() {
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
	g.mockScreen.On(UpdateLabelString, "Added Ingredient: PEAR FIG AVOCADO")
	err := m.RunReminder(g.mockScreen, ing)
	g.EqualError(err, "error adding the following ingredient=PEAR FIG AVOCADO err=exit status 1")
}

func (g *workflowTest) Test_terminal_RunReminder_Pass() {
	m := TerminalFakeWorkflow{}
	ing := recipe.Ingredient{
		UnitSize:       fruit.Watermelon,
		UnitType:       "CHERRY",
		IngredientName: "PITAYA",
	}
	g.mockScreen.On(UpdateLabelString, "Added Ingredient: WATERMELON CHERRY PITAYA")
	err := m.RunReminder(g.mockScreen, ing)
	g.Nil(err)
}

func (g *workflowTest) Test_NewWorkflow_macWorkflow_Pass() {
	g.mockFileChecker.On(checkWorkflowExistsString, g.mockFileChecker).Return(true, nil)

	wf, err := NewWorkflow(g.mockFileChecker, macOSName)
	g.Nil(err)
	g.Equal(&macWorkflow{}, wf)
}

func (g *workflowTest) Test_NewWorkflow_checkWorkflowExists_Error() {
	g.mockFileChecker.On(checkWorkflowExistsString, g.mockFileChecker).Return(false, fmt.Errorf("file check error"))

	wf, err := NewWorkflow(g.mockFileChecker, macOSName)
	g.EqualError(err, "file check error")
	g.Nil(wf)
}

func (g *workflowTest) Test_NewWorkflow_termWorkflow_Pass() {
	g.mockFileChecker.On(checkWorkflowExistsString, g.mockFileChecker).Return(false, nil)

	wf, err := NewWorkflow(g.mockFileChecker, macOSName)
	g.Nil(err)
	g.Equal(&TerminalFakeWorkflow{}, wf)
}

func (g *workflowTest) Test_NewWorkflow_termWorkflow_workflowPresent_Pass() {
	g.mockFileChecker.On(checkWorkflowExistsString, g.mockFileChecker).Return(true, nil)

	wf, err := NewWorkflow(g.mockFileChecker, "windows")
	g.Nil(err)
	g.Equal(&TerminalFakeWorkflow{}, wf)
}

func (g *workflowTest) Test_checkWorkflowExistsImpl_Present_Pass() {
	g.mockFileChecker.On(statString, workflowName).Return(nil, nil)

	present, err := checkWorkflowExistsImpl(g.mockFileChecker)
	g.Nil(err)
	g.Equal(true, present)

	w := CheckWorkflow{}
	present, err = w.checkWorkflowExists(g.mockFileChecker)
	g.Nil(err)
	g.Equal(true, present)
}

func (g *workflowTest) Test_checkWorkflowExistsImpl_stat_Error() {
	statError := "stat error"
	g.mockFileChecker.On(statString, workflowName).Return(nil, fmt.Errorf(statError))
	g.mockFileChecker.On("isNotExist", fmt.Errorf(statError)).Return(false)

	present, err := checkWorkflowExistsImpl(g.mockFileChecker)
	g.EqualError(err, statError)
	g.Equal(false, present)
}

func (g *workflowTest) Test_checkWorkflowExistsImpl_NotPresent_Pass() {
	g.mockFileChecker.On(statString, workflowName).Return(nil, fs.ErrNotExist)
	g.mockFileChecker.On("isNotExist", fs.ErrNotExist).Return(true)

	present, err := checkWorkflowExistsImpl(g.mockFileChecker)
	g.Nil(err)
	g.Equal(false, present)
}

func (g *workflowTest) Test_SubmitShoppingList_Pass() {
	testRecipe := recipe.Recipe{
		Name: melon,
		Ings: []recipe.Ingredient{
			recipe.Ingredient{
				IngredientName: "PEACH",
				UnitSize:       "1",
				UnitType:       apple,
			},
		},
	}
	recipeMap := map[string]recipe.Recipe{
		testRecipe.Name: testRecipe,
	}
	recipeString := []string{testRecipe.Name}
	g.mockScreen.On(UpdateProgessBarString, progressBarEmpty)
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockScreen.On(UpdateLabelString, recipeFinishLabel)
	g.mockFileReader.On(incrementPopularityString, g.mockFileReader, melon).Return(nil)
	g.mockWorkflow.On(AddIngredientsToRemindersString, []recipe.Ingredient{testRecipe.Ings[firstArgumentInSlice]}, g.mockScreen, g.mockWorkflow).Return(nil)

	err := SubmitShoppingList(g.mockScreen, g.mockWorkflow, g.mockFileReader, recipeString, recipeMap)
	g.Nil(err)

	t := TerminalFakeWorkflow{}
	err = t.SubmitShoppingList(g.mockScreen, g.mockWorkflow, g.mockFileReader, recipeString, recipeMap)
	g.Nil(err)

	m := macWorkflow{}
	err = m.SubmitShoppingList(g.mockScreen, g.mockWorkflow, g.mockFileReader, recipeString, recipeMap)
	g.Nil(err)

	e := excelWorkflow{}
	err = e.SubmitShoppingList(g.mockScreen, g.mockWorkflow, g.mockFileReader, recipeString, recipeMap)
	g.Nil(err)
}

func (g *workflowTest) Test_SubmitShoppingList_Error() {
	testRecipe := recipe.Recipe{
		Name: melon,
		Ings: []recipe.Ingredient{
			recipe.Ingredient{
				IngredientName: "PEACH",
				UnitSize:       "5",
				UnitType:       apple,
			},
		},
	}
	recipeMap := map[string]recipe.Recipe{
		testRecipe.Name: testRecipe,
	}
	recipeString := []string{testRecipe.Name}
	g.mockScreen.On(UpdateProgessBarString, progressBarEmpty)
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockScreen.On(UpdateLabelString, recipeFinishLabel)
	g.mockFileReader.On(incrementPopularityString, g.mockFileReader, melon).Return(fmt.Errorf("increment pop error"))
	g.mockWorkflow.On(AddIngredientsToRemindersString, []recipe.Ingredient{testRecipe.Ings[firstArgumentInSlice]}, g.mockScreen, g.mockWorkflow).Return(nil)
	err := SubmitShoppingList(g.mockScreen, g.mockWorkflow, g.mockFileReader, recipeString, recipeMap)
	g.EqualError(err, "increment pop error")
}
