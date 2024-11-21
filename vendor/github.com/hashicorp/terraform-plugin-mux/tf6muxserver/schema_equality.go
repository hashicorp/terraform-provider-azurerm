// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// schemaCmpOptions ensures comparisons of SchemaAttribute and
// SchemaNestedBlock slices are considered equal despite ordering differences.
var schemaCmpOptions = []cmp.Option{
	cmpopts.SortSlices(func(i, j *tfprotov6.SchemaAttribute) bool {
		return i.Name < j.Name
	}),
	cmpopts.SortSlices(func(i, j *tfprotov6.SchemaNestedBlock) bool {
		return i.TypeName < j.TypeName
	}),
	cmpopts.IgnoreFields(tfprotov6.SchemaNestedBlock{}, "MinItems", "MaxItems"),
}

// schemaDiff outputs the difference between schemas while accounting for
// inconsequential ordering differences in SchemaAttribute and
// SchemaNestedBlock slices.
func schemaDiff(i, j *tfprotov6.Schema) string {
	return cmp.Diff(i, j, schemaCmpOptions...)
}

// schemaEquals asserts equality between schemas by normalizing inconsequential
// ordering differences in SchemaAttribute and SchemaNestedBlock slices.
func schemaEquals(i, j *tfprotov6.Schema) bool {
	return cmp.Equal(i, j, schemaCmpOptions...)
}
