package AZRMS001_test

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurermproviderlint/passes/AZRMS001"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZRMS001(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, AZRMS001.Analyzer, "a")
}

func TestAZRMS001Fixes(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, AZRMS001.Analyzer, "a")
}
