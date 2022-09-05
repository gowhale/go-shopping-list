// Package gui is responsible for code relating to user interaction
package gui

import (
	"go-shopping-list/pkg/recipe"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/suite"
)

// TODO: Test fyne.io properly()
type guiTest struct {
	suite.Suite
	mockScreen     *mockScreenInterface
	mockWorkflow   *mockWorkflowInterface
	mockFileReader *recipe.MockFileReader
}

func (g *guiTest) SetupTest() {
	g.mockScreen = new(mockScreenInterface)
	g.mockWorkflow = new(mockWorkflowInterface)
	g.mockFileReader = new(recipe.MockFileReader)
}

func TestGuiTest(t *testing.T) {
	suite.Run(t, new(guiTest))
}

func (*guiTest) Test_mockFileInfo() {
	testRecipe := []recipe.Recipe{}
	_ = NewApp(testRecipe, nil, &TerminalFakeWorkflow{})
}

func (g *guiTest) Test_buttonPress() {
	b := createSubmitButton(g.mockScreen, g.mockWorkflow, g.mockFileReader, []string{}, map[string]recipe.Recipe{})
	g.mockScreen.On(updateProgessBarString, progressBarEmpty)
	g.mockScreen.On(updateProgessBarString, progressBarFull)
	g.mockScreen.On(updateLabelString, recipeFinishLabel)
	g.mockWorkflow.On("submitShoppingList", g.mockScreen, g.mockWorkflow, g.mockFileReader, []string{}, map[string]recipe.Recipe{}).Return(nil)
	test.Tap(b)
}

