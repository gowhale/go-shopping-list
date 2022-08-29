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
)

func (i *Ingredients) String() string {
	return fmt.Sprintf("%s %s %s", i.Unit_size, i.Unit_type, i.Ingredient_name)
}

func (f *FileInteractionImpl) LoadPopularityFile() (Popularity, error) {
	file, err := ioutil.ReadFile(popularityFileName)
	if err != nil {
		return Popularity{}, err
	}
	pop := Popularity{}
	return pop, json.Unmarshal([]byte(file), &pop)
}

func (f *FileInteractionImpl) WritePopularityFile(pop Popularity) error {
	newFile, err := json.MarshalIndent(pop, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(popularityFileName, newFile, 0644)
}

func (f *FileInteractionImpl) IncrementPopularity(recipeName string) error {
	return IncrementPopularity(f, recipeName)
}

func IncrementPopularity(f FileReader, recipeName string) error {
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

//go:generate go run github.com/vektra/mockery/cmd/mockery -name FileReader -inpkg --filename file_reader_mock.go
type FileReader interface {
	ReadRecipeDirectory(recipeFolder string) ([]fs.FileInfo, error)
	ReadRecipeFile(fileName fs.FileInfo) (Recipe, error)
	ReadFile(fileName fs.FileInfo) ([]byte, error)
	GetPopularity(recipeName string) (int, error)
	LoadPopularityFile() (Popularity, error)
	WritePopularityFile(pop Popularity) error
	UnmarshallJSON(file []byte) (Recipe, error)
}

type FileInteractionImpl struct{}

func (f *FileInteractionImpl) GetPopularity(recipeName string) (int, error) {
	return GetPopularityImpl(f, recipeName)
}

func GetPopularityImpl(f FileReader, recipeName string) (int, error) {
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

func (*FileInteractionImpl) ReadRecipeDirectory(recipeFolder string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(recipeFolder)
}

func (*FileInteractionImpl) ReadFile(fileName fs.FileInfo) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("recipes/%s", fileName.Name()))
}

func (*FileInteractionImpl) UnmarshallJSON(file []byte) (Recipe, error) {
	recipe := Recipe{}
	return recipe, json.Unmarshal([]byte(file), &recipe)
}

func (f *FileInteractionImpl) ReadRecipeFile(fileName fs.FileInfo) (Recipe, error) {
	return ReadRecipeFile(f, fileName)
}

func ReadRecipeFile(f FileReader, fileName fs.FileInfo) (Recipe, error) {
	file, err := f.ReadFile(fileName)
	if err != nil {
		return Recipe{}, err
	}
	return f.UnmarshallJSON(file)
}

func ProcessIngredients(f FileReader, recipeFolder string) ([]Recipe, error) {
	files, err := f.ReadRecipeDirectory(recipeFolder)
	if err != nil {
		return nil, err
	}

	// Process every file and put into Recipe strucr
	allRecipes := []Recipe{}
	for _, fileName := range files {
		if !fileName.IsDir() {
			recipe, err := f.ReadRecipeFile(fileName)
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
