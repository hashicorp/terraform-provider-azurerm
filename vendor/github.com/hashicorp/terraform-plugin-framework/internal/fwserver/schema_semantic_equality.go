// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// SchemaSemanticEqualityRequest represents a request for a schema to run all
// semantic equality logic.
type SchemaSemanticEqualityRequest struct {
	// PriorData is the prior schema-based data.
	PriorData fwschemadata.Data

	// ProposedNewData is the proposed new schema-based data. The response
	// NewData contains the results of any modifications.
	ProposedNewData fwschemadata.Data
}

// SchemaSemanticEqualityResponse represents a response to a
// SchemaSemanticEqualityRequest.
type SchemaSemanticEqualityResponse struct {
	// NewData is the new schema-based data after any modifications.
	NewData fwschemadata.Data

	// Diagnostics report errors or warnings related to running all attribute
	// plan modifiers. Returning an empty slice indicates a successful
	// plan modification with no warnings or errors generated.
	Diagnostics diag.Diagnostics
}

// SchemaSemanticEquality runs semantic equality logic for all schema attributes
// and blocks.
//
// MAINTAINER NOTE: Since semantic equality is purely value based, where
// attributes and blocks cannot currently introduce semantic equality logic
// based on those schema concepts, this logic immediately delegates to value
// based handling. On the off chance that the framework is enhanced with
// attribute and block level semantic equality support (not recommended since
// value types should really be the correct provider developer abstraction,
// rather than potentially causing confusing or duplicated provider logic), this
// logic will need to be redesigned similar to the plan modification and
// validation logic which walks the schema. That schema walk may interfere with
// the value based recursion for collection and structural types, so additional
// design may be necessary so that provider developer data handling intentions
// are kept based on both the value based logic and schema based logic.
func SchemaSemanticEquality(ctx context.Context, req SchemaSemanticEqualityRequest, resp *SchemaSemanticEqualityResponse) {
	var diags diag.Diagnostics

	for name := range req.ProposedNewData.Schema.GetAttributes() {
		valueReq := fwschemadata.ValueSemanticEqualityRequest{
			Path: path.Root(name),
		}

		valueReq.PriorValue, diags = req.PriorData.ValueAtPath(ctx, valueReq.Path)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		valueReq.ProposedNewValue, diags = req.ProposedNewData.ValueAtPath(ctx, valueReq.Path)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		valueResp := &fwschemadata.ValueSemanticEqualityResponse{
			NewValue: valueReq.ProposedNewValue,
		}

		fwschemadata.ValueSemanticEquality(ctx, valueReq, valueResp)

		resp.Diagnostics.Append(valueResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}

		// If the response value equals the original proposed new value, move
		// to next attribute.
		if valueResp.NewValue.Equal(valueReq.ProposedNewValue) {
			continue
		}

		resp.Diagnostics.Append(resp.NewData.SetAtPath(ctx, valueReq.Path, valueResp.NewValue)...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	for name := range req.ProposedNewData.Schema.GetBlocks() {
		valueReq := fwschemadata.ValueSemanticEqualityRequest{
			Path: path.Root(name),
		}

		valueReq.PriorValue, diags = req.PriorData.ValueAtPath(ctx, valueReq.Path)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		valueReq.ProposedNewValue, diags = req.ProposedNewData.ValueAtPath(ctx, valueReq.Path)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		valueResp := &fwschemadata.ValueSemanticEqualityResponse{
			NewValue: valueReq.ProposedNewValue,
		}

		fwschemadata.ValueSemanticEquality(ctx, valueReq, valueResp)

		resp.Diagnostics.Append(valueResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}

		// If the response value equals the original proposed new value, move
		// to next block.
		if valueResp.NewValue.Equal(valueReq.ProposedNewValue) {
			continue
		}

		resp.Diagnostics.Append(resp.NewData.SetAtPath(ctx, valueReq.Path, valueResp.NewValue)...)

		if resp.Diagnostics.HasError() {
			return
		}
	}
}
