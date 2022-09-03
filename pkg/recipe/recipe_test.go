// Package gui is responsible for visual output
// File terminal_test.go tests the terminal.go file
package recipe

import (
	"encoding/json"
	"fmt"
	"io/fs"

	"testing"

	"github.com/stretchr/testify/suite"
)

type recipeTest struct {
	suite.Suite

	mockFileReader *MockFileReader
}

func (r *recipeTest) SetupTest() {
	r.mockFileReader = new(MockFileReader)
}

func TestRecipeTest(r *testing.T) {
	suite.Run(r, new(recipeTest))
}

func (r *recipeTest) Test_String() {
	i := Ingredient{
		UnitSize:       "DURIAN",
		UnitType:       "APPLE",
		IngredientName: "PEAR",
	}
	r.Equal("DURIAN APPLE PEAR", i.String())
}

func (r *recipeTest) Test_loadPopularityFileImpl_Pass() {
	testByte := []byte{}
	expectedPop := Popularity{
		Pop: []RecipePopularity{
			{
				Name:  "DURIAN",
				Count: 5,
			},
		},
	}
	r.mockFileReader.On("readFile", popularityFileName).Return(testByte, nil)
	r.mockFileReader.On("unmarshallJSONToPopularity", testByte).Return(expectedPop, nil)
	pop, err := loadPopularityFileImpl(r.mockFileReader)
	r.Nil(err)
	r.Equal(expectedPop, pop)
}

func (r *recipeTest) Test_loadPopularityFileImpl_Error() {
	testByte := []byte{}
	expectedPop := Popularity{
		Pop: []RecipePopularity{
			{
				Name:  "DURIAN",
				Count: 5,
			},
		},
	}
	r.mockFileReader.On("readFile", popularityFileName).Return(testByte, fmt.Errorf("read error"))
	r.mockFileReader.On("unmarshallJSONToPopularity", testByte).Return(expectedPop, nil)
	pop, err := loadPopularityFileImpl(r.mockFileReader)
	r.EqualError(err, "read error")
	r.Equal(Popularity{}, pop)
}

func (r *recipeTest) Test_ReadRecipeFile_Pass() {
	testByte := []byte{}
	expectedRecipe := Recipe{
		Name: "APPLE",
	}
	r.mockFileReader.On("readFile", "recipes/DURIAN").Return(testByte, nil)
	r.mockFileReader.On("unmarshallJSONToRecipe", testByte).Return(expectedRecipe, nil)
	recipe, err := loadRecipeFileImpl(r.mockFileReader, &mockFileInfo{})
	r.Nil(err)
	r.Equal(expectedRecipe, recipe)
}

func (r *recipeTest) Test_ReadRecipeFile_Error() {
	testByte := []byte{}
	r.mockFileReader.On("readFile", "recipes/DURIAN").Return(testByte, fmt.Errorf("read file error"))
	recipe, err := loadRecipeFileImpl(r.mockFileReader, &mockFileInfo{})
	r.EqualError(err, "read file error")
	r.Equal(Recipe{}, recipe)
}

func (r *recipeTest) Test_IncrementPopularity_Pass() {
	rp1 := []RecipePopularity{{
		Name:  "DURIAN",
		Count: 5,
	}}
	mockPopBefore := Popularity{
		Pop: rp1,
	}
	rp2 := []RecipePopularity{{
		Name:  "DURIAN",
		Count: 6,
	}}
	mockPopAfter := Popularity{
		Pop: rp2,
	}

	r.mockFileReader.On("loadPopularityFile").Return().Return(mockPopBefore, nil)
	r.mockFileReader.On("writePopularityFile", mockPopAfter).Return(nil)
	err := incrementPopularityImpl(r.mockFileReader, "DURIAN")
	r.Nil(err)
}

func (r *recipeTest) Test_IncrementPopularity_Error() {
	r.mockFileReader.On("loadPopularityFile").Return().Return(Popularity{}, fmt.Errorf("load error"))
	err := incrementPopularityImpl(r.mockFileReader, "DURIAN")
	r.EqualError(err, "load error")
}

func (r *recipeTest) Test_GetPopularityImpl_Present_Pass() {
	rp := []RecipePopularity{{
		Name:  "DURIAN",
		Count: 5,
	}}
	mockPop := Popularity{
		Pop: rp,
	}

	r.mockFileReader.On("loadPopularityFile").Return().Return(mockPop, nil)
	pop, err := getPopularityImpl(r.mockFileReader, "DURIAN")
	r.Equal(pop, 5)
	r.Nil(err)
}

func (r *recipeTest) Test_GetPopularityImpl_Error() {
	rp := []RecipePopularity{{
		Name:  "DURIAN",
		Count: 5,
	}}
	mockPop := Popularity{
		Pop: rp,
	}

	r.mockFileReader.On("loadPopularityFile").Return().Return(mockPop, fmt.Errorf("load file error"))
	pop, err := getPopularityImpl(r.mockFileReader, "DURIAN")
	r.Equal(pop, -1)
	r.EqualError(err, "load file error")
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

	r.mockFileReader.On("loadPopularityFile").Return(mockPop, nil)
	r.mockFileReader.On("writePopularityFile", Popularity{
		Pop: writePopularity,
	}).Return(nil)
	pop, err := getPopularityImpl(r.mockFileReader, "DURIAN")
	r.Equal(0, pop)
	r.Nil(err)
}

