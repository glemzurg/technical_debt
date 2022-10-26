package technical_debt

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type PrefixSuite struct{}

var _ = Suite(&PrefixSuite{})

// Add the tests.

func (s *PrefixSuite) Test_LongestCommonPrefix(c *C) {
	tests := []struct {
		values []string
		prefix string
	}{
		{
			values: nil,
			prefix: "",
		},
		{
			values: []string{},
			prefix: "",
		},
		{
			values: []string{"path"},
			prefix: "",
		},
		{
			values: []string{"path", "path"},
			prefix: "",
		},
		{
			values: []string{"path/a", "path"},
			prefix: "",
		},
		{
			values: []string{"path", "path/a"},
			prefix: "",
		},
		{
			values: []string{"path/thing", "path/other", "path/swing"},
			prefix: "path",
		},
		{
			values: []string{"path/a/thing", "path/a/other", "path/a/swing"},
			prefix: "path/a",
		},
		{
			values: []string{"path/a/thing", "other/a/other", "path/a/swing"},
			prefix: "",
		},
	}
	for i, test := range tests {
		comment := Commentf("Case %v: %v", i, test)
		c.Check(LongestCommonPrefix(test.values), Equals, test.prefix, comment)
	}
}
