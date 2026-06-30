// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package passes

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/bflad/tfproviderlint/passes/commentignore"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/reporting"
	"golang.org/x/tools/go/analysis"
)

const AZBP005Doc = `check that Go source files have the correct licensing header

The AZBP005 analyzer reports cases where Go source files do not have the
required licensing header at the very beginning of the file.

Required header format (no preceding blank lines):
  // Copyright IBM Corp. 2014, 2025
  // SPDX-License-Identifier: MPL-2.0

Example violation:
  package main  // Missing license header!

  func main() {}

Valid usage:
  // Copyright IBM Corp. 2014, 2025
  // SPDX-License-Identifier: MPL-2.0

  package main

  func main() {}`

const azbp005Name = "AZBP005"

var expectedLicenseLines = []string{
	"// Copyright IBM Corp. 2014, 2025",
	"// SPDX-License-Identifier: MPL-2.0",
}

var AZBP005Analyzer = &analysis.Analyzer{
	Name:     azbp005Name,
	Doc:      AZBP005Doc,
	Run:      runAZBP005,
	Requires: []*analysis.Analyzer{commentignore.Analyzer},
}

func runAZBP005(pass *analysis.Pass) (interface{}, error) {
	ignorer, ok := pass.ResultOf[commentignore.Analyzer].(*commentignore.Ignorer)
	if !ok {
		return nil, nil
	}

	for _, file := range pass.Files {
		checkLicenseHeader(pass, file, ignorer)
	}
	return nil, nil
}

func checkLicenseHeader(pass *analysis.Pass, file *ast.File, ignorer *commentignore.Ignorer) {
	filename := pass.Fset.Position(file.Pos()).Filename
	if !strings.HasSuffix(filename, ".go") || !loader.IsFileChanged(filename) {
		return
	}

	if ignorer.ShouldIgnore(azbp005Name, file) || (file.Name != nil && ignorer.ShouldIgnore(azbp005Name, file.Name)) {
		return
	}

	expectedHeader := strings.Join(expectedLicenseLines, "\n")

	// Check: must have comments, first comment before package, starts at line 1
	if len(file.Comments) == 0 || file.Comments[0].Pos() > file.Package {
		reporting.Report(pass, reporting.ReportOptions{
			Rule:          azbp005Name,
			ReportPos:     file.Pos(),
			Message:       fmt.Sprintf("%s: missing license header. Add at the beginning:\n%s\n", azbp005Name, helper.FixedCode(expectedHeader)),
			EvidenceFile:  filename,
			EvidenceLines: []int{1},
			MatchMode:     reporting.MatchModeExactAdded,
		})
		return
	}

	firstComment := file.Comments[0]
	if pass.Fset.Position(firstComment.Pos()).Line != 1 {
		reporting.Report(pass, reporting.ReportOptions{
			Rule:          azbp005Name,
			ReportPos:     file.Pos(),
			Message:       fmt.Sprintf("%s: license header must start at line 1. Expected:\n%s\n", azbp005Name, helper.FixedCode(expectedHeader)),
			EvidenceFile:  filename,
			EvidenceLines: []int{1},
			MatchMode:     reporting.MatchModeExactAdded,
		})
		return
	}

	// Check content matches
	comments := firstComment.List
	if len(comments) < len(expectedLicenseLines) {
		reporting.Report(pass, reporting.ReportOptions{
			Rule:          azbp005Name,
			ReportPos:     firstComment.Pos(),
			Message:       fmt.Sprintf("%s: incomplete license header. Expected:\n%s\n", azbp005Name, helper.FixedCode(expectedHeader)),
			EvidenceFile:  filename,
			EvidenceLines: []int{1},
			MatchMode:     reporting.MatchModeExactAdded,
		})
		return
	}

	for i, expected := range expectedLicenseLines {
		if strings.TrimSpace(comments[i].Text) != expected {
			reporting.Report(pass, reporting.ReportOptions{
				Rule:          azbp005Name,
				ReportPos:     comments[i].Pos(),
				Message:       fmt.Sprintf("%s: incorrect license header. Expected:\n%s\n", azbp005Name, helper.FixedCode(expectedHeader)),
				EvidenceFile:  filename,
				EvidenceLines: []int{1},
				MatchMode:     reporting.MatchModeExactAdded,
			})
			return
		}
	}
}
