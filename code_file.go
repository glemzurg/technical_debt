package technical_debt

import (
// "fmt"
)

// CodeFile is a single file in the dependency graph.
type CodeFile struct {
	Name              string          // The file name with path.
	DependsOn         map[string]bool // The files this file depends on. Set represented as a map.
	DependedOnBy      map[string]bool // The files this file depends on. Set represented as a map.
	VisibilityFanIn   int
	VisibilityFanOut  int
	CyclicFingerprint string // The md5 fingerprint of the dependencies.
	Index             int    // The position in in the whole display this code file is (starting at zero).
}

// CreateCodeFiles creates the dependency map for all the files.
func CreateCodeFiles(folders []packageFolder) (codeFiles map[string]CodeFile) {

	// Create lookup of which declarations are in which files.
	declarationLookup := map[string]map[string]string{}
	for _, folder := range folders {
		declarationLookup[folder.importPath] = map[string]string{}
		for _, file := range folder.files {
			for _, declaration := range file.declarations {
				declarationLookup[folder.importPath][declaration.name] = folder.importPath + "/" + file.name
			}
		}
	}

	// fmt.Println("===========")
	// for importPath, declarations := range declarationLookup {
	// 	fmt.Printf(importPath + "\n")
	// 	for declaration, filename := range declarations {
	// 		fmt.Printf("\t" + declaration + " -> " + filename + "\n")
	// 	}
	// }

	// Go through folders and compose the file dependencies.
	codeFiles = map[string]CodeFile{}
	for _, folder := range folders {
		for _, file := range folder.files {
			codeFile := CodeFile{
				Name:         folder.importPath + "/" + file.name,
				DependsOn:    map[string]bool{},
				DependedOnBy: map[string]bool{},
			}
			for _, unresolved := range file.unresolved {
				if filename, found := declarationLookup[unresolved.packageName][unresolved.name]; found {
					codeFile.DependsOn[filename] = true
				}
			}
			codeFiles[codeFile.Name] = codeFile
		}
	}

	// fmt.Println("===========")
	// for _, codeFile := range codeFiles {
	// 	fmt.Printf(codeFile.Name + "\n")
	// 	for filename := range codeFile.DependsOn {
	// 		fmt.Printf("\t" + filename + "\n")
	// 	}
	// }

	// Turn the folders into a raw list of files and their dependencies.

	// Examine each

	// fmt.Println("===========")
	// fmt.Println(folder.String())
	// fmt.Println("===========")

	// Create files.
	//
	//fmt.Println(codeFiles)

	return codeFiles
}
