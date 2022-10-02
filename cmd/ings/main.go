package main

import (
	"go-shopping-list/pkg/gui"
	"go-shopping-list/pkg/recipe"
	"log"

	"github.com/bradfitz/slice"
)

func main() {
	log.Println("Hello World!")

	uniqueIngredients := map[string]bool{}
	recipes, _, err := recipe.ProcessRecipes(&recipe.FileInteractionImpl{})
	if err != nil {
		log.Panicln(err)
	}

	uniqueIngredientsSlice := []string{}
	for _, r := range recipes {
		for _, i := range r.Ings {
			if ok := uniqueIngredients[i.IngredientName]; !ok {
				uniqueIngredientsSlice = append(uniqueIngredientsSlice, i.IngredientName)
				uniqueIngredients[i.IngredientName] = true
			}
		}
	}

	for count, i := range uniqueIngredientsSlice {
		log.Println(count, i)
	}

	slice.Sort(uniqueIngredientsSlice[:], func(i, j int) bool { //nolint:all
		return uniqueIngredientsSlice[i] > uniqueIngredientsSlice[j]
	})

	screen := gui.NewStockScreen(uniqueIngredientsSlice)
	screen.ShowAndRun()
}
