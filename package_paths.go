package technical_debt

import (
	"os"
	"path/filepath"
)

// PackagePaths gets all the paths that belong to this go path.
func PackagePaths(gopath string, paths []string) (packagePaths []string, err error) {

	// We're building a list of package paths.
	var packagePathSet map[string]bool = map[string]bool{} // A map as a set.

	// The tree walk function.
	var walkFunc filepath.WalkFunc = func(path string, info os.FileInfo, err2 error) (err3 error) {
		if err2 != nil {
			return err2
		}
		// Is this a go file?
		if filepath.Ext(path) == ".go" {
			// Add this package path.
			var fullPackagePath string = filepath.Dir(path)
			// This is a full path. Chop off the source root path.
			var packagePath string
			if packagePath, err3 = filepath.Rel(gopath, fullPackagePath); err3 != nil {
				return err3
			}
			packagePathSet[packagePath] = true
		}
		return nil
	}

	// Gather every folder that has at least a single .go file.
	for _, path := range paths {
		if err = filepath.Walk(gopath+"/"+path, walkFunc); err != nil {
			return nil, Error(err)
		}
	}

	// Convert the set ot a list.
	for path := range packagePathSet {
		packagePaths = append(packagePaths, path)
	}

	return packagePaths, nil
}
