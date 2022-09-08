package recipe

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	tsbp      = "tsbp"
	tsp       = "tsp"
	oliveOil  = "Olive Oil"
	largeType = "large"
	oneUnit   = "1"
	onion     = "onion"
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
	ings, err := CombineRecipesToIngredients(nil)
	i.Equal([]Ingredient{}, ings)
	i.Nil(err)
}

func (i *ingredientsTest) Test_CombineRecipesToIngredients_SameIng_Combine() {
	r1 := Recipe{
		Ings: []Ingredient{
			Ingredient{
				UnitSize:       oneUnit,
				UnitType:       tsbp,
				IngredientName: oliveOil,
			},
		},
	}
	r2 := Recipe{
		Ings: []Ingredient{
			Ingredient{
				UnitSize:       "3",
				UnitType:       tsbp,
				IngredientName: oliveOil,
			},
		},
	}
	expected := []Ingredient{{
		UnitSize:       "4.00",
		UnitType:       tsbp,
		IngredientName: oliveOil,
	}}
	ings, err := CombineRecipesToIngredients([]Recipe{r1, r2})
	i.Equal(expected, ings)
	i.Nil(err)
}

func (i *ingredientsTest) Test_CombineRecipesToIngredients_DiffIng_Combine() {
	r1 := Recipe{
		Ings: []Ingredient{
			Ingredient{
				UnitSize:       oneUnit,
				UnitType:       largeType,
				IngredientName: onion,
			},
		},
	}
	r2 := Recipe{
		Ings: []Ingredient{
			Ingredient{
				UnitSize:       "3",
				UnitType:       largeType,
				IngredientName: onion,
			},
			Ingredient{
				UnitSize:       oneUnit,
				UnitType:       tsp,
				IngredientName: "oil",
			},
		},
	}
	expected := []Ingredient{{
		UnitSize:       "4.00",
		UnitType:       largeType,
		IngredientName: onion,
	}, {
		UnitSize:       oneUnit,
		UnitType:       tsp,
		IngredientName: "oil",
	},
	}
	ings, err := CombineRecipesToIngredients([]Recipe{r1, r2})
	i.Nil(err)
	i.Equal(expected, ings)
}

func (i *ingredientsTest) Test_CombineRecipesToIngredients_DiffIng_Combine_Error() {
	r1 := Recipe{
		Ings: []Ingredient{
			Ingredient{
				UnitSize:       oneUnit,
				UnitType:       largeType,
				IngredientName: onion,
			},
		},
	}
	r2 := Recipe{
		Ings: []Ingredient{
			Ingredient{
				UnitSize:       "EGG",
				UnitType:       largeType,
				IngredientName: onion,
			},
			Ingredient{
				UnitSize:       "EGG",
				UnitType:       tsp,
				IngredientName: "oil",
			},
		},
	}
	ings, err := CombineRecipesToIngredients([]Recipe{r1, r2})
	i.Nil(ings)
	i.EqualError(err, "strconv.ParseFloat: parsing \"EGG\": invalid syntax")
}
