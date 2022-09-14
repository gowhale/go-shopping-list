package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var execCommand = exec.Command
var excludedPkgs = map[string]bool{
	"go-shopping-list":                  true,
	"go-shopping-list/cmd/pkg-cover":    true,
	"go-shopping-list/cmd/authenticate": true,
}

func main() {
	output, err := runGoTest()
	if err != nil {
		log.Println(output)
		log.Fatalln(err)
	}

	tl, err := covertOutputToCoverage(output)
	if err != nil {
		log.Fatalln(err)
	}

	if err := validateTestOutput(tl); err != nil {
		log.Fatalln(err)
	}
}

func runGoTest() (string, error) {
	cmd := execCommand("go", "test", "./...", "--cover")
	output, err := cmd.CombinedOutput()
	termOutput := string(output)
	return termOutput, err
}

type testLine struct {
	pkgName  string
	coverage float64
}

func covertOutputToCoverage(termOutput string) ([]testLine, error) {
	testStruct := []testLine{}
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
					return nil, err
				}
				testStruct = append(testStruct, testLine{pkgName: pkgName, coverage: s})
			} else {
				testStruct = append(testStruct, testLine{pkgName: pkgName, coverage: -1})
			}
		}
	}
	return testStruct, nil
}

func validateTestOutput(tl []testLine) error {
	invalidOutputs := []string{}
	for _, test := range tl {
		switch {
		case test.coverage == -1:
			invalidOutputs = append(invalidOutputs, fmt.Sprintf("pkg=%s is missing tests", test.pkgName))
		case test.coverage < 80:
			invalidOutputs = append(invalidOutputs, fmt.Sprintf("pkg=%s cov=%f under the %f%% minimum line coverage", test.pkgName, test.coverage, 80.0))
		}
	}
	for _, invalid := range invalidOutputs {
		log.Println(invalid)
	}
	if len(invalidOutputs) == 0 {
		return nil
	}
	return fmt.Errorf("the following pkgs are not valid: %+v", invalidOutputs)
}

func pkgCover(termOutput string) error {
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
