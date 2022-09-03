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
)

//go:generate go run github.com/vektra/mockery/cmd/mockery -name FileReader -inpkg --filename file_reader_mock.go
type FileReader interface {
	getPopularity(recipeName string) (int, error)
	IncrementPopularity(recipeName string) error
	loadPopularityFile() (Popularity, error)
	marshallJSON(pop Popularity) ([]byte, error)
	readRecipeDirectory() ([]fs.FileInfo, error)
	loadRecipeFile(fileName fs.FileInfo) (Recipe, error)
	readFile(filePath string) ([]byte, error)
	unmarshallJSONToPopularity(file []byte) (Popularity, error)
	unmarshallJSONToRecipe(file []byte) (Recipe, error)
	writePopularityFile(pop Popularity) error
	writeFile(newFile []byte) error
}

type FileInteractionImpl struct{}

func (f *FileInteractionImpl) getPopularity(recipeName string) (int, error) {
	return getPopularityImpl(f, recipeName)
}

// IncrementPopularity incrementes the popularity count of a recipe by 1
func (f *FileInteractionImpl) IncrementPopularity(recipeName string) error {
	return incrementPopularityImpl(f, recipeName)
}

func (f *FileInteractionImpl) loadPopularityFile() (Popularity, error) {
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

func (*FileInteractionImpl) unmarshallJSONToPopularity(file []byte) (Popularity, error) {
	popularity := Popularity{}
	return popularity, json.Unmarshal([]byte(file), &popularity)
}

func (*FileInteractionImpl) unmarshallJSONToRecipe(file []byte) (Recipe, error) {
	recipe := Recipe{}
	return recipe, json.Unmarshal([]byte(file), &recipe)
}

func (*FileInteractionImpl) marshallJSON(pop Popularity) ([]byte, error) {
	return json.MarshalIndent(pop, "", " ")
}

func (*FileInteractionImpl) writeFile(newFile []byte) error {
	return ioutil.WriteFile(popularityFileName, newFile, 0644)
}

func (f *FileInteractionImpl) writePopularityFile(pop Popularity) error {
	return writePopularityFileImpl(f, pop)
}

// Implementations of functions
func incrementPopularityImpl(f FileReader, recipeName string) error {
	pop, err := f.loadPopularityFile()
	if err != nil {
		return err
	}

	updateIndex := -1
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
		return -1, err
	}

	mapOfPops := map[string]int{}
	for _, p := range pop.Pop {
		log.Printf("name=%s count=%d", p.Name, p.Count)
		mapOfPops[p.Name] = p.Count
	}

	if val, ok := mapOfPops[recipeName]; ok {
		return val, nil
	}
	pop.Pop = append(pop.Pop, RecipePopularity{Name: recipeName, Count: 0})
	return 0, f.writePopularityFile(pop)
}

func loadPopularityFileImpl(f FileReader) (Popularity, error) {
	file, err := f.readFile(popularityFileName)
	if err != nil {
		return Popularity{}, err
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

func writePopularityFileImpl(f FileReader, pop Popularity) error {
	newFile, err := f.marshallJSON(pop)
	if err != nil {
		return err
	}
	return f.writeFile(newFile)
}

func ProcessIngredients(f FileReader) ([]Recipe, error) {
	files, err := f.readRecipeDirectory()
	if err != nil {
		return nil, err
	}

	// Process every file and put into Recipe strucr
	allRecipes := []Recipe{}
	for _, fileName := range files {
		if !fileName.IsDir() {
			recipe, err := f.loadRecipeFile(fileName)
			if err != nil {
				return nil, err
			}
			recipe.Count, err = f.getPopularity(recipe.Name)
			if err != nil {
				return nil, err
			}
			allRecipes = append(allRecipes, recipe)
		}
	}

	slice.Sort(allRecipes[:], func(i, j int) bool {
		return allRecipes[i].Count > allRecipes[j].Count
	})

	log.Printf("amount of recipes=%d", len(allRecipes))
	return allRecipes, nil
}

func (i *Ingredient) String() string {
	return fmt.Sprintf("%s %s %s", i.UnitSize, i.UnitType, i.IngredientName)
}
