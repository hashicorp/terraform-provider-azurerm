// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ValidPathExpression returns true if the given expression is valid for the
// schema underlying the Data. This can be used to determine if there was an
// expression implementation error versus an expression returning no path
// matches based on implementation details of the underlying data storage.
func (d Data) ValidPathExpression(ctx context.Context, expression path.Expression) bool {
	expressionSteps := expression.Resolve().Steps()

	if len(expressionSteps) == 0 {
		return false
	}

	return validatePathExpressionSteps(ctx, d.Schema.Type(), expressionSteps)
}

// validatePathExpressionSteps is a recursive function which returns true if
// the path expression steps can be applied to the type.
func validatePathExpressionSteps(ctx context.Context, currentType attr.Type, currentExpressionSteps path.ExpressionSteps) bool {
	currentExpressionStep, nextSteps := currentExpressionSteps.NextStep()

	// Generate a tftypes step based on the expression. For type definitions,
	// any value should be acceptable for element steps.
	var currentTfStep tftypes.AttributePathStep

	switch step := currentExpressionStep.(type) {
	case nil:
		// There are no more expression steps.
		return true
	case path.ExpressionStepAttributeNameExact:
		currentTfStep = tftypes.AttributeName(step)
	case path.ExpressionStepElementKeyIntAny:
		currentTfStep = tftypes.ElementKeyInt(0)
	case path.ExpressionStepElementKeyIntExact:
		currentTfStep = tftypes.ElementKeyInt(step)
	case path.ExpressionStepElementKeyStringAny:
		currentTfStep = tftypes.ElementKeyString("")
	case path.ExpressionStepElementKeyStringExact:
		currentTfStep = tftypes.ElementKeyString(step)
	case path.ExpressionStepElementKeyValueAny:
		tfValue := tftypes.NewValue(
			currentType.TerraformType(ctx),
			nil,
		)
		currentTfStep = tftypes.ElementKeyValue(tfValue)
	case path.ExpressionStepElementKeyValueExact:
		// Best effort
		tfValue, err := step.Value.ToTerraformValue(ctx)

		if err != nil {
			tfValue = tftypes.NewValue(
				currentType.TerraformType(ctx),
				nil,
			)
		}

		currentTfStep = tftypes.ElementKeyValue(tfValue)
	default:
		// If new, resolved path.ExpressionStep are introduced, they must be
		// added as cases to this switch statement.
		panic(fmt.Sprintf("unimplemented path.ExpressionStep type: %T", currentExpressionStep))
	}

	nextTypeIface, err := currentType.ApplyTerraform5AttributePathStep(currentTfStep)

	if err != nil {
		// Debug, not error, log entry for troubleshooting as validation may
		// be running in a scenario where invalid expressions are okay.
		logging.FrameworkDebug(
			ctx,
			fmt.Sprintf("Returning false due to error while calling %T ApplyTerraform5AttributePathStep with %T", currentType, currentTfStep),
			map[string]any{
				logging.KeyError: err,
			},
		)

		return false
	}

	nextType, ok := nextTypeIface.(attr.Type)

	if !ok {
		// Raise a more descriptive panic message instead of the type assertion
		// panic.
		panic(fmt.Sprintf("%T returned unexpected type %T from ApplyTerraform5AttributePathStep", currentType, nextTypeIface))
	}

	return validatePathExpressionSteps(ctx, nextType, nextSteps)
}
