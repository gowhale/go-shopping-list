package recipe

import "fmt"

// Recipe represents a recipe JSON file
type Recipe struct {
	Name  string       `json:"recipe_name"`
	Ings  []Ingredient `json:"ingredients"`
	Meth  []string     `json:"method"`
	Count int
}

// PopularityFile represents the popularity.json file
type PopularityFile struct {
	Pop []Popularity `json:"popularity"`
}

// Popularity contains a recipe name and popularity count
type Popularity struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// Ingredient represents an ingredient from a recipe
type Ingredient struct {
	UnitSize       string `json:"unit_size"`
	UnitType       string `json:"unit_type"`
	IngredientName string `json:"ingredient_name"`
}

func (i *Ingredient) String() string {
	return fmt.Sprintf("%s %s %s", i.UnitSize, i.UnitType, i.IngredientName)
}
