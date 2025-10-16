// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	schema2 "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/schema"
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

func TestComputedMissingInAttributes(t *testing.T) {
	testSchema := &schema.Resource{
		ResourceType: "azurerm_test_resource",
		Schema: &schema2.Resource{
			Schema: map[string]*schema2.Schema{
				"name": {
					Type:     schema2.TypeString,
					Required: true,
				},
				"computed_field": {
					Type:     schema2.TypeString,
					Computed: true,
				},
			},
		},
	}

	testDoc := &model.ResourceDoc{
		Args: model.Properties{
			"name": &model.Field{
				Name: "name",
				Pos:  model.PosArgs,
			},
		},
		Attr: model.Properties{},
	}

	checkers := crossCheckProperty(testSchema, testDoc)

	foundComputedMissing := false
	for _, checker := range checkers {
		if missChecker, ok := checker.(*propertyMissDiff); ok {
			if missChecker.key == "computed_field" && missChecker.MissType == MissInDocAttr {
				foundComputedMissing = true
				break
			}
		}
	}

	if !foundComputedMissing {
		t.Error("Expected to find MissInDocAttr error for computed field 'computed_field', but didn't find it")
	}
}

func TestNestedComputedMissingInAttributes(t *testing.T) {
	testSchema := &schema.Resource{
		ResourceType: "azurerm_test_resource",
		Schema: &schema2.Resource{
			Schema: map[string]*schema2.Schema{
				"name": {
					Type:     schema2.TypeString,
					Required: true,
				},
				"block_field": {
					Type:     schema2.TypeList,
					Optional: true,
					Elem: &schema2.Resource{
						Schema: map[string]*schema2.Schema{
							"normal_field": {
								Type:     schema2.TypeString,
								Optional: true,
							},
							"nested_computed_field": {
								Type:     schema2.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}

	testDoc := &model.ResourceDoc{
		Args: model.Properties{
			"name": &model.Field{
				Name: "name",
				Pos:  model.PosArgs,
			},
			"block_field": &model.Field{
				Name: "block_field",
				Pos:  model.PosArgs,
				Subs: model.Properties{
					"normal_field": &model.Field{
						Name: "normal_field",
						Pos:  model.PosArgs,
					},
				},
			},
		},
		Attr: model.Properties{
		},
	}

	checkers := crossCheckProperty(testSchema, testDoc)

	foundNestedComputedMissing := false
	for _, checker := range checkers {
		if missChecker, ok := checker.(*propertyMissDiff); ok {
			if missChecker.key == "block_field.nested_computed_field" && missChecker.MissType == MissInDocAttr {
				foundNestedComputedMissing = true
				break
			}
		}
	}

	if !foundNestedComputedMissing {
		t.Error("Expected to find MissInDocAttr error for nested computed field 'block_field.nested_computed_field', but didn't find it")
	}
}

func TestOptionalComputedExistsInAttributes(t *testing.T) {
	testSchema := &schema.Resource{
		ResourceType: "azurerm_test_resource",
		Schema: &schema2.Resource{
			Schema: map[string]*schema2.Schema{
				"oc_field": {
					Type:     schema2.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}

	testDoc := &model.ResourceDoc{
		Args: model.Properties{
			"oc_field": &model.Field{
				Name: "name",
				Pos:  model.PosArgs,
			},
		},
		Attr: model.Properties{
		},
	}

	checkers := crossCheckProperty(testSchema, testDoc)

	foundNestedComputedMissing := false
	for _, checker := range checkers {
		if missChecker, ok := checker.(*propertyMissDiff); ok {
			if missChecker.key == "oc_field" && missChecker.MissType == MissInDocAttr {
				foundNestedComputedMissing = true
				break
			}
		}
	}

	if foundNestedComputedMissing {
		t.Error("Expected `oc_field` to exist in Argument field, but receive MissInDocAttr")
	}
}
