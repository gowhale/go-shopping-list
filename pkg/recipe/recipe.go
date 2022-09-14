package recipe

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"

	"github.com/bradfitz/slice"
)

const (
	popularityFileName = "popularity.json"
	recipeFolder       = "recipes/"

	errorIntReturn = -1
	defaultCount   = 0

	writePermissionCode = 0644
)

// FileReader deals with interactions with files
//
//go:generate go run github.com/vektra/mockery/cmd/mockery -name FileReader -inpkg --filename file_reader_mock.go
type FileReader interface {
	getPopularity(f FileReader, recipeName string) (int, error)
	IncrementPopularity(f FileReader, recipeName string) error
	loadPopularityFile(f FileReader) (PopularityFile, error)
	marshallJSON(pop PopularityFile) ([]byte, error)
	readRecipeDirectory() ([]fs.FileInfo, error)
	loadRecipeFile(f FileReader, fileName fs.FileInfo) (Recipe, error)
	readFile(filePath string) ([]byte, error)
	unmarshallJSONToPopularity(file []byte) (PopularityFile, error)
	unmarshallJSONToRecipe(file []byte) (Recipe, error)
	writePopularityFile(f FileReader, pop PopularityFile) error
	writeFile(newFile []byte) error
}

// FileInteractionImpl is a struct to implement FileReader
type FileInteractionImpl struct{}

func (*FileInteractionImpl) getPopularity(f FileReader, recipeName string) (int, error) {
	return getPopularityImpl(f, recipeName)
}

// IncrementPopularity incrementes the popularity count of a recipe by 1
func (*FileInteractionImpl) IncrementPopularity(f FileReader, recipeName string) error {
	return incrementPopularityImpl(f, recipeName)
}

func (*FileInteractionImpl) loadPopularityFile(f FileReader) (PopularityFile, error) {
	return loadPopularityFileImpl(f)
}

func (*FileInteractionImpl) readRecipeDirectory() ([]fs.FileInfo, error) {
	return ioutil.ReadDir(recipeFolder)
}

func (*FileInteractionImpl) loadRecipeFile(f FileReader, fileName fs.FileInfo) (Recipe, error) {
	return loadRecipeFileImpl(f, fileName)
}

func (*FileInteractionImpl) readFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func (*FileInteractionImpl) unmarshallJSONToPopularity(file []byte) (PopularityFile, error) {
	popularity := PopularityFile{}
	return popularity, json.Unmarshal([]byte(file), &popularity)
}

func (*FileInteractionImpl) unmarshallJSONToRecipe(file []byte) (Recipe, error) {
	recipe := Recipe{}
	return recipe, json.Unmarshal([]byte(file), &recipe)
}

func (*FileInteractionImpl) marshallJSON(pop PopularityFile) ([]byte, error) {
	return json.MarshalIndent(pop, "", " ")
}

func (*FileInteractionImpl) writeFile(newFile []byte) error {
	return ioutil.WriteFile(popularityFileName, newFile, writePermissionCode)
}

func (*FileInteractionImpl) writePopularityFile(f FileReader, pop PopularityFile) error {
	return writePopularityFileImpl(f, pop)
}

// Implementations of functions
func incrementPopularityImpl(f FileReader, recipeName string) error {
	pop, err := f.loadPopularityFile(f)
	if err != nil {
		return err
	}

	updateIndex := errorIntReturn
	for i, p := range pop.Pop {
		if recipeName == p.Name {
			updateIndex = i
		}
	}
	pop.Pop[updateIndex].Count++

	return f.writePopularityFile(f, pop)
}

func getPopularityImpl(f FileReader, recipeName string) (int, error) {
	pop, err := f.loadPopularityFile(f)
	if err != nil {
		log.Printf("error unmarshalling file=%s", popularityFileName)
		return errorIntReturn, err
	}

	mapOfPops := map[string]int{}
	for _, p := range pop.Pop {
		mapOfPops[p.Name] = p.Count
	}

	if val, ok := mapOfPops[recipeName]; ok {
		return val, nil
	}
	pop.Pop = append(pop.Pop, Popularity{Name: recipeName, Count: defaultCount})
	return defaultCount, f.writePopularityFile(f, pop)
}

func loadPopularityFileImpl(f FileReader) (PopularityFile, error) {
	file, err := f.readFile(popularityFileName)
	if err != nil {
		return PopularityFile{}, err
	}
	return f.unmarshallJSONToPopularity(file)
}

func loadRecipeFileImpl(f FileReader, fileName fs.FileInfo) (Recipe, error) {
	file, err := f.readFile(fmt.Sprintf("recipes/%s", fileName.Name()))
	if err != nil {
		return Recipe{}, err
	}
	return f.unmarshallJSONToRecipe(file)
}

func writePopularityFileImpl(f FileReader, pop PopularityFile) error {
	newFile, err := f.marshallJSON(pop)
	if err != nil {
		return err
	}
	return f.writeFile(newFile)
}

func validateRecipe(f FileReader, uniqueRecipeNames map[string]Recipe, fileName fs.FileInfo) (Recipe, error) {
	recipe, err := f.loadRecipeFile(f, fileName)
	if err != nil {
		return Recipe{}, err
	}

	if recipe.Name == "" {
		return Recipe{}, fmt.Errorf("file-name=%s is missing a recipe name", fileName.Name())
	}

	// Check if duplicate recipe name
	if _, ok := uniqueRecipeNames[recipe.Name]; ok {
		return Recipe{}, fmt.Errorf("duplicate recipe name detected. file-name=%s", fileName.Name())
	}

	if len(recipe.Ings) < 1 {
		return Recipe{}, fmt.Errorf("file-name=%s has 0 ingredients", fileName.Name())
	}

	for _, ing := range recipe.Ings {
		if ing.IngredientName == "" {
			return Recipe{}, fmt.Errorf("file-name=%s ing=%s with nil name", fileName.Name(), ing)
		}
	}

	if len(recipe.Ings) < 1 {
		return Recipe{}, fmt.Errorf("recipe=%s has 0 ingredients", recipe.Name)
	}

	for _, ing := range recipe.Ings {
		if ing.IngredientName == "" {
			return Recipe{}, fmt.Errorf("recipe=%s ing=%s with nil name", recipe.Name, ing)
		}
	}

	recipe.Count, err = f.getPopularity(f, recipe.Name)
	if err != nil {
		return Recipe{}, err
	}
	return recipe, nil
}

// ProcessRecipes processes recipe JSON files from the recipe folder
func ProcessRecipes(f FileReader) ([]Recipe, map[string]Recipe, error) {
	files, err := f.readRecipeDirectory()
	if err != nil {
		return nil, nil, err
	}

	uniqueRecipeNames := map[string]Recipe{}

	// Process every file and put into Recipe strucr
	allRecipes := []Recipe{}
	for _, fileName := range files {
		if !fileName.IsDir() {
			recipe, err := validateRecipe(f, uniqueRecipeNames, fileName)
			if err != nil {
				return nil, nil, err
			}
			allRecipes = append(allRecipes, recipe)
			uniqueRecipeNames[recipe.Name] = recipe
		}
	}

	slice.Sort(allRecipes[:], func(i, j int) bool {
		return allRecipes[i].Count > allRecipes[j].Count
	})

	log.Printf("amount of recipes=%d", len(allRecipes))
	return allRecipes, uniqueRecipeNames, nil
}
