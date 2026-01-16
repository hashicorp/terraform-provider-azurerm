package passes

import (
	"go/ast"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
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
	Name: azbp005Name,
	Doc:  AZBP005Doc,
	Run:  runAZBP005,
}

func runAZBP005(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		checkLicenseHeader(pass, file)
	}
	return nil, nil
}

func checkLicenseHeader(pass *analysis.Pass, file *ast.File) {
	pos := pass.Fset.Position(file.Pos())
	if !loader.ShouldReport(pos.Filename, 1) {
		return
	}

	expectedHeader := strings.Join(expectedLicenseLines, "\n")

	// Check: must have comments, first comment before package, starts at line 1
	if len(file.Comments) == 0 || file.Comments[0].Pos() > file.Package {
		pass.Reportf(file.Pos(), "%s: missing license header. Add at the beginning:\n%s\n",
			azbp005Name, helper.FixedCode(expectedHeader))
		return
	}

	firstComment := file.Comments[0]
	if pass.Fset.Position(firstComment.Pos()).Line != 1 {
		pass.Reportf(file.Pos(), "%s: license header must start at line 1. Expected:\n%s\n",
			azbp005Name, helper.FixedCode(expectedHeader))
		return
	}

	// Check content matches
	comments := firstComment.List
	if len(comments) < len(expectedLicenseLines) {
		pass.Reportf(firstComment.Pos(), "%s: incomplete license header. Expected:\n%s\n",
			azbp005Name, helper.FixedCode(expectedHeader))
		return
	}

	for i, expected := range expectedLicenseLines {
		if strings.TrimSpace(comments[i].Text) != expected {
			pass.Reportf(comments[i].Pos(), "%s: incorrect license header. Expected:\n%s\n",
				azbp005Name, helper.FixedCode(expectedHeader))
			return
		}
	}
}
