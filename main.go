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

	recipeNames := []Recipe{}

	files, err := ioutil.ReadDir("recipes/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())

		file, _ := ioutil.ReadFile(fmt.Sprintf("recipes/%s", file.Name()))

		data := Recipe{}

		_ = json.Unmarshal([]byte(file), &data)

		fmt.Println(data.Name)
		recipeNames = append(recipeNames, data)
		for i, ing := range data.Ings {
			fmt.Println(i, ing)
		}
		for i, step := range data.Meth {
			fmt.Println(i, step)
		}

	}

	myApp := app.New()
	myWindow := myApp.NewWindow("List Widget")

	p := widget.NewProgressBar()

	list2 := widget.NewList(
		func() int {
			return len(recipeNames[1].Ings)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("template", func() {})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Button).SetText(recipeNames[i].Name)
			o.(*widget.Button).OnTapped = func() { addIngredientsToReminders(recipeNames[i], p) }
			test := widget.NewButton(fmt.Sprintf("%d", i), func() { log.Println(i) })
			o = test
		})

	// data := binding.BindStringList(
	// 	&recipeNames[1].Ings.Ingredient_name,
	// )

	// list := widget.NewListWithData(data,
	// 	func() fyne.CanvasObject {
	// 		return widget.NewLabel("template")
	// 	},
	// 	func(i binding.DataItem, o fyne.CanvasObject) {
	// 		o.(*widget.Label).Bind(i.(binding.String))
	// 	})
	// list := widget.NewList(
	// 	func() int {
	// 		return len(recipeNames)
	// 	},
	// 	func() fyne.CanvasObject {
	// 		return widget.NewButton("template", func() {
	// 			list2 = widget.NewList(func() int {
	// 				return len(recipeNames)
	// 			},
	// 				func() fyne.CanvasObject {
	// 					return widget.NewButton("template", func() { log.Println("go") })
	// 				},
	// 				func(i widget.ListItemID, o fyne.CanvasObject) {

	// 				})
	// 		})
	// 	},
	// 	func(i widget.ListItemID, o fyne.CanvasObject) {
	// 		o.(*widget.Button).SetText(recipeNames[i].Name)
	// 		o.(*widget.Button).OnTapped = func() { log.Println(fmt.Sprintf("%d epic", i)) }
	// 	})

	grid := container.New(layout.NewGridLayout(1), list2)
	gridTop := container.New(layout.NewGridWrapLayout(fyne.NewSize(600, 50)), p)
	p.Resize(fyne.Size{100, 100})

	masterGrid := container.New(layout.NewGridLayoutWithColumns(1), gridTop, grid)

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
