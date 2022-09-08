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

const (
	marshallJSONString        = "marshallJSON"
	loadRecipeFileString      = "loadRecipeFile"
	readRecipeDirectoryString = "readRecipeDirectory"
	loadPopularityFileString  = "loadPopularityFile"
	readFileString            = "readFile"

	readFileError       = "read file error"
	testCount           = 5
	countAfterIncrement = 6
)

type recipeTest struct {
	suite.Suite

	mockFileReader *MockFileReader
	fileImpl       FileInteractionImpl
}

func (r *recipeTest) SetupTest() {
	r.mockFileReader = new(MockFileReader)
	r.fileImpl = FileInteractionImpl{}
}

func TestRecipeTest(r *testing.T) {
	suite.Run(r, new(recipeTest))
}

func (r *recipeTest) Test_String() {
	i := Ingredient{
		UnitSize:       "PEACH",
		UnitType:       "WATERMELON",
		IngredientName: "PEAR",
	}
	r.Equal("PEACH WATERMELON PEAR", i.String())
}

func (r *recipeTest) Test_loadPopularityFileImpl_Pass() {
	testByte := []byte{}
	expectedPop := PopularityFile{
		Pop: []Popularity{
			{
				Name:  "LEMON",
				Count: testCount,
			},
		},
	}
	r.mockFileReader.On(readFileString, popularityFileName).Return(testByte, nil)
	r.mockFileReader.On("unmarshallJSONToPopularity", testByte).Return(expectedPop, nil)

	pop, err := loadPopularityFileImpl(r.mockFileReader)
	r.Nil(err)
	r.Equal(expectedPop, pop)

	pop, err = r.fileImpl.loadPopularityFile(r.mockFileReader)
	r.Nil(err)
	r.Equal(expectedPop, pop)
}

func (r *recipeTest) Test_loadPopularityFileImpl_Error() {
	testByte := []byte{}
	expectedPop := PopularityFile{
		Pop: []Popularity{
			{
				Name:  "CRANBERRY",
				Count: testCount,
			},
		},
	}
	r.mockFileReader.On(readFileString, popularityFileName).Return(testByte, fmt.Errorf("read error"))
	r.mockFileReader.On("unmarshallJSONToPopularity", testByte).Return(expectedPop, nil)
	pop, err := loadPopularityFileImpl(r.mockFileReader)
	r.EqualError(err, "read error")
	r.Equal(PopularityFile{}, pop)
}

func (r *recipeTest) Test_loadRecipeFileImpl_Pass() {
	testByte := []byte{}
	expectedRecipe := Recipe{
		Name: "APPLE",
	}
	r.mockFileReader.On(readFileString, "recipes/DURIAN").Return(testByte, nil)
	r.mockFileReader.On("unmarshallJSONToRecipe", testByte).Return(expectedRecipe, nil)

	recipe, err := loadRecipeFileImpl(r.mockFileReader, &mockFileInfo{})
	r.Nil(err)
	r.Equal(expectedRecipe, recipe)

	recipe, err = r.fileImpl.loadRecipeFile(r.mockFileReader, &mockFileInfo{})
	r.Nil(err)
	r.Equal(expectedRecipe, recipe)
}

func (r *recipeTest) Test_ReadRecipeFile_Error() {
	testByte := []byte{}
	r.mockFileReader.On(readFileString, "recipes/DURIAN").Return(testByte, fmt.Errorf(readFileError))
	recipe, err := loadRecipeFileImpl(r.mockFileReader, &mockFileInfo{})
	r.EqualError(err, readFileError)
	r.Equal(Recipe{}, recipe)
}

