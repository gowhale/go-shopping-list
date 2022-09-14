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
	if err := execute(); err != nil {
		log.Fatalln(err)
	}
}

func execute() error {
	output, err := runGoTest()
	if err != nil {
		log.Println(output)
		return err
	}

	tl, err := covertOutputToCoverage(output)
	if err != nil {
		return err
	}

	return validateTestOutput(tl, output)
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
			coverageIndex := strings.Index(line, "coverage: ")
			if coverageIndex != -1 {
				lineFields := strings.Fields(line[coverageIndex:])
				pkgPercentStr := lineFields[1][:len(lineFields[1])-1]
				pkgPercentFloat, err := strconv.ParseFloat(pkgPercentStr, 64)
				if err != nil {
					return nil, err
				}
				testStruct = append(testStruct, testLine{pkgName: pkgName, coverage: pkgPercentFloat})
			} else {
				testStruct = append(testStruct, testLine{pkgName: pkgName, coverage: -1})
			}
		}
	}
	return testStruct, nil
}

func validateTestOutput(tl []testLine, o string) error {
	invalidOutputs := []string{}
	for _, line := range tl {
		switch {
		case line.coverage == -1:
			invalidOutputs = append(invalidOutputs, fmt.Sprintf("pkg=%s is missing tests", line.pkgName))
		case line.coverage < 80:
			invalidOutputs = append(invalidOutputs, fmt.Sprintf("pkg=%s cov=%f under the %f%% minimum line coverage", line.pkgName, line.coverage, 80.0))
		}
	}
	for _, invalid := range invalidOutputs {
		log.Println(invalid)
	}
	if len(invalidOutputs) == 0 {
		return nil
	}
	log.Println(o)
	return fmt.Errorf("the following pkgs are not valid: %+v", invalidOutputs)
}
