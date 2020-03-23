package AZURERMR001

import (
	"testing"

	"github.com/bflad/tfproviderlint/helper/analysisfixtest"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZURERMR001(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "a")
}

func TestAZURERMR001Fixes(t *testing.T) {
	testdata := analysistest.TestData()
	analysisfixtest.Run(t, testdata, Analyzer, "a")
}
