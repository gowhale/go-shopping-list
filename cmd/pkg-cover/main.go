package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	log.Println(pkgCover())
}

var execCommand = exec.Command
var excludedPkgs = map[string]bool{
	"go-shopping-list":                  true,
	"go-shopping-list/cmd/pkg-cover":    true,
	"go-shopping-list/cmd/authenticate": true,
}

type testLine struct {
}

func pkgCover() error {
	cmd := execCommand("go", "test", "./...", "--cover")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	termOutput := string(output)
	lines := strings.Split(termOutput, "\n")
	for _, line := range lines[:len(lines)-1] {
		pkgName := strings.Fields(line)[1]
		if _, ok := excludedPkgs[pkgName]; !ok {
			coverageLine := strings.Index(line, "coverage: ")
			if coverageLine != -1 {
				words := strings.Fields(line[coverageLine:])
				percentageOfPkgCovered := words[1][:len(words[1])-1]
				s, err := strconv.ParseFloat(percentageOfPkgCovered, 64)
				if err != nil {
					return err
				}
				log.Printf("pkg=%s cov=%f", pkgName, s)
				if s < 80 {
					return fmt.Errorf("pkg=%s cov=%f does not meet %f pkg coverage", pkgName, s, 80.0)
				}
			} else {
				return fmt.Errorf("pkg=%s missing test files", pkgName)
			}
		}
	}
	return nil
}
