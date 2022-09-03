package recipe

type Recipe struct {
	Name  string       `json:"recipe_name"`
	Ings  []Ingredient `json:"ingredients"`
	Meth  []string     `json:"method"`
	Count int
}

type Ingredient struct {
	UnitSize       string `json:"unit_size"`
	UnitType       string `json:"unit_type"`
	IngredientName string `json:"ingredient_name"`
}

type Popularity struct {
	Pop []RecipePopularity `json:"popularity"`
}

type RecipePopularity struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
