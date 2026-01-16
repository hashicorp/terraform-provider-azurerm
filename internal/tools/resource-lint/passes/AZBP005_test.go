package passes

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZBP005(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, AZBP005Analyzer, "testdata/src/azbp005")
}
