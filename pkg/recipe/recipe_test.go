// Package gui is responsible for visual output
// File terminal_test.go tests the terminal.go file
package recipe

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type recipeTest struct {
	suite.Suite

	mockFileReader *MockFileReader
}

func (t *recipeTest) SetupTest() {
	t.mockFileReader = new(MockFileReader)
}

func TestRecipeTest(t *testing.T) {
	suite.Run(t, new(recipeTest))
}

func (r *recipeTest) Test_String() {
	i := Ingredients{
		Unit_size:       "DURIAN",
		Unit_type:       "APPLE",
		Ingredient_name: "PEAR",
	}
	r.Equal("DURIAN APPLE PEAR", i.String())
}

func (r *recipeTest) Test_GetPopularityImpl_Present_Pass() {
	rp := []RecipePopularity{{
		Name:  "DURIAN",
		Count: 5,
	}}
	mockPop := Popularity{
		Pop: rp,
	}

	r.mockFileReader.On("LoadPopularityFile").Return().Return(mockPop, nil)
	pop, err := GetPopularityImpl(r.mockFileReader, "DURIAN")
	r.Equal(pop, 5)
	r.Nil(err)
}

func (r *recipeTest) Test_GetPopularityImpl_NotPresent_Pass() {
	rp := []RecipePopularity{{
		Name:  "Apple",
		Count: 5,
	}}
	mockPop := Popularity{
		Pop: rp,
	}
	writePopularity := append(rp, RecipePopularity{Name: "DURIAN", Count: 0})

	r.mockFileReader.On("LoadPopularityFile").Return().Return(mockPop, nil)
	r.mockFileReader.On("WritePopularityFile", Popularity{
		Pop: writePopularity,
	}).Return(nil)
	pop, err := GetPopularityImpl(r.mockFileReader, "DURIAN")
	r.Equal(0, pop)
	r.Nil(err)
}
