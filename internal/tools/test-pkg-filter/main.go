// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// test-pkg-filter scans *_test.go files using go/parser and Go's testing -run regexp rules
// to emit only the package directories that contain matching tests. This avoids compiling
// and linking hundreds of ~300-500 MB full-provider test binaries for packages with zero
// matching tests.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func main() {
	mode := flag.String("mode", "unit", "filter mode: 'unit' (non-TestAcc tests repo-wide) or 'acc' (scoped service tests matching -run)")
	rootPath := flag.String("path", ".", "root directory to scan")
	runExpr := flag.String("run", "", "go test -run regular expression (used in 'acc' mode)")

	flag.Parse()

	pkgs, err := findMatchingPackages(*mode, *rootPath, *runExpr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "test-pkg-filter: %v\n", err)
		os.Exit(1)
	}

	if *mode == "unit" && len(pkgs) == 0 {
		fmt.Fprintln(os.Stderr, "test-pkg-filter: found 0 unit-test packages")
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		fmt.Println(pkg)
	}
}

func findMatchingPackages(mode, rootPath, runExpr string) ([]string, error) {
	cleanRoot := filepath.Clean(rootPath)

	info, err := os.Stat(cleanRoot)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", cleanRoot)
	}

	var topRe *regexp.Regexp

	hasSubtestFilter := false

	if mode == "acc" && runExpr != "" {
		parts := splitTestingRegexp(runExpr)
		hasSubtestFilter = len(parts) > 1

		re, err := regexp.Compile(parts[0])
		if err != nil {
			// Let `go test` itself report invalid -run syntax on the root package without linking subpackages.
			return []string{formatRelPkg(cleanRoot)}, nil
		}

		topRe = re
	} else if mode != "unit" && mode != "acc" {
		return nil, fmt.Errorf("unsupported -mode %q (expected 'unit' or 'acc')", mode)
	}

	filesByDir := make(map[string][]string)

	err = filepath.WalkDir(cleanRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			if path != cleanRoot && shouldSkipDir(d.Name()) {
				return filepath.SkipDir
			}

			return nil
		}

		if strings.HasSuffix(d.Name(), "_test.go") {
			dir := filepath.Dir(path)
			filesByDir[dir] = append(filesByDir[dir], path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	var matched []string

	fset := token.NewFileSet()
	for dir, files := range filesByDir {
		if dirMatches(fset, files, mode, topRe, hasSubtestFilter) {
			matched = append(matched, formatRelPkg(dir))
		}
	}

	sort.Strings(matched)

	if mode == "acc" && len(matched) == 0 {
		// Emit the service root package (which compiles no subpackage binaries) so `go test -json`
		// still emits a valid package event for teamcity-test-reporter when a filter matches 0 tests.
		return []string{formatRelPkg(cleanRoot)}, nil
	}

	return matched, nil
}

func shouldSkipDir(name string) bool {
	if name == "vendor" || name == "testdata" || name == "examples" {
		return true
	}

	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}

func formatRelPkg(dir string) string {
	clean := filepath.ToSlash(filepath.Clean(dir))
	if clean == "." {
		return "."
	}

	if strings.HasPrefix(clean, "./") || filepath.IsAbs(clean) {
		return clean
	}

	return "./" + clean
}

func dirMatches(fset *token.FileSet, files []string, mode string, topRe *regexp.Regexp, hasSubtestFilter bool) bool {
	for _, file := range files {
		parsed, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			// If a test file has a syntax error, include its package so `go test` reports the compile failure.
			return true
		}

		for _, decl := range parsed.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || fn.Name == nil {
				continue
			}

			name := fn.Name.Name

			switch mode {
			case "unit":
				if isUnitTestTarget(name) {
					return true
				}
			case "acc":
				if isGoTestFunc(name, "Test") && name != "TestMain" {
					if topRe == nil || topRe.MatchString(name) {
						return true
					}
				}
			}
		}

		if mode == "acc" && hasSubtestFilter && fileCallsRunTestsInSequence(parsed) {
			return true
		}
	}

	return false
}

func isUnitTestTarget(name string) bool {
	if name == "TestMain" || strings.HasPrefix(name, "TestAcc") {
		return false
	}

	return isGoTestFunc(name, "Test") ||
		isGoTestFunc(name, "Benchmark") ||
		isGoTestFunc(name, "Example") ||
		isGoTestFunc(name, "Fuzz")
}

// isGoTestFunc mirrors Go's cmd/go/internal/load.isTest constraint: prefix followed by
// end-of-string or a non-lowercase rune (e.g. TestFoo, Test_foo, Test1, but not Testing).
func isGoTestFunc(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}

	if len(name) == len(prefix) {
		return true
	}

	r := rune(name[len(prefix)])

	return r < 'a' || r > 'z'
}

func fileCallsRunTestsInSequence(file *ast.File) bool {
	found := false

	ast.Inspect(file, func(n ast.Node) bool {
		if found {
			return false
		}

		sel, ok := n.(*ast.SelectorExpr)
		if ok && sel.Sel != nil && sel.Sel.Name == "RunTestsInSequence" {
			found = true

			return false
		}

		return true
	})

	return found
}

// splitTestingRegexp mirrors Go stdlib testing/match.go splitRegexp so subtest slashes
// (outside parentheses/brackets) split cleanly from the top-level Test function pattern.
func splitTestingRegexp(s string) []string {
	a := make([]string, 0, 2)
	cs := 0
	cp := 0

	for i := 0; i < len(s); {
		switch s[i] {
		case '[':
			cs++
		case ']':
			if cs--; cs < 0 {
				cs = 0
			}
		case '(':
			if cs == 0 {
				cp++
			}
		case ')':
			if cs == 0 {
				cp--
			}
		case '\\':
			i++
		case '/':
			if cs == 0 && cp == 0 {
				a = append(a, s[:i])
				s = s[i+1:]
				i = 0

				continue
			}
		}

		i++
	}

	return append(a, s)
}
