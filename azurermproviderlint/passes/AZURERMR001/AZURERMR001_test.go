package AZURERMR001_test

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurermproviderlint/passes/AZURERMR001"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZURERMR001(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, AZURERMR001.Analyzer, "a")
}

func TestAZURERMR001Fixes(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, AZURERMR001.Analyzer, "a")
}
