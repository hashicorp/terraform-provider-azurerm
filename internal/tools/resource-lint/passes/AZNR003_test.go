// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package passes

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZNR003(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, AZNR003Analyzer, "testdata/src/aznr003")
}
