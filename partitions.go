package technical_debt

import (
	"sort"
)

const (
	VIEW_CORE_PERIPHERY = "core-periphery"
	VIEW_MEDIAN         = "median"
)

type Partition struct {
	LowestIndex  int
	HighestIndex int
	FileCount    int
	Groups       []CyclicalGroup
}

// FindCorePeripheryVisibilityFanInOut finds the visibility fan in/fan out for core perifery view.
func FindCorePeripheryVisibilityFanInOut(groups []CyclicalGroup, coreCount int) (visibilityFanIn, visibilityFanOut int) {

	// Find the Visibility fan in and fan out for the core group.
	for _, group := range groups {
		if group.FileCount == coreCount {
			visibilityFanIn = group.VisibilityFanIn
			visibilityFanOut = group.VisibilityFanOut
			break
		}
	}

	return visibilityFanIn, visibilityFanOut
}

// FindMedianVisibilityFanInOut finds the visibility fan in/fan out for core perifery view.
func FindMedianVisibilityFanInOut(fileLookup map[string]CodeFile) (visibilityFanIn, visibilityFanOut int) {

	// Sort code files.
	fileCount := len(fileLookup)
	var codeFiles []CodeFile
	for _, codeFile := range fileLookup {
		codeFiles = append(codeFiles, codeFile)
	}
	sort.Sort(byCyclicalPoperties(codeFiles))

	// Is there an odd or even numbered file count?
	if fileCount%2 == 0 {
		// An even number, we need to find the middle two files. Examples:
		//   2 files ---> indexes *0 *1
		//   4 files ---> indexes 0 *1 *2 3
		//   6 files ---> indexes 0 1 *2 *3 4 5
		//   8 files ---> indexes 0 1 2 *3 *4 5 6 7
		// Or:
		//   higher index = len(files) / 2
		//   lower index is one less
		i := len(codeFiles) / 2

		// Integer division will drop remainder but it will be fine for our graphing of hundreds of files.
		visibilityFanIn = (codeFiles[i].VisibilityFanIn + codeFiles[i-1].VisibilityFanIn) / 2
		visibilityFanOut = (codeFiles[i].VisibilityFanOut + codeFiles[i-1].VisibilityFanOut) / 2
	} else {
		// An odd number we need to get the middle file. Examples:
		//   1 file ---> indexes 0
		//   3 file ---> indexes 0 *1 2
		//   5 file ---> indexes 0 1 *2 3 4
		//   7 file ---> indexes 0 1 2 *3 4 5 6
		// Or:
		//   index = len(files) - 1 / 2
		i := (len(codeFiles) - 1) / 2
		visibilityFanIn = codeFiles[i].VisibilityFanIn
		visibilityFanOut = codeFiles[i].VisibilityFanOut
	}

	return visibilityFanIn, visibilityFanOut
}

// CreateViewPartions creates partitions for the display of the code dependencies.
func CreateViewPartions(groups []CyclicalGroup, visibilityFanIn, visibilityFanOut int) (partitions []Partition) {

	// All the groups are in the proper order, but we need to partition them.
	// a) "Shared" elements have VFI ≥ VFIC and VFO < VFOC.
	// b) "Peripheral" elements have VFI < VFIC and VFO < VFOC.
	// c) "Control" elements have VFI < VFIC and VFO ≥ VFOC.
	var core, shared, periphery, control Partition
	for _, group := range groups {
		if group.VisibilityFanIn >= visibilityFanIn && group.VisibilityFanOut < visibilityFanOut {
			// This is s shared group.
			shared.Groups = append(shared.Groups, group)
		} else if group.VisibilityFanIn < visibilityFanIn && group.VisibilityFanOut < visibilityFanOut {
			// This is a periphery group.
			periphery.Groups = append(periphery.Groups, group)
		} else if group.VisibilityFanIn < visibilityFanIn && group.VisibilityFanOut >= visibilityFanOut {
			// This is a control group.
			control.Groups = append(control.Groups, group)
		} else {
			// This is a core group.
			core.Groups = append(core.Groups, group)
		}
	}

	// Put them in order.
	partitions = []Partition{
		shared,
		core,
		periphery,
		control,
	}

	// All the files need to be update to have proper indexes.
	var i int
	for a, partition := range partitions {
		for b, group := range partition.Groups {
			for c := range group.Files {
				partitions[a].Groups[b].Files[c].Index = i
				partitions[a].FileCount++
				if a > 0 {
					if partitions[a].LowestIndex == 0 || i < partitions[a].LowestIndex {
						partitions[a].LowestIndex = i
					}
				}
				if i > partitions[a].HighestIndex {
					partitions[a].HighestIndex = i
				}
				i++
			}
		}
	}

	return partitions
}
