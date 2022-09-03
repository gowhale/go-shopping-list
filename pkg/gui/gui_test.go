// Package gui is responsible for code relating to user interaction
package gui

import (
	"go-shopping-list/pkg/recipe"
	"testing"

	"github.com/stretchr/testify/suite"
)

type guiTest struct {
	suite.Suite
}

func (g *guiTest) SetupTest() {

}

func TestGuiTest(t *testing.T) {
	suite.Run(t, new(guiTest))
}

func (*guiTest) Test_mockFileInfo() {
	testRecipe := []recipe.Recipe{}
	_ = NewApp(testRecipe, &TerminalFakeWorkflow{})
}
