package recipe

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TODO: Test fyne.io properly()
type ingredientsTest struct {
	suite.Suite
}

func (*ingredientsTest) SetupTest() {

}

func TestGuiTest(t *testing.T) {
	suite.Run(t, new(ingredientsTest))
}

func (i *ingredientsTest) Test_CombineRecipesToIngredients_Nil() {
	ings := CombineRecipesToIngredients(nil)
	i.Equal([]Ingredient{}, ings)
}

func (i *ingredientsTest) Test_CombineRecipesToIngredients_SameIng_Combine() {
	r1 := Recipe{
		Ings: []Ingredient{
			Ingredient{
				UnitSize:       "1",
				UnitType:       "tsbp",
				IngredientName: "Olive Oil",
			},
		},
	}
	r2 := Recipe{
		Ings: []Ingredient{
			Ingredient{
				UnitSize:       "3",
				UnitType:       "tsbp",
				IngredientName: "Olive Oil",
			},
		},
	}
	expected := []Ingredient{{
		UnitSize:       "4.00",
		UnitType:       "tsbp",
		IngredientName: "Olive Oil",
	}}
	ings := CombineRecipesToIngredients([]Recipe{r1, r2})
	i.Equal(expected, ings)
}
func (i *ingredientsTest) Test_CombineRecipesToIngredients_DiffIng_Combine() {
	r1 := Recipe{
		Ings: []Ingredient{
			Ingredient{
				UnitSize:       "1",
				UnitType:       "large",
				IngredientName: "onion",
			},
		},
	}
	r2 := Recipe{
		Ings: []Ingredient{
			Ingredient{
				UnitSize:       "3",
				UnitType:       "large",
				IngredientName: "onion",
			},
			Ingredient{
				UnitSize:       "1",
				UnitType:       "tsp",
				IngredientName: "oil",
			},
		},
	}
	expected := []Ingredient{{
		UnitSize:       "4.00",
		UnitType:       "large",
		IngredientName: "onion",
	}, {
		UnitSize:       "1",
		UnitType:       "tsp",
		IngredientName: "oil",
	},
	}
	ings := CombineRecipesToIngredients([]Recipe{r1, r2})
	i.Equal(expected, ings)
}
