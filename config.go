package technical_debt

import (
	"encoding/json"
	"io/ioutil"
)

// Config is the information we need to run.
type Config struct {
	Gopath       string
	RootPath     string
	Paths        []string
	View         string
	IncludeTests bool
}

// LoadConfig loads a json config.
func LoadConfig(filename string) (config Config, err error) {

	// Load the config.
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, Error(err)
	}

	// Parse the data.
	if err = json.Unmarshal(bytes, &config); err != nil {
		return Config{}, Error(err)
	}

	// Verify the
	if err = config.validate(); err != nil {
		return Config{}, Error(err)
	}

	return config, nil
}

// validate confirms a well-formed config.
func (c Config) validate() (err error) {
	if c.Gopath == "" {
		return Errorf(`config requires Gopath`)
	}
	if c.RootPath == "" {
		return Errorf(`config requires RootPath`)
	}
	if len(c.Paths) == 0 {
		return Errorf(`config requires Paths`)
	}
	for _, path := range c.Paths {
		if path == "" {
			return Errorf(`config requires each Path to be non-blank`)
		}
	}
	if !(c.View == VIEW_CORE_PERIPHERY || c.View == VIEW_MEDIAN) {
		return Errorf(`config View must be either '%s' or '%s'`, VIEW_CORE_PERIPHERY, VIEW_MEDIAN)
	}

	return nil
}
