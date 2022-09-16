// Package gui is responsible for code relating to user interaction
package gui

import (
	"go-shopping-list/pkg/common"
	"go-shopping-list/pkg/recipe"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/suite"
)

const (
	UpdateProgessBarString = "UpdateProgessBar"
	UpdateLabelString      = "UpdateLabel"
	recipeFinishLabel      = "Finished. Select another recipe to add more."
)

// TODO: Test fyne.io properly()
type guiTest struct {
	suite.Suite
	mockScreen     *common.MockScreenInterface
	mockWorkflow   *common.MockWorkflowInterface
	mockFileReader *recipe.MockFileReader
}

func (g *guiTest) SetupTest() {
	g.mockScreen = new(common.MockScreenInterface)
	g.mockWorkflow = new(common.MockWorkflowInterface)
	g.mockFileReader = new(recipe.MockFileReader)
}

func TestGuiTest(t *testing.T) {
	suite.Run(t, new(guiTest))
}

// func (*guiTest) Test_mockFileInfo() {
// 	testRecipe := []recipe.Recipe{}
// 	_ = NewApp(testRecipe, nil, &workflow.TerminalFakeWorkflow{})
// }

func (g *guiTest) Test_buttonPress() {
	b := createSubmitButton(g.mockScreen, g.mockWorkflow, g.mockFileReader, &[]string{}, map[string]recipe.Recipe{})
	g.mockScreen.On(UpdateProgessBarString, progressBarEmpty)
	g.mockScreen.On(UpdateProgessBarString, progressBarFull)
	g.mockScreen.On(UpdateLabelString, recipeFinishLabel)
	g.mockWorkflow.On("SubmitShoppingList", g.mockScreen, g.mockWorkflow, g.mockFileReader, []string{}, map[string]recipe.Recipe{}).Return(nil)
	test.Tap(b)
}
