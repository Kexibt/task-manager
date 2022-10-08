package config

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

// TestConfig tests config
func TestConfig(t *testing.T) {
	cfgCopy := GetConfig()
	cfgExpected := Config{}

	b, err := os.ReadFile(cfgFilename)
	if err != nil {
		t.Error(err)
	}

	err = yaml.Unmarshal(b, &cfgExpected)
	if err != nil {
		t.Error(err)
	}

	if cfgExpected != cfgCopy {
		t.Errorf("Expected:\n\t%+v, but got:\n\t%+v", cfgExpected, cfgCopy)
	}
}
