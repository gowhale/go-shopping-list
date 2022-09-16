// Package main authenticates that the recipes are in the correct format with valid data
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
