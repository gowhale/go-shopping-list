package workflow

import (
	"go-shopping-list/pkg/common"
	"go-shopping-list/pkg/recipe"
	"testing"

	"github.com/stretchr/testify/suite"
)

type terminalTest struct {
	suite.Suite
	mockScreen      *common.MockScreenInterface
	mockFileReader  *recipe.MockFileReader
	mockWorkflow    *common.MockWorkflowInterface
	mockFileChecker *mockFileChecker
}

func (g *terminalTest) SetupTest() {
	g.mockScreen = new(common.MockScreenInterface)
	g.mockFileReader = new(recipe.MockFileReader)
	g.mockWorkflow = new(common.MockWorkflowInterface)
	g.mockFileChecker = new(mockFileChecker)
}

func Test_terminalTest(t *testing.T) {
	suite.Run(t, new(terminalTest))
}

func (g *terminalTest) Test_terminal_RunReminder_Pass() {
	m := TerminalFakeWorkflow{}
	ing := recipe.Ingredient{
		UnitSize:       "WATERMELON",
		UnitType:       "CHERRY",
		IngredientName: "PITAYA",
	}
	g.mockScreen.On(UpdateLabelString, "Added Ingredient: WATERMELON CHERRY PITAYA")
	err := m.RunReminder(g.mockScreen, ing)
	g.Nil(err)
}
