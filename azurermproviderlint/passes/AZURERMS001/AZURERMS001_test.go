package AZURERMS001_test

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurermproviderlint/passes/AZURERMS001"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZURERMS001(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, AZURERMS001.Analyzer, "a")
}

func TestAZURERMS001Fixes(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, AZURERMS001.Analyzer, "a")
}
