// Package main runs the shopping list
package main

import (
	"go-shopping-list/pkg/gui"
	"go-shopping-list/pkg/recipe"
	"log"
)

func main() {
	//Fetch Recipes
	allRecipes, err := recipe.ProcessIngredients(&recipe.FileInteractionImpl{}, "recipes/")
	if err != nil {
		log.Fatalf("error getting all recipes err=%e", err)
	}

	// Show Window
	myWindow := gui.NewApp(allRecipes)
	myWindow.ShowAndRun()
}
