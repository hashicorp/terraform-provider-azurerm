// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSliceDiff(t *testing.T) {
	type args struct {
		want []string
		got  []string
	}
	tests := []struct {
		name     string
		args     args
		wantDiff int
	}{
		{
			name: "aaa",
			args: args{
				want: []string{"abc", "def"},
				got:  []string{"def", "abc"},
			},
			wantDiff: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if missed, odd := SliceDiff(tt.args.want, tt.args.got, true); len(missed)+len(odd) != tt.wantDiff {
				t.Errorf("SliceDiff() missed: %v, odd: %v", missed, odd)
			}
			if diff := cmp.Diff(tt.args.want, tt.args.got); diff == "" {
				t.Errorf("%s should have diff in cmp.Diff", tt.name)
			}
		})
	}
}

func TestDiffAll(t *testing.T) {
	t.Skipf("skip for fix all documents")
	result := DiffAll(AzurermAllResources("", "", "", ""), true)
	if err := result.FixDocuments(); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", result.ToString())
}

func TestMinimumOf3(t *testing.T) {
	args := [][4]int{
		{1, 2, 3, 1},
		{3, 2, 1, 1},
		{10, 32, -1, -1},
		{0, 0, 0, 0},
		{-1, -2, -3, -3},
	}
	for _, arg := range args {
		if got := minimumOf3(arg[0], arg[1], arg[2]); got != arg[3] {
			t.Fatalf("want min %d got %d", arg[3], got)
		}
	}
}

func TestLevenshteinDist(t *testing.T) {
	args := []struct {
		str1 string
		str2 string
		dist int
	}{
		{"Python", "Peithen", 3},
		{"azure", "Azure", 1},
		{"azurerm", "awsrm", 4},
	}

	for _, arg := range args {
		got := levenshteinDist(arg.str1, arg.str2)
		if got != arg.dist {
			t.Fatalf("want dist %d, got %d", arg.dist, got)
		}
	}
}
