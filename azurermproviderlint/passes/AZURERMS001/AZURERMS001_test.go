package AZURERMS001_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurermproviderlint/passes/AZURERMS001"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMain(m *testing.M) {
	// Setup
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(pwd, "testdata", "src", "a", "vendor", "github.com", "terraform-providers", "terraform-provider-azurerm")
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatal(err)
	}
	symlink := filepath.Join(path, "azurerm")
	target := pwd
	for i := 0; i < 3; i++ {
		target = filepath.Dir(target)
	}
	target = filepath.Join(target, "azurerm")
	if err := os.Symlink(target, symlink); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	// Teardown
	os.RemoveAll(path)

	os.Exit(code)
}

func TestAZURERMS001(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, AZURERMS001.Analyzer, "a")
}

func TestAZURERMS001Fixes(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, AZURERMS001.Analyzer, "a")
}
