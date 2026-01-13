package passes_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZNR002(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, passes.AZNR002Analyzer, "testdata/src/aznr002")
}
