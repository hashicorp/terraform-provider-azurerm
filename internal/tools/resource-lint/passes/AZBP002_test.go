// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package passes_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZBP002(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, passes.AZBP002Analyzer, "testdata/src/azbp002")
}
