package recipe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"fyne.io/fyne/v2/widget"
	"github.com/bradfitz/slice"
)

const (
	popularityFileName = "popularity.json"
)

type Recipe struct {
	Name  string        `json:"recipe_name"`
	Ings  []Ingredients `json:"ingredients"`
	Meth  []string      `json:"method"`
	Count int
}

type Ingredients struct {
	Unit_size       string `json:"unit_size"`
	Unit_type       string `json:"unit_type"`
	Ingredient_name string `json:"ingredient_name"`
}

type Popularity struct {
	Pop []RecipePopularity `json:"popularity"`
}

type RecipePopularity struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (i *Ingredients) String() string {
	return fmt.Sprintf("%s %s %s", i.Unit_size, i.Unit_type, i.Ingredient_name)
}

func IncrementPopularity(recipeName string) error {
	file, err := ioutil.ReadFile(popularityFileName)
	if err != nil {
		return err
	}
	pop := Popularity{}
	err = json.Unmarshal([]byte(file), &pop)
	if err != nil {
		log.Printf("error unmarshalling file=%s", popularityFileName)
		return err
	}
	updateIndex := -1
	for i, p := range pop.Pop {
		if recipeName == p.Name {
			updateIndex = i
		}
	}
	pop.Pop[updateIndex].Count++
	newFile, _ := json.MarshalIndent(pop, "", " ")

	_ = ioutil.WriteFile(popularityFileName, newFile, 0644)
	return nil
}

func GetPopularity(recipeName string) (int, error) {
	file, err := ioutil.ReadFile(popularityFileName)
	if err != nil {
		return -1, err
	}
	pop := Popularity{}
	err = json.Unmarshal([]byte(file), &pop)
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
	newFile, _ := json.MarshalIndent(pop, "", " ")

	_ = ioutil.WriteFile(popularityFileName, newFile, 0644)
	return 0, nil
}

func ProcessIngredients(recipeFolder string) ([]Recipe, error) {
	allRecipes := []Recipe{}

	// Get name for all recipe files
	files, err := ioutil.ReadDir(recipeFolder)
	if err != nil {
		return nil, err
	}

	// Process every file and put into Recipe strucr
	for _, fileName := range files {
		if !fileName.IsDir() {
			file, err := ioutil.ReadFile(fmt.Sprintf("recipes/%s", fileName.Name()))
			if err != nil {
				return nil, err
			}
			recipe := Recipe{}
			err = json.Unmarshal([]byte(file), &recipe)
			if err != nil {
				log.Printf("error for file=%s", fileName)
				return nil, err
			}
			recipe.Count, err = GetPopularity(recipe.Name)
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

var execCommand = exec.Command

func AddIngredientsToReminders(r Recipe, p *widget.ProgressBar) error {
	if err := IncrementPopularity(r.Name); err != nil {
		return err
	}
	progress := 0.0
	for i, ing := range r.Ings {
		progress = float64(i) / float64(len(r.Ings))
		p.SetValue(progress)
		log.Printf("progress=%.2f adding ing='%s'", progress, ing.String())
		cmd := execCommand("automator", "-i", fmt.Sprintf(`"%s"`, ing.String()), "shopping.workflow")
		_, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error adding the following ingredient=%s err=%e", ing.String(), err)
		}
	}
	progress = 1
	log.Printf("progress=%.2f", progress)
	p.SetValue(progress)
	return nil
}
