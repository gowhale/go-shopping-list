// Package gui is responsible for code relating to user interaction
package gui

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TODO: Test fyne.io properly()
type guiTest struct {
	suite.Suite
}

func (*guiTest) SetupTest() {

}

func TestGuiTest(t *testing.T) {
	suite.Run(t, new(guiTest))
}

// func (*guiTest) Test_mockFileInfo() {
// 	testRecipe := []recipe.Recipe{}
// 	_ = NewApp(testRecipe, nil, &TerminalFakeWorkflow{})
// }
