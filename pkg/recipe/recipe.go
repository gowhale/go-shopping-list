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
	GetPopularity(recipeName string) (int, error)
	IncrementPopularity(recipeName string) error
	LoadPopularityFile() (Popularity, error)
	MarshallJSON(pop Popularity) ([]byte, error)
	ReadRecipeDirectory() ([]fs.FileInfo, error)
	LoadRecipeFile(fileName fs.FileInfo) (Recipe, error)
	ReadFile(filePath string) ([]byte, error)
	UnmarshallJSONToPopularity(file []byte) (Popularity, error)
	UnmarshallJSONToRecipe(file []byte) (Recipe, error)
	WritePopularityFile(pop Popularity) error
	WriteFile(newFile []byte) error
}

type FileInteractionImpl struct{}

func (f *FileInteractionImpl) GetPopularity(recipeName string) (int, error) {
	return getPopularityImpl(f, recipeName)
}

func (f *FileInteractionImpl) IncrementPopularity(recipeName string) error {
	return incrementPopularityImpl(f, recipeName)
}

func (f *FileInteractionImpl) LoadPopularityFile() (Popularity, error) {
	return loadPopularityFileImpl(f)
}

func (*FileInteractionImpl) ReadRecipeDirectory() ([]fs.FileInfo, error) {
	return ioutil.ReadDir(recipeFolder)
}

func (f *FileInteractionImpl) LoadRecipeFile(fileName fs.FileInfo) (Recipe, error) {
	return loadRecipeFileImpl(f, fileName)
}

func (*FileInteractionImpl) ReadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func (*FileInteractionImpl) UnmarshallJSONToPopularity(file []byte) (Popularity, error) {
	popularity := Popularity{}
	return popularity, json.Unmarshal([]byte(file), &popularity)
}

func (*FileInteractionImpl) UnmarshallJSONToRecipe(file []byte) (Recipe, error) {
	recipe := Recipe{}
	return recipe, json.Unmarshal([]byte(file), &recipe)
}

func (*FileInteractionImpl) MarshallJSON(pop Popularity) ([]byte, error) {
	return json.MarshalIndent(pop, "", " ")
}

func (f *FileInteractionImpl) WriteFile(newFile []byte) error {
	return ioutil.WriteFile(popularityFileName, newFile, 0644)
}

func (f *FileInteractionImpl) WritePopularityFile(pop Popularity) error {
	return writePopularityFileImpl(f, pop)
}

// Implementations of functions
func incrementPopularityImpl(f FileReader, recipeName string) error {
	pop, err := f.LoadPopularityFile()
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

	return f.WritePopularityFile(pop)
}

func getPopularityImpl(f FileReader, recipeName string) (int, error) {
	pop, err := f.LoadPopularityFile()
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
	return 0, f.WritePopularityFile(pop)
}

func loadPopularityFileImpl(f FileReader) (Popularity, error) {
	file, err := f.ReadFile(popularityFileName)
	if err != nil {
		return Popularity{}, err
	}
	return f.UnmarshallJSONToPopularity(file)
}

func loadRecipeFileImpl(f FileReader, fileName fs.FileInfo) (Recipe, error) {
	file, err := f.ReadFile(fmt.Sprintf("recipes/%s", fileName.Name()))
	if err != nil {
		return Recipe{}, err
	}
	return f.UnmarshallJSONToRecipe(file)
}

func writePopularityFileImpl(f FileReader, pop Popularity) error {
	newFile, err := f.MarshallJSON(pop)
	if err != nil {
		return err
	}
	return f.WriteFile(newFile)
}

func ProcessIngredients(f FileReader) ([]Recipe, error) {
	files, err := f.ReadRecipeDirectory()
	if err != nil {
		return nil, err
	}

	// Process every file and put into Recipe strucr
	allRecipes := []Recipe{}
	for _, fileName := range files {
		if !fileName.IsDir() {
			recipe, err := f.LoadRecipeFile(fileName)
			if err != nil {
				return nil, err
			}
			recipe.Count, err = f.GetPopularity(recipe.Name)
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

func (i *Ingredients) String() string {
	return fmt.Sprintf("%s %s %s", i.UnitSize, i.UnitType, i.IngredientName)
}
