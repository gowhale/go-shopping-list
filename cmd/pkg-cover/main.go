package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const (
	minPercentCov          = 80.0
	coverageStringNotFound = -1
	firstItemIndex             = 1
)

var execCommand = exec.Command
var excludedPkgs = map[string]bool{
	"go-shopping-list":                  true,
	"go-shopping-list/cmd/pkg-cover":    true,
	"go-shopping-list/cmd/authenticate": true,
	"go-shopping-list/pkg/common":       true,
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
		if !strings.Contains(line, "go: downloading") {
			pkgName := strings.Fields(line)[firstItemIndex]
			if _, ok := excludedPkgs[pkgName]; !ok {
				coverageIndex := strings.Index(line, "coverage: ")
				if coverageIndex != coverageStringNotFound {
					lineFields := strings.Fields(line[coverageIndex:])
					pkgPercentStr := lineFields[firstItemIndex][:len(lineFields[firstItemIndex])-1]
					pkgPercentFloat, err := strconv.ParseFloat(pkgPercentStr, 64)
					if err != nil {
						return nil, err
					}
					testStruct = append(testStruct, testLine{pkgName: pkgName, coverage: pkgPercentFloat})
				} else {
					testStruct = append(testStruct, testLine{pkgName: pkgName, coverage: coverageStringNotFound})
				}
			}
		}
	}
	return testStruct, nil
}

func validateTestOutput(tl []testLine, o string) error {
	invalidOutputs := []string{}
	for _, line := range tl {
		switch {
		case line.coverage == coverageStringNotFound:
			invalidOutputs = append(invalidOutputs, fmt.Sprintf("pkg=%s is missing tests", line.pkgName))
		case line.coverage < minPercentCov:
			invalidOutputs = append(invalidOutputs, fmt.Sprintf("pkg=%s cov=%f under the %f%% minimum line coverage", line.pkgName, line.coverage, minPercentCov))
		}
	}
	if len(invalidOutputs) == 0 {
		return nil
	}
	log.Println(o)
	log.Println("###############################")
	log.Println("###############################")
	log.Println("invalid pkg's:")
	for i, invalid := range invalidOutputs {
		log.Printf("id=%d problem=%s", i, invalid)
	}
	return fmt.Errorf("the following pkgs are not valid: %+v", invalidOutputs)
}
