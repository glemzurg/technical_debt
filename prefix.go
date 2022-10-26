package technical_debt

import (
	"path"
	"strings"
)

// LongestFilenamePrefix gets the longest prefix shared by the filenames.
func LongestFilenamePrefix(files map[string]CodeFile) (prefix string) {

	// Get the filenames as strings.
	var filenames []string
	for filename := range files {
		filenames = append(filenames, filename)
	}

	return LongestCommonPrefix(filenames)
}

// LongestCommonPrefix finds the longest common prefix from a list of strings.
func LongestCommonPrefix(values []string) (prefix string) {

	// If there are one or no values, assume no prefix.
	if len(values) <= 1 {
		return ""
	}

	// Find the prefix.
	prefix, _ = path.Split(values[0]) // Start with the entire first value.
	prefix = strings.TrimSuffix(prefix, "/")
	for _, value := range values {
		valuePrefix, _ := path.Split(value)
		// Is this prefix a complete subset of the value prefix?
		if !strings.HasPrefix(valuePrefix, prefix) {
			// Keep trimming prefix until we have found a subset or have no prefix.
			for prefix != "" {
				prefix, _ = path.Split(prefix)
				prefix = strings.TrimSuffix(prefix, "/")
				if strings.HasPrefix(valuePrefix, prefix) {
					// This prefix works. Continue.
					break
				}
			}
			// If no more prefix no need to coniue.
			if prefix == "" {
				break
			}
		}
	}

	return prefix
}
