// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfsdk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Config represents a Terraform config.
type Config struct {
	Raw    tftypes.Value
	Schema fwschema.Schema
}

// Get populates the struct passed as `target` with the entire config.
func (c Config) Get(ctx context.Context, target interface{}) diag.Diagnostics {
	return c.data().Get(ctx, target)
}

// GetAttribute retrieves the attribute or block found at `path` and populates
// the `target` with the value. This method is intended for top level schema
// attributes or blocks. Use `types` package methods or custom types to step
// into collections.
//
// Attributes or elements under null or unknown collections return null
// values, however this behavior is not protected by compatibility promises.
func (c Config) GetAttribute(ctx context.Context, path path.Path, target interface{}) diag.Diagnostics {
	return c.data().GetAtPath(ctx, path, target)
}

// PathMatches returns all matching path.Paths from the given path.Expression.
//
// If a parent path is null or unknown, which would prevent a full expression
// from matching, the parent path is returned rather than no match to prevent
// false positives.
func (c Config) PathMatches(ctx context.Context, pathExpr path.Expression) (path.Paths, diag.Diagnostics) {
	return c.data().PathMatches(ctx, pathExpr)
}

func (c Config) data() fwschemadata.Data {
	return fwschemadata.Data{
		Description:    fwschemadata.DataDescriptionConfiguration,
		Schema:         c.Schema,
		TerraformValue: c.Raw,
	}
}
