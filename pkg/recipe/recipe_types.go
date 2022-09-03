package recipe

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
