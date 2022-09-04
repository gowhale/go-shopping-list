package recipe

import (
	"fmt"
	"strconv"
)

func CombineRecipesToIngredients(recipes []Recipe) []Ingredient {
	ingsToReturn := []Ingredient{}
	uniqueIngredients := map[string]Ingredient{}
	for _, r := range recipes {
		for _, i := range r.Ings {
			combineTypeName := fmt.Sprintf("%s-%s", i.UnitType, i.IngredientName)
			if _, ok := uniqueIngredients[combineTypeName]; !ok {
				uniqueIngredients[combineTypeName] = i
			} else {
				currentSize, err := strconv.ParseFloat(uniqueIngredients[combineTypeName].UnitSize, 32)
				if err != nil {
					fmt.Println() // 3.1415927410125732
				}

				sizeToAdd, err := strconv.ParseFloat(i.UnitSize, 32)
				if err != nil {
					fmt.Println() // 3.1415927410125732
				}
				oldIng := uniqueIngredients[combineTypeName]
				newSize := fmt.Sprintf("%.2f", currentSize+sizeToAdd)
				oldIng.UnitSize = newSize
				uniqueIngredients[combineTypeName] = oldIng
			}
		}
	}

	for _, ing := range uniqueIngredients {
		ingsToReturn = append(ingsToReturn, ing)
	}

	return ingsToReturn
}
