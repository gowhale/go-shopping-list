package main

import (
	"go-shopping-list/pkg/recipe"
	"log"
)

func main() {
	if _, _, err := recipe.ProcessRecipes(&recipe.FileInteractionImpl{}); err != nil {
		log.Fatalf("error getting all recipes err=%s", err)
	}
	log.Println("Recipes valid")
}
