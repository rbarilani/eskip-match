package cli

import (
	"testing"
)

func TestLoader(t *testing.T) {
	scenarios := []struct {
		title               string
		defaultFile         string
		file                string
		expCustomFiltersLen int
	}{
		{
			title:               "Load file",
			file:                "./testdata/config.yml",
			expCustomFiltersLen: 2,
		},
		{
			title:               "Load default file",
			defaultFile:         configDefaultFile,
			expCustomFiltersLen: 2,
		},
		{
			title:               "Load not existent default file",
			defaultFile:         "blue.yaml",
			expCustomFiltersLen: 0,
		},
	}

	for _, s := range scenarios {
		t.Run(s.title, func(t *testing.T) {
			loader := newConfigLoader(s.defaultFile)
			config := loader.Load(s.file)
			if len(config.CustomFilters) != s.expCustomFiltersLen {
				t.Errorf("expect CustomFilters to have been populated by %d items but got %d", s.expCustomFiltersLen, len(config.CustomFilters))
			}
		})
	}
}