func (r *recipeTest) Test_IncrementPopularity_Pass() {
	testFruit := "GUAVA"
	rp1 := []Popularity{{
		Name:  testFruit,
		Count: testCount,
	}}
	mockPopBefore := PopularityFile{
		Pop: rp1,
	}
	rp2 := []Popularity{{
		Name:  testFruit,
		Count: countAfterIncrement,
	}}
	mockPopAfter := PopularityFile{
		Pop: rp2,
	}

	r.mockFileReader.On(loadPopularityFileString, r.mockFileReader).Return().Return(mockPopBefore, nil)
	r.mockFileReader.On("writePopularityFile", r.mockFileReader, mockPopAfter).Return(nil)

	err := incrementPopularityImpl(r.mockFileReader, testFruit)
	r.Nil(err)

	rp3 := []Popularity{{
		Name:  testFruit,
		Count: 7,
	}}

	mockPopAfter2 := PopularityFile{
		Pop: rp3,
	}
	r.mockFileReader.On("writePopularityFile", r.mockFileReader, mockPopAfter2).Return(nil)

	err = r.fileImpl.IncrementPopularity(r.mockFileReader, testFruit)
	r.Nil(err)
}

func (r *recipeTest) Test_IncrementPopularity_Error() {
	r.mockFileReader.On(loadPopularityFileString, r.mockFileReader).Return().Return(PopularityFile{}, fmt.Errorf("load error"))
	err := incrementPopularityImpl(r.mockFileReader, "STRAWBERRY")
	r.EqualError(err, "load error")
}

func (r *recipeTest) Test_GetPopularityImpl_Present_Pass() {
	rp := []Popularity{{
		Name:  "JACKFRUIT",
		Count: testCount,
	}}
	mockPop := PopularityFile{
		Pop: rp,
	}

	r.mockFileReader.On(loadPopularityFileString, r.mockFileReader).Return().Return(mockPop, nil)

	pop, err := getPopularityImpl(r.mockFileReader, "JACKFRUIT")
	r.Equal(pop, testCount)
	r.Nil(err)

	pop, err = r.fileImpl.getPopularity(r.mockFileReader, "JACKFRUIT")
	r.Equal(pop, testCount)
	r.Nil(err)
}

func (r *recipeTest) Test_GetPopularityImpl_Error() {
	rp := []Popularity{{
		Name:  "PLANTAIN",
		Count: testCount,
	}}
	mockPop := PopularityFile{
		Pop: rp,
	}

	r.mockFileReader.On(loadPopularityFileString, r.mockFileReader).Return().Return(mockPop, fmt.Errorf("load file error"))
	pop, err := getPopularityImpl(r.mockFileReader, "PLANTAIN")
	r.Equal(pop, errorIntReturn)
	r.EqualError(err, "load file error")
}

func (r *recipeTest) Test_GetPopularityImpl_NotPresent_Pass() {
	rp := []Popularity{{
		Name:  "Apple",
		Count: testCount,
	}}
	mockPop := PopularityFile{
		Pop: rp,
	}
	writePopularity := append(rp, Popularity{Name: "LIME", Count: defaultCount})

	r.mockFileReader.On(loadPopularityFileString, r.mockFileReader).Return(mockPop, nil)
	r.mockFileReader.On("writePopularityFile", r.mockFileReader, PopularityFile{
		Pop: writePopularity,
	}).Return(nil)
	pop, err := getPopularityImpl(r.mockFileReader, "LIME")
	r.Equal(defaultCount, pop)
	r.Nil(err)
}

func (r *recipeTest) Test_WritePopularityFileImpl_Pass() {
	expectedRecipe := PopularityFile{}

	r.mockFileReader.On(marshallJSONString, expectedRecipe).Return([]byte{}, nil)
	r.mockFileReader.On("writeFile", []byte{}).Return(nil)

	err := writePopularityFileImpl(r.mockFileReader, expectedRecipe)
	r.Nil(err)

	err = r.fileImpl.writePopularityFile(r.mockFileReader, expectedRecipe)
	r.Nil(err)
}

func (r *recipeTest) Test_WritePopularityFileImpl_WriteFile_Error() {
	expectedRecipe := PopularityFile{}

	r.mockFileReader.On(marshallJSONString, expectedRecipe).Return([]byte{}, nil)
	r.mockFileReader.On("writeFile", []byte{}).Return(fmt.Errorf("write error"))
	err := writePopularityFileImpl(r.mockFileReader, expectedRecipe)
	r.EqualError(err, "write error")
}

