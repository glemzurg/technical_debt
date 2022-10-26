package technical_debt

import (
	"crypto/md5"
	"fmt"
	"io"
	"sort"
	"strings"
)

// CyclicalGroup is a distince set of files with the same dependencies.
type CyclicalGroup struct {
	FileCount         int
	VisibilityFanIn   int
	VisibilityFanOut  int
	CyclicFingerprint string     // The md5 fingerprint of the dependencies.
	Files             []CodeFile // The sorted code files in this group.
}

// AddCyclicFingerPrints computes a fingerprint for each dependency set.
func AddCyclicFingerPrints(codeFiles map[string]CodeFile) {
	for filename := range codeFiles {
		codeFile := codeFiles[filename]

		// Start with filenames.
		var dependsOnFilenames []string
		for dependsOnFilename := range codeFile.DependsOn {
			dependsOnFilenames = append(dependsOnFilenames, dependsOnFilename)
		}
		sort.Strings(dependsOnFilenames)

		// Make an md5 of the information.
		h := md5.New()
		io.WriteString(h, strings.Join(dependsOnFilenames, ";"))
		md5 := fmt.Sprintf("%x", h.Sum(nil))

		// Remember it.
		codeFile.CyclicFingerprint = md5
		codeFiles[filename] = codeFile
	}
}

// CreateCyclicalGroups gathers code files into their cyclical groups.
func CreateCyclicalGroups(fileLookup map[string]CodeFile) (groups []CyclicalGroup, coreCount, fileCount int) {

	// Sort code files.
	fileCount = len(fileLookup)
	var codeFiles []CodeFile
	for _, codeFile := range fileLookup {
		codeFiles = append(codeFiles, codeFile)
	}
	sort.Sort(byCyclicalPoperties(codeFiles))

	// Gather into cyclical groups.
	var currentGroup CyclicalGroup
	for i, codeFile := range codeFiles {

		// Add indices for future template display.
		codeFile.Index = i
		if codeFile.CyclicFingerprint != currentGroup.CyclicFingerprint {
			currentGroup = CyclicalGroup{
				VisibilityFanIn:   codeFile.VisibilityFanIn,
				VisibilityFanOut:  codeFile.VisibilityFanOut,
				CyclicFingerprint: codeFile.CyclicFingerprint,
				Files:             []CodeFile{codeFile},
			}
		} else {
			currentGroup.Files = append(currentGroup.Files, codeFile)
		}

		// If we are ar the end of the list, or we're changing groups, add this group to our results.
		if i == len(codeFiles)-1 || codeFiles[i+1].CyclicFingerprint != codeFile.CyclicFingerprint {
			currentGroup.FileCount = len(currentGroup.Files)
			groups = append(groups, currentGroup)
			if currentGroup.FileCount > coreCount {
				coreCount = currentGroup.FileCount
			}
		}
	}

	// Sort the groups.
	sort.Sort(byCyclicalGroupPoperties(groups))

	return groups, coreCount, fileCount
}

// byCyclicalPoperties implements sort.Interface to sort for finding cyclical groups.
// Example: sort.Sort(byCyclicalPoperties(codeFiles))
type byCyclicalPoperties []CodeFile

func (a byCyclicalPoperties) Len() int      { return len(a) }
func (a byCyclicalPoperties) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byCyclicalPoperties) Less(i, j int) bool {
	if a[i].VisibilityFanIn == a[j].VisibilityFanIn {
		if a[i].VisibilityFanOut == a[j].VisibilityFanOut {
			if a[i].CyclicFingerprint == a[j].CyclicFingerprint {
				return a[i].Name < a[j].Name
			}
			return a[i].CyclicFingerprint < a[j].CyclicFingerprint
		}
		return a[i].VisibilityFanOut < a[j].VisibilityFanOut
	}
	return a[i].VisibilityFanIn > a[j].VisibilityFanIn
}

// byCyclicalGroupPoperties implements sort.Interface to sort for ordering cyclical groups.
// Example: sort.Sort(byCyclicalGroupPoperties(groups))
type byCyclicalGroupPoperties []CyclicalGroup

func (a byCyclicalGroupPoperties) Len() int      { return len(a) }
func (a byCyclicalGroupPoperties) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byCyclicalGroupPoperties) Less(i, j int) bool {
	if a[i].VisibilityFanIn == a[j].VisibilityFanIn {
		if a[i].VisibilityFanOut == a[j].VisibilityFanOut {
			if a[i].FileCount == a[j].FileCount {
				return a[i].CyclicFingerprint < a[j].CyclicFingerprint
			}
			return a[i].FileCount > a[j].FileCount
		}
		return a[i].VisibilityFanOut < a[j].VisibilityFanOut
	}
	return a[i].VisibilityFanIn > a[j].VisibilityFanIn
}
