package technical_debt

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck

	"sort"
)

// Create a suite.
type CyclicalGroupSuite struct{}

var _ = Suite(&CyclicalGroupSuite{})

// Add the tests.

func (s *CyclicalGroupSuite) Test_SortCyclicalGroups(c *C) {

	groups := []CyclicalGroup{
		CyclicalGroup{
			FileCount:         10,
			VisibilityFanIn:   9,
			VisibilityFanOut:  11,
			CyclicFingerprint: "CyclicC",
		},
		CyclicalGroup{
			FileCount:         1,
			VisibilityFanIn:   8,
			VisibilityFanOut:  10,
			CyclicFingerprint: "CyclicA",
		},
		CyclicalGroup{
			FileCount:         9,
			VisibilityFanIn:   9,
			VisibilityFanOut:  11,
			CyclicFingerprint: "CyclicA",
		},
		CyclicalGroup{
			FileCount:         1,
			VisibilityFanIn:   9,
			VisibilityFanOut:  10,
			CyclicFingerprint: "CyclicA",
		},
		CyclicalGroup{
			FileCount:         1,
			VisibilityFanIn:   10,
			VisibilityFanOut:  10,
			CyclicFingerprint: "CyclicA",
		},
		CyclicalGroup{
			FileCount:         9,
			VisibilityFanIn:   9,
			VisibilityFanOut:  11,
			CyclicFingerprint: "CyclicB",
		},
	}

	sort.Sort(byCyclicalGroupPoperties(groups))

	c.Assert(groups, DeepEquals, []CyclicalGroup{
		CyclicalGroup{
			FileCount:         1,
			VisibilityFanIn:   10, // First sort by fan in descending..
			VisibilityFanOut:  10,
			CyclicFingerprint: "CyclicA",
		},
		CyclicalGroup{
			FileCount:         1,
			VisibilityFanIn:   9,
			VisibilityFanOut:  10, // Next by fan out ascending.
			CyclicFingerprint: "CyclicA",
		},
		CyclicalGroup{
			FileCount:         10, // Next by file count descending.
			VisibilityFanIn:   9,
			VisibilityFanOut:  11,
			CyclicFingerprint: "CyclicC",
		},
		CyclicalGroup{
			FileCount:         9,
			VisibilityFanIn:   9,
			VisibilityFanOut:  11,
			CyclicFingerprint: "CyclicA", // Next fingerprint ascending.
		},
		CyclicalGroup{
			FileCount:         9,
			VisibilityFanIn:   9,
			VisibilityFanOut:  11,
			CyclicFingerprint: "CyclicB",
		},
		CyclicalGroup{
			FileCount:         1,
			VisibilityFanIn:   8,
			VisibilityFanOut:  10,
			CyclicFingerprint: "CyclicA",
		},
	})
}

func (s *CyclicalGroupSuite) Test_SortCodeFiles(c *C) {

	codeFiles := []CodeFile{
		CodeFile{
			Name:              "NameD",
			VisibilityFanIn:   8,
			VisibilityFanOut:  10,
			CyclicFingerprint: "CyclicA",
		},
		CodeFile{
			Name:              "NameA",
			VisibilityFanIn:   9,
			VisibilityFanOut:  11,
			CyclicFingerprint: "CyclicA",
		},
		CodeFile{
			Name:              "NameM",
			VisibilityFanIn:   9,
			VisibilityFanOut:  10,
			CyclicFingerprint: "CyclicA",
		},
		CodeFile{
			Name:              "NameB",
			VisibilityFanIn:   9,
			VisibilityFanOut:  11,
			CyclicFingerprint: "CyclicA",
		},
		CodeFile{
			Name:              "NameX",
			VisibilityFanIn:   9,
			VisibilityFanOut:  11,
			CyclicFingerprint: "CyclicB",
		},
		CodeFile{
			Name:              "NameZ",
			VisibilityFanIn:   10,
			VisibilityFanOut:  10,
			CyclicFingerprint: "CyclicA",
		},
	}

	sort.Sort(byCyclicalPoperties(codeFiles))

	c.Assert(codeFiles, DeepEquals, []CodeFile{
		CodeFile{
			Name:              "NameZ",
			VisibilityFanIn:   10, // First sort by fan in descending..
			VisibilityFanOut:  10,
			CyclicFingerprint: "CyclicA",
		},
		CodeFile{
			Name:              "NameM",
			VisibilityFanIn:   9,
			VisibilityFanOut:  10, // Next by fan out ascending.
			CyclicFingerprint: "CyclicA",
		},
		CodeFile{
			Name:              "NameA",
			VisibilityFanIn:   9,
			VisibilityFanOut:  11,
			CyclicFingerprint: "CyclicA", // Next fingerprint ascending.
		},
		CodeFile{
			Name:              "NameB", // Next by name ascending.
			VisibilityFanIn:   9,
			VisibilityFanOut:  11,
			CyclicFingerprint: "CyclicA",
		},
		CodeFile{
			Name:              "NameX",
			VisibilityFanIn:   9,
			VisibilityFanOut:  11,
			CyclicFingerprint: "CyclicB",
		},
		CodeFile{
			Name:              "NameD",
			VisibilityFanIn:   8,
			VisibilityFanOut:  10,
			CyclicFingerprint: "CyclicA",
		},
	})
}
