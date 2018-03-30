package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	config := Load("./config_test.yml")
	if len(config.CustomFilters) != 2 {
		t.Error("expect CustomFilters to have been populated by two items")
	}
}

func TestDefault(t *testing.T) {
	config := Default()
	if len(config.CustomFilters) != 2 {
		t.Error("expect CustomFilters to have been populated from default file by two items")
	}
}
