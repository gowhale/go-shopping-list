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
//go:generate go run github.com/vektra/mockery/cmd/mockery -name FileReader -inpkg --filename file_reader_mock.go
type FileReader interface {
	getPopularity(recipeName string) (int, error)
	IncrementPopularity(recipeName string) error
	loadPopularityFile() (PopularityFile, error)
	marshallJSON(pop PopularityFile) ([]byte, error)
	readRecipeDirectory() ([]fs.FileInfo, error)
	loadRecipeFile(fileName fs.FileInfo) (Recipe, error)
	readFile(filePath string) ([]byte, error)
	unmarshallJSONToPopularity(file []byte) (PopularityFile, error)
	unmarshallJSONToRecipe(file []byte) (Recipe, error)
	writePopularityFile(pop PopularityFile) error
	writeFile(newFile []byte) error
}

// FileInteractionImpl is a struct to implement FileReader
type FileInteractionImpl struct{}

func (f *FileInteractionImpl) getPopularity(recipeName string) (int, error) {
	return getPopularityImpl(f, recipeName)
}

// IncrementPopularity incrementes the popularity count of a recipe by 1
func (f *FileInteractionImpl) IncrementPopularity(recipeName string) error {
	return incrementPopularityImpl(f, recipeName)
}

func (f *FileInteractionImpl) loadPopularityFile() (PopularityFile, error) {
	return loadPopularityFileImpl(f)
}

func (*FileInteractionImpl) readRecipeDirectory() ([]fs.FileInfo, error) {
	return ioutil.ReadDir(recipeFolder)
}

func (f *FileInteractionImpl) loadRecipeFile(fileName fs.FileInfo) (Recipe, error) {
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

func (f *FileInteractionImpl) writePopularityFile(pop PopularityFile) error {
	return writePopularityFileImpl(f, pop)
}

// Implementations of functions
func incrementPopularityImpl(f FileReader, recipeName string) error {
	pop, err := f.loadPopularityFile()
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

	return f.writePopularityFile(pop)
}

func getPopularityImpl(f FileReader, recipeName string) (int, error) {
	pop, err := f.loadPopularityFile()
	if err != nil {
		log.Printf("error unmarshalling file=%s", popularityFileName)
		return errorIntReturn, err
	}

	mapOfPops := map[string]int{}
	for _, p := range pop.Pop {
		log.Printf("name=%s count=%d", p.Name, p.Count)
		mapOfPops[p.Name] = p.Count
	}

	if val, ok := mapOfPops[recipeName]; ok {
		return val, nil
	}
	pop.Pop = append(pop.Pop, Popularity{Name: recipeName, Count: defaultCount})
	return defaultCount, f.writePopularityFile(pop)
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
			recipe, err := f.loadRecipeFile(fileName)
			if err != nil {
				return nil, nil, err
			}
			
			// Check if duplicate recipe name
			if _, ok := uniqueRecipeNames[recipe.Name]; ok {
				return nil, nil, fmt.Errorf("duplicate recipe name detected. name=%s", recipe.Name)
			}

			recipe.Count, err = f.getPopularity(recipe.Name)
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
