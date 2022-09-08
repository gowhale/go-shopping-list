package recipe

import (
	"fmt"
	"strconv"
)

func convertMapToSlice(im map[string]Ingredient) []Ingredient {
	ingsToReturn := []Ingredient{}
	for _, ing := range im {
		ingsToReturn = append(ingsToReturn, ing)
	}
	return ingsToReturn
}

func calculateNewSize(uniqueIngredients map[string]Ingredient, combineTypeName, currentIngSize string) (string, error) {
	currentSize, err := strconv.ParseFloat(uniqueIngredients[combineTypeName].UnitSize, 32)
	if err != nil {
		return "", err
	}

	sizeToAdd, err := strconv.ParseFloat(currentIngSize, 32)
	if err != nil {
		return "nil", err
	}
	return fmt.Sprintf("%.2f", currentSize+sizeToAdd), nil
}

// CombineRecipesToIngredients combines the ingredients within mutiple recipes
func CombineRecipesToIngredients(recipes []Recipe) ([]Ingredient, error) {
	uniqueIngredients := map[string]Ingredient{}
	for _, r := range recipes {
		for _, i := range r.Ings {

			// unique identifier is type and name, as could have grams, ounces, kg's of same ing
			combineTypeName := fmt.Sprintf("%s-%s", i.UnitType, i.IngredientName)
			if _, ok := uniqueIngredients[combineTypeName]; !ok {
				uniqueIngredients[combineTypeName] = i
			} else {
				oldIng := uniqueIngredients[combineTypeName]
				newSize, err := calculateNewSize(uniqueIngredients, combineTypeName, i.UnitSize)
				if err != nil {
					return nil, err
				}
				oldIng.UnitSize = newSize
				uniqueIngredients[combineTypeName] = oldIng
			}
		}
	}

	return convertMapToSlice(uniqueIngredients), nil
}
