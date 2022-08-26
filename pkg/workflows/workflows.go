package workflows

import (
	"fmt"
	"go-shopping-list/pkg/recipe"
	"log"
	"os/exec"
	"sync"

	"fyne.io/fyne/v2/widget"
)

var execCommand = exec.Command

func runReminder(p *widget.ProgressBar, l *widget.Label, currentIng recipe.Ingredients) {
	cmd := execCommand("automator", "-i", fmt.Sprintf(`"%s"`, currentIng.String()), "shopping.workflow")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Sprintf("error adding the following ingredient=%s err=%e", currentIng.String(), err)
	}
	l.SetText(fmt.Sprintf("Added Ingredient: %s", currentIng.String()))
	l.Refresh()
}

func AddIngredientsToReminders(r recipe.Recipe, p *widget.ProgressBar, l *widget.Label) error {
	l.SetText(fmt.Sprintf("Starting to add ingredients for Recipe: %s", r.Name))
	l.Refresh()
	f := recipe.FileInteractionImpl{}
	if err := f.IncrementPopularity(r.Name); err != nil {
		return err
	}
	progress := 0.0
	p.SetValue(progress)
	p.Refresh()
	ingAdded := []recipe.Ingredients{}
	var wg sync.WaitGroup
	for _, ing := range r.Ings {
		wg.Add(1)
		ing := ing
		go func() {
			runReminder(p, l, ing)
			defer func() {
				wg.Done()
				ingAdded = append(ingAdded, ing)
				progress = float64(len(ingAdded)) / float64(len(r.Ings))
				p.SetValue(progress)
				log.Printf("progress=%.2f adding ing='%s'", progress, ing.String())
			}()
		}()
	}
	wg.Wait()
	progress = 1
	log.Printf("progress=%.2f", progress)
	p.SetValue(progress)
	l.SetText("Finished. Select another recipe to add more.")
	l.Refresh()
	return nil
}
