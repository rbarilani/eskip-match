package cli

import (
	"testing"
)

func TestLoader(t *testing.T) {
	tests := []struct {
		name             string
		defaultFile      string
		file             string
		customFiltersLen int
	}{
		{
			name:             "Load file",
			file:             "./testdata/config.yml",
			customFiltersLen: 2,
		},
		{
			name:             "Load default file",
			defaultFile:      configDefaultFile,
			customFiltersLen: 2,
		},
		{
			name:             "Load not existent default file",
			defaultFile:      "blue.yaml",
			customFiltersLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := newConfigLoader(tt.defaultFile)
			config := loader.Load(tt.file)
			if len(config.CustomFilters) != tt.customFiltersLen {
				t.Errorf("expect CustomFilters to have been populated by %d items but got %d", tt.customFiltersLen, len(config.CustomFilters))
			}
		})
	}
}
