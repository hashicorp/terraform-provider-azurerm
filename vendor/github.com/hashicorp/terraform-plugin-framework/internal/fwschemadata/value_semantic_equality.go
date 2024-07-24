// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ValueSemanticEqualityRequest represents a request for the provider to
// perform semantic equality logic on a value.
type ValueSemanticEqualityRequest struct {
	// Path is the schema-based path of the value.
	Path path.Path

	// PriorValue is the prior value.
	PriorValue attr.Value

	// ProposedNewValue is the proposed new value. NewValue in the response
	// contains the results of semantic equality logic.
	ProposedNewValue attr.Value
}

// ValueSemanticEqualityResponse represents a response to a
// ValueSemanticEqualityRequest.
type ValueSemanticEqualityResponse struct {
	// NewValue contains the new value based on the semantic equality logic.
	NewValue attr.Value

	// Diagnostics contains any errors and warnings for the logic.
	Diagnostics diag.Diagnostics
}

// ValueSemanticEquality runs all semantic equality logic for a value, including
// recursive checking against collection and structural types.
func ValueSemanticEquality(ctx context.Context, req ValueSemanticEqualityRequest, resp *ValueSemanticEqualityResponse) {
	ctx = logging.FrameworkWithAttributePath(ctx, req.Path.String())

	// Ensure the response NewValue always starts with the proposed new value.
	// This is purely defensive coding to prevent subtle data handling bugs.
	resp.NewValue = req.ProposedNewValue

	// If the prior value is null or unknown, no need to check semantic equality
	// as the proposed new value is always correct. There is also no need to
	// descend further into any nesting.
	if req.PriorValue.IsNull() || req.PriorValue.IsUnknown() {
		return
	}

	// If the proposed new value is null or unknown, no need to check semantic
	// equality as it should never be changed back to the prior value. There is
	// also no need to descend further into any nesting.
	if req.ProposedNewValue.IsNull() || req.ProposedNewValue.IsUnknown() {
		return
	}

	switch req.ProposedNewValue.(type) {
	case basetypes.BoolValuable:
		ValueSemanticEqualityBool(ctx, req, resp)
	case basetypes.Float64Valuable:
		ValueSemanticEqualityFloat64(ctx, req, resp)
	case basetypes.Int64Valuable:
		ValueSemanticEqualityInt64(ctx, req, resp)
	case basetypes.ListValuable:
		ValueSemanticEqualityList(ctx, req, resp)
	case basetypes.MapValuable:
		ValueSemanticEqualityMap(ctx, req, resp)
	case basetypes.NumberValuable:
		ValueSemanticEqualityNumber(ctx, req, resp)
	case basetypes.ObjectValuable:
		ValueSemanticEqualityObject(ctx, req, resp)
	case basetypes.SetValuable:
		ValueSemanticEqualitySet(ctx, req, resp)
	case basetypes.StringValuable:
		ValueSemanticEqualityString(ctx, req, resp)
	case basetypes.DynamicValuable:
		ValueSemanticEqualityDynamic(ctx, req, resp)
	}

	if resp.NewValue.Equal(req.PriorValue) {
		logging.FrameworkDebug(ctx, "Value switched to prior value due to semantic equality logic")
	}
}
