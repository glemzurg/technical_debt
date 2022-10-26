package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/glemzurg/technical_debt"
)

func main() {
	var err error

	// Example call: go/bin/technical_debt -config /path/to/technical_debt/root/config/config.json

	var configFilename string
	flag.StringVar(&configFilename, "config", "", "the config for this technical debt")
	flag.Parse()

	if configFilename == "" {
		panic(`-config required`)
	}

	config, err := technical_debt.LoadConfig(configFilename)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("\n\nconfig: \n%+v\n\n", config)

	// Get all the paths we are graphing.
	packagePaths, err := technical_debt.PackagePaths(config.Gopath, config.Paths)
	if err != nil {
		panic(err.Error())
	}

	// Process each package in out project.
	folders, err := technical_debt.ProcessPackage(config.Gopath, packagePaths, config.Paths, config.IncludeTests)
	if err != nil {
		panic(err.Error())
	}

	// Create the codeFiles.
	codeFiles := technical_debt.CreateCodeFiles(folders)

	// Drive the codeFiles deep into the data structure.
	technical_debt.CalculateDeepCodeFiles(codeFiles)

	// Calculate important numbers.
	propogationCost := technical_debt.CalculateMetrics(codeFiles)

	// Add cyclic fingerprints.
	technical_debt.AddCyclicFingerPrints(codeFiles)

	// What is the longest prefix shared by filenames?
	prefix := technical_debt.LongestFilenamePrefix(codeFiles)

	// Compose into cyclical groups.
	groups, coreCount, fileCount := technical_debt.CreateCyclicalGroups(codeFiles)

	// Find the fan in and fan out for the view we want.
	var visibilityFanIn, visibilityFanOut int
	if config.View == technical_debt.VIEW_MEDIAN {
		visibilityFanIn, visibilityFanOut = technical_debt.FindMedianVisibilityFanInOut(codeFiles)
	} else {
		visibilityFanIn, visibilityFanOut = technical_debt.FindCorePeripheryVisibilityFanInOut(groups, coreCount)
	}

	// Partition for display.
	partitions := technical_debt.CreateViewPartions(groups, visibilityFanIn, visibilityFanOut)

	// Define some function for our template.
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"isDependency": func(file, potential technical_debt.CodeFile) bool {
			return file.DependsOn[potential.Name]
		},
		"trimPrefix": func(filename string) string {
			return strings.TrimPrefix(filename, prefix+"/")
		},
	}

	// Load the template.
	t := template.Must(template.New("grid.template").Funcs(funcMap).ParseFiles(config.RootPath + "/template/grid.template"))

	// Get the output bytes.
	var outputBuffer bytes.Buffer // A Buffer needs no initialization.

	// Generate text.
	err = t.Execute(&outputBuffer, struct {
		FileCount  int
		Partitions []technical_debt.Partition
	}{
		FileCount:  fileCount,
		Partitions: partitions,
	})
	if err != nil {
		panic(err.Error())
	}

	// Write the text to a file.
	if err = ioutil.WriteFile(config.RootPath+"/output/grid.svg", outputBuffer.Bytes(), os.ModePerm); err != nil {
		panic(err.Error())
	}

	fmt.Println("propogation cost:", propogationCost)
	fmt.Printf("core size: %d / %d == %.2f\n\n", coreCount, fileCount, float64(coreCount)/float64(fileCount))
}
