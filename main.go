package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"fyne.io/fyne/v2/widget"
)

var data = []string{"a", "string", "list"}

type Recipe struct {
	Name string        `json:"recipe_name"`
	Ings []Ingredients `json:"ingredients"`
	Meth []string      `json:"method"`
}

type Ingredients struct {
	Unit_size       string `json:"unit_size"`
	Unit_type       string `json:"unit_type"`
	Ingredient_name string `json:"ingredient_name"`
}

func (i *Ingredients) String() string {
	return fmt.Sprintf("%s %s %s", i.Unit_size, i.Unit_type, i.Ingredient_name)
}

func main() {

	allRecipes := []Recipe{}

	// Get name for all recipe files
	files, err := ioutil.ReadDir("recipes/")
	if err != nil {
		log.Fatal(err)
	}

	// Process every file and put into Recipe strucr
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		file, _ := ioutil.ReadFile(fmt.Sprintf("recipes/%s", file.Name()))
		recipe := Recipe{}
		_ = json.Unmarshal([]byte(file), &recipe)

		allRecipes = append(allRecipes, recipe)
	}
	log.Printf("amount of recipes=%d", len(files))

	myApp := app.New()
	myWindow := myApp.NewWindow("List Widget")

	// Progress bar for adding ings
	p := widget.NewProgressBar()

	// Recipe list with all recipes
	recipeList := widget.NewList(
		func() int {
			return len(allRecipes[1].Ings)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("template", func() {})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Button).SetText(allRecipes[i].Name)
			o.(*widget.Button).OnTapped = func() { addIngredientsToReminders(allRecipes[i], p) }
			test := widget.NewButton(fmt.Sprintf("%d", i), func() { log.Println(i) })
			o = test
		})

	// Create content grid
	grid := container.New(layout.NewGridLayout(1), recipeList)
	gridTop := container.New(layout.NewGridWrapLayout(fyne.NewSize(600, 50)), p)
	masterGrid := container.New(layout.NewGridLayoutWithColumns(1), gridTop, grid)

	// Set Window and execute
	myWindow.SetFixedSize(true)
	myWindow.Resize(fyne.Size{600, 1200})
	myWindow.SetContent(masterGrid)
	myWindow.ShowAndRun()

}

var execCommand = exec.Command

func addIngredientsToReminders(r Recipe, p *widget.ProgressBar) error {
	progress := 0.0
	for i, ing := range r.Ings {
		progress = float64(i) / float64(len(r.Ings))
		p.SetValue(progress)
		log.Printf("progress=%.2f adding ing='%s'", progress, ing.String())
		cmd := execCommand("automator", "-i", fmt.Sprintf(`"%s"`, ing.String()), "shopping.workflow")
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalln(err)
		}
	}
	progress = 1
	log.Printf("progress=%.2f", progress)
	p.SetValue(progress)
	return nil
}
