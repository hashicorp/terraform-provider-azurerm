package AZURERMR001

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZURERMR001(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "a")
}
