// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// Unit tests for test-pkg-filter are included because this tool gates which packages
// CI and TeamCity pass to `go test`; regressions in Go test-name boundary rules,
// TestMain/TestAcc exclusion, or subtest slash splitting could silently skip tests.
package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestIsUnitTestTarget(t *testing.T) {
	cases := map[string]bool{
		"Test":                     true,
		"TestFoo":                  true,
		"Test_foo":                 true,
		"Test123":                  true,
		"TestAccess":               false, // CI's `-run='^Test[^A]|^TestA[^c]|^TestAc[^c]'` excludes all `TestAcc*` prefixes
		"TestAccount":              false,
		"TestAccFoo_basic":         false,
		"TestMain":                 false,
		"TestMaintenance":          true,
		"TestingHelper":            false,
		"BenchmarkDecode":          true,
		"ExampleParse":             true,
		"FuzzResourceID":           true,
		"benchmarkUnexported":      false,
		"testAccSequentialSubtest": false,
	}

	for name, want := range cases {
		if got := isUnitTestTarget(name); got != want {
			t.Errorf("isUnitTestTarget(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestSplitTestingRegexp(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"TestAccLinuxVM", []string{"TestAccLinuxVM"}},
		{"TestAccSeq/group/basic", []string{"TestAccSeq", "group", "basic"}},
		{"(TestAccA/sub|TestAccB)/case", []string{"(TestAccA/sub|TestAccB)", "case"}},
		{"TestAcc[/]Name/sub", []string{"TestAcc[/]Name", "sub"}},
	}
	for _, tc := range cases {
		if got := splitTestingRegexp(tc.in); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("splitTestingRegexp(%q) = %#v, want %#v", tc.in, got, tc.want)
		}
	}
}

func TestFindMatchingPackages(t *testing.T) {
	tmp := t.TempDir()

	write := func(rel, content string) {
		t.Helper()

		full := filepath.Join(tmp, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("parse/id_test.go", "package parse\nimport \"testing\"\nfunc TestID(t *testing.T) {}\n")
	write("resources/res_a/resource_test.go", "package res_a_test\nimport \"testing\"\nfunc TestAccResA_basic(t *testing.T) {}\n")
	write("resources/res_b/resource_test.go", "package res_b_test\nimport \"testing\"\nfunc TestAccResB_basic(t *testing.T) {}\n")
	write("resources/res_seq/resource_test.go", "package res_seq_test\nimport (\"testing\"; \"example/acceptance\")\nfunc TestAccSeq(t *testing.T) { acceptance.RunTestsInSequence(t, nil) }\n")
	write(filepath.Join("testdata", "ignored_test.go"), "package testdata\nimport \"testing\"\nfunc TestIgnored(t *testing.T) {}\n")

	unitPkgs, err := findMatchingPackages("unit", tmp, "")
	if err != nil {
		t.Fatalf("unit mode failed: %v", err)
	}

	wantUnit := []string{filepath.ToSlash(filepath.Join(tmp, "parse"))}
	if !reflect.DeepEqual(unitPkgs, wantUnit) {
		t.Errorf("unit mode = %v, want %v", unitPkgs, wantUnit)
	}

	accScoped, err := findMatchingPackages("acc", tmp, "TestAccResB")
	if err != nil {
		t.Fatalf("acc scoped failed: %v", err)
	}

	wantScoped := []string{filepath.ToSlash(filepath.Join(tmp, "resources/res_b"))}
	if !reflect.DeepEqual(accScoped, wantScoped) {
		t.Errorf("acc scoped = %v, want %v", accScoped, wantScoped)
	}

	accNone, err := findMatchingPackages("acc", tmp, "TestAccDoesNotExist")
	if err != nil {
		t.Fatalf("acc empty match failed: %v", err)
	}

	wantFallback := []string{filepath.ToSlash(tmp)}
	if !reflect.DeepEqual(accNone, wantFallback) {
		t.Errorf("acc empty fallback = %v, want %v", accNone, wantFallback)
	}
}
