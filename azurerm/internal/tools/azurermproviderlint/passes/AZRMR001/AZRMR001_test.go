package AZRMR001_test

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tools/azurermproviderlint/passes/AZRMR001"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZRMR001(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, AZRMR001.Analyzer, "a")
}

func TestAZRMR001Fixes(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, AZRMR001.Analyzer, "a")
}