func (r *recipeTest) Test_WritePopularityFileImpl_Pass() {
	expectedRecipe := Popularity{}

	r.mockFileReader.On("marshallJSON", expectedRecipe).Return([]byte{}, nil)
	r.mockFileReader.On("writeFile", []byte{}).Return(nil)
	err := writePopularityFileImpl(r.mockFileReader, expectedRecipe)
	r.Nil(err)
}

func (r *recipeTest) Test_WritePopularityFileImpl_WriteFile_Error() {
	expectedRecipe := Popularity{}

	r.mockFileReader.On("marshallJSON", expectedRecipe).Return([]byte{}, nil)
	r.mockFileReader.On("writeFile", []byte{}).Return(fmt.Errorf("write error"))
	err := writePopularityFileImpl(r.mockFileReader, expectedRecipe)
	r.EqualError(err, "write error")
}

func (r *recipeTest) Test_WritePopularityFileImpl_MarshallJSON_Error() {
	expectedRecipe := Popularity{}

	r.mockFileReader.On("marshallJSON", expectedRecipe).Return([]byte{}, fmt.Errorf("marshall error"))
	err := writePopularityFileImpl(r.mockFileReader, expectedRecipe)
	r.EqualError(err, "marshall error")
}

func (r *recipeTest) Test_ProcessIngredients_Pass() {
	expectedRecipe := Recipe{
		Name:  "APPLE",
		Count: 5,
	}
	expectedResult := []Recipe{expectedRecipe, expectedRecipe}
	filesReturned := []fs.FileInfo{&mockFileInfo{}, &mockFileInfo{}}
	r.mockFileReader.On("readRecipeDirectory").Return(filesReturned, nil)
	r.mockFileReader.On("loadRecipeFile", &mockFileInfo{}).Return(expectedRecipe, nil)
	r.mockFileReader.On("getPopularity", expectedRecipe.Name).Return(5, nil)
	recipes, err := ProcessIngredients(r.mockFileReader)
	r.Nil(err)
	r.Equal(expectedResult, recipes)
}

func (r *recipeTest) Test_ProcessIngredients_GetPopularity_Error() {
	expectedRecipe := Recipe{
		Name:  "APPLE",
		Count: 5,
	}
	filesReturned := []fs.FileInfo{&mockFileInfo{}, &mockFileInfo{}}
	r.mockFileReader.On("readRecipeDirectory").Return(filesReturned, nil)
	r.mockFileReader.On("loadRecipeFile", &mockFileInfo{}).Return(expectedRecipe, nil)
	r.mockFileReader.On("getPopularity", expectedRecipe.Name).Return(-1, fmt.Errorf("pop error"))
	recipes, err := ProcessIngredients(r.mockFileReader)
	r.Nil(recipes)
	r.EqualError(err, "pop error")
}

func (r *recipeTest) Test_ProcessIngredients_ReadRecipeFile_Error() {
	filesReturned := []fs.FileInfo{&mockFileInfo{}, &mockFileInfo{}}
	r.mockFileReader.On("readRecipeDirectory").Return(filesReturned, nil)
	r.mockFileReader.On("loadRecipeFile", &mockFileInfo{}).Return(Recipe{}, fmt.Errorf("read file error"))
	recipes, err := ProcessIngredients(r.mockFileReader)
	r.Nil(recipes)
	r.EqualError(err, "read file error")
}

func (r *recipeTest) Test_ReadRecipeDirectory_Error() {
	r.mockFileReader.On("readRecipeDirectory").Return(nil, fmt.Errorf("read dir error"))
	recipes, err := ProcessIngredients(r.mockFileReader)
	r.Nil(recipes)
	r.EqualError(err, "read dir error")
}

func (r *recipeTest) Test_UnmarshallJSONToRecipe_Pass() {
	mf := FileInteractionImpl{}
	expected := Recipe{
		Name: "DURIAN",
		Ings: nil,
		Meth: nil,
	}
	bytes, err := json.Marshal(expected)
	r.Nil(err)
	result, err := mf.unmarshallJSONToRecipe(bytes)
	r.Equal(expected, result)
	r.Nil(err)
}

func (r *recipeTest) Test_UnmarshallJSONToPopularity_Pass() {
	mf := FileInteractionImpl{}
	rp := []RecipePopularity{{
		Name:  "Apple",
		Count: 5,
	}}
	mockPop := Popularity{
		Pop: rp,
	}
	bytes, err := json.Marshal(mockPop)
	r.Nil(err)
	result, err := mf.unmarshallJSONToPopularity(bytes)
	r.Equal(mockPop, result)
	r.Nil(err)
}
