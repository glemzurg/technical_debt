package technical_debt

// CalculateMetrics calcualtes important nubmer for the algorithm.
func CalculateMetrics(codeFiles map[string]CodeFile) (propogationCost float64) {

	for filename := range codeFiles {
		codeFile := codeFiles[filename]
		codeFile.VisibilityFanIn = len(codeFile.DependedOnBy)
		codeFile.VisibilityFanOut = len(codeFile.DependsOn)
		codeFiles[filename] = codeFile
	}

	// What is the total fan in and fan out?
	var totalFanIn, totalFanOut int
	for _, codeFile := range codeFiles {
		totalFanIn += codeFile.VisibilityFanIn
		totalFanOut += codeFile.VisibilityFanOut
	}

	// Sanity check.
	if totalFanIn != totalFanOut {
		panic(`total fan in and fan out should be identical`)
	}

	propogationCost = float64(totalFanIn) / float64(len(codeFiles)*len(codeFiles))
	return propogationCost
}
