package recipe

import (
	"fmt"
	"log"
	"strconv"
)

func CombineRecipesToIngredients(recipes []Recipe) []Ingredient {
	uniqueIngredients := map[string]Ingredient{}
	for _, r := range recipes {
		for _, i := range r.Ings {

			// unique identifier is type and name, as could have grams, ounces, kg's of same ing
			combineTypeName := fmt.Sprintf("%s-%s", i.UnitType, i.IngredientName)
			if _, ok := uniqueIngredients[combineTypeName]; !ok {
				uniqueIngredients[combineTypeName] = i
			} else {
				currentSize, err := strconv.ParseFloat(uniqueIngredients[combineTypeName].UnitSize, 32)
				if err != nil {
				}

				sizeToAdd, err := strconv.ParseFloat(i.UnitSize, 32)
				if err != nil {
					log.Fatalln(err)
				}
				oldIng := uniqueIngredients[combineTypeName]
				newSize := fmt.Sprintf("%.2f", currentSize+sizeToAdd)
				oldIng.UnitSize = newSize
				uniqueIngredients[combineTypeName] = oldIng
			}
		}
	}

	ingsToReturn := []Ingredient{}
	for _, ing := range uniqueIngredients {
		ingsToReturn = append(ingsToReturn, ing)
	}

	return ingsToReturn
}
