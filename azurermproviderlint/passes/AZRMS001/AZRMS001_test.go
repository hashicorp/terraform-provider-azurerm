package AZRMS001_test

import (
	"log"
	"os"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurermproviderlint/helper/testsetup"
	"github.com/terraform-providers/terraform-provider-azurerm/azurermproviderlint/passes/AZRMS001"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMain(m *testing.M) {
	code, err := testsetup.TestSetup(m)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func TestAZRMS001(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, AZRMS001.Analyzer, "a")
}

func TestAZRMS001Fixes(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, AZRMS001.Analyzer, "a")
}
