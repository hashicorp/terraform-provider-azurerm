package passes_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZRE001(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, passes.AZRE001Analyzer, "testdata/src/azre001")
}