func (r *recipeTest) Test_WritePopularityFileImpl_MarshallJSON_Error() {
	expectedRecipe := PopularityFile{}

	r.mockFileReader.On(marshallJSONString, expectedRecipe).Return([]byte{}, fmt.Errorf("marshall error"))
	err := writePopularityFileImpl(r.mockFileReader, expectedRecipe)
	r.EqualError(err, "marshall error")
}

func (r *recipeTest) Test_ProcessIngredients_Pass() {
	r1 := Recipe{
		Name:  "APPLE",
		Count: testCount,
		Ings: []Ingredient{
			Ingredient{
				UnitSize:       "1",
				UnitType:       "g",
				IngredientName: "salt",
			},
		},
	}
	r2 := Recipe{
		Name:  "PEAR",
		Count: testCount,
		Ings: []Ingredient{
			Ingredient{
				UnitSize:       "1",
				UnitType:       "g",
				IngredientName: "salt",
			},
		},
	}
	expectedResult := []Recipe{r1, r2}
	filesReturned := []fs.FileInfo{&mockFileInfo{}, &mockFileInfo{}}
	r.mockFileReader.On(readRecipeDirectoryString).Return(filesReturned, nil)
	r.mockFileReader.On(loadRecipeFileString, r.mockFileReader, &mockFileInfo{}).Return(r1, nil).Once()
	r.mockFileReader.On(loadRecipeFileString, r.mockFileReader, &mockFileInfo{}).Return(r2, nil).Once()
	r.mockFileReader.On("getPopularity", r.mockFileReader, r1.Name).Return(testCount, nil)
	r.mockFileReader.On("getPopularity", r.mockFileReader, r2.Name).Return(testCount, nil)
	recipes, _, err := ProcessRecipes(r.mockFileReader)
	r.Nil(err)
	r.Equal(expectedResult, recipes)
}

func (r *recipeTest) Test_ProcessIngredients_GetPopularity_Error() {
	expectedRecipe := Recipe{
		Name:  "DURIAN",
		Count: testCount,
	}
	filesReturned := []fs.FileInfo{&mockFileInfo{}, &mockFileInfo{}}
	r.mockFileReader.On(readRecipeDirectoryString).Return(filesReturned, nil)
	r.mockFileReader.On(loadRecipeFileString, r.mockFileReader, &mockFileInfo{}).Return(expectedRecipe, nil)
	r.mockFileReader.On("getPopularity", r.mockFileReader, expectedRecipe.Name).Return(errorIntReturn, fmt.Errorf("pop error"))
	recipes, _, err := ProcessRecipes(r.mockFileReader)
	r.Nil(recipes)
	r.EqualError(err, "pop error")
}

func (r *recipeTest) Test_ProcessIngredients_ReadRecipeFile_Error() {
	filesReturned := []fs.FileInfo{&mockFileInfo{}, &mockFileInfo{}}
	r.mockFileReader.On(readRecipeDirectoryString).Return(filesReturned, nil)
	r.mockFileReader.On(loadRecipeFileString, r.mockFileReader, &mockFileInfo{}).Return(Recipe{}, fmt.Errorf(readFileError))
	recipes, _, err := ProcessRecipes(r.mockFileReader)
	r.Nil(recipes)
	r.EqualError(err, readFileError)
}

func (r *recipeTest) Test_ReadRecipeDirectory_Error() {
	r.mockFileReader.On(readRecipeDirectoryString).Return(nil, fmt.Errorf("read dir error"))
	recipes, _, err := ProcessRecipes(r.mockFileReader)
	r.Nil(recipes)
	r.EqualError(err, "read dir error")
}

func (r *recipeTest) Test_UnmarshallJSONToRecipe_Pass() {
	mf := FileInteractionImpl{}
	expected := Recipe{
		Name: "LEMON",
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
	rp := []Popularity{{
		Name:  "Apple",
		Count: testCount,
	}}
	mockPop := PopularityFile{
		Pop: rp,
	}
	bytes, err := json.Marshal(mockPop)
	r.Nil(err)
	result, err := mf.unmarshallJSONToPopularity(bytes)
	r.Equal(mockPop, result)
	r.Nil(err)
}
