// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func coerceBoolTypable(ctx context.Context, schemaPath path.Path, valuable basetypes.BoolValuable) (basetypes.BoolTypable, diag.Diagnostics) {
	typable, ok := valuable.Type(ctx).(basetypes.BoolTypable)

	// Type() of a Valuable should always be a Typable to recreate the Valuable,
	// but if for some reason it is not, raise an implementation error instead
	// of a panic.
	if !ok {
		return nil, diag.Diagnostics{
			attributePlanModificationTypableError(schemaPath, valuable),
		}
	}

	return typable, nil
}

func coerceFloat64Typable(ctx context.Context, schemaPath path.Path, valuable basetypes.Float64Valuable) (basetypes.Float64Typable, diag.Diagnostics) {
	typable, ok := valuable.Type(ctx).(basetypes.Float64Typable)

	// Type() of a Valuable should always be a Typable to recreate the Valuable,
	// but if for some reason it is not, raise an implementation error instead
	// of a panic.
	if !ok {
		return nil, diag.Diagnostics{
			attributePlanModificationTypableError(schemaPath, valuable),
		}
	}

	return typable, nil
}

func coerceInt64Typable(ctx context.Context, schemaPath path.Path, valuable basetypes.Int64Valuable) (basetypes.Int64Typable, diag.Diagnostics) {
	typable, ok := valuable.Type(ctx).(basetypes.Int64Typable)

	// Type() of a Valuable should always be a Typable to recreate the Valuable,
	// but if for some reason it is not, raise an implementation error instead
	// of a panic.
	if !ok {
		return nil, diag.Diagnostics{
			attributePlanModificationTypableError(schemaPath, valuable),
		}
	}

	return typable, nil
}

func coerceListTypable(ctx context.Context, schemaPath path.Path, valuable basetypes.ListValuable) (basetypes.ListTypable, diag.Diagnostics) {
	typable, ok := valuable.Type(ctx).(basetypes.ListTypable)

	// Type() of a Valuable should always be a Typable to recreate the Valuable,
	// but if for some reason it is not, raise an implementation error instead
	// of a panic.
	if !ok {
		return nil, diag.Diagnostics{
			attributePlanModificationTypableError(schemaPath, valuable),
		}
	}

	return typable, nil
}

func coerceMapTypable(ctx context.Context, schemaPath path.Path, valuable basetypes.MapValuable) (basetypes.MapTypable, diag.Diagnostics) {
	typable, ok := valuable.Type(ctx).(basetypes.MapTypable)

	// Type() of a Valuable should always be a Typable to recreate the Valuable,
	// but if for some reason it is not, raise an implementation error instead
	// of a panic.
	if !ok {
		return nil, diag.Diagnostics{
			attributePlanModificationTypableError(schemaPath, valuable),
		}
	}

	return typable, nil
}

func coerceNumberTypable(ctx context.Context, schemaPath path.Path, valuable basetypes.NumberValuable) (basetypes.NumberTypable, diag.Diagnostics) {
	typable, ok := valuable.Type(ctx).(basetypes.NumberTypable)

	// Type() of a Valuable should always be a Typable to recreate the Valuable,
	// but if for some reason it is not, raise an implementation error instead
	// of a panic.
	if !ok {
		return nil, diag.Diagnostics{
			attributePlanModificationTypableError(schemaPath, valuable),
		}
	}

	return typable, nil
}

func coerceObjectTypable(ctx context.Context, schemaPath path.Path, valuable basetypes.ObjectValuable) (basetypes.ObjectTypable, diag.Diagnostics) {
	typable, ok := valuable.Type(ctx).(basetypes.ObjectTypable)

	// Type() of a Valuable should always be a Typable to recreate the Valuable,
	// but if for some reason it is not, raise an implementation error instead
	// of a panic.
	if !ok {
		return nil, diag.Diagnostics{
			attributePlanModificationTypableError(schemaPath, valuable),
		}
	}

	return typable, nil
}

func coerceSetTypable(ctx context.Context, schemaPath path.Path, valuable basetypes.SetValuable) (basetypes.SetTypable, diag.Diagnostics) {
	typable, ok := valuable.Type(ctx).(basetypes.SetTypable)

	// Type() of a Valuable should always be a Typable to recreate the Valuable,
	// but if for some reason it is not, raise an implementation error instead
	// of a panic.
	if !ok {
		return nil, diag.Diagnostics{
			attributePlanModificationTypableError(schemaPath, valuable),
		}
	}

	return typable, nil
}

func coerceStringTypable(ctx context.Context, schemaPath path.Path, valuable basetypes.StringValuable) (basetypes.StringTypable, diag.Diagnostics) {
	typable, ok := valuable.Type(ctx).(basetypes.StringTypable)

	// Type() of a Valuable should always be a Typable to recreate the Valuable,
	// but if for some reason it is not, raise an implementation error instead
	// of a panic.
	if !ok {
		return nil, diag.Diagnostics{
			attributePlanModificationTypableError(schemaPath, valuable),
		}
	}

	return typable, nil
}

func coerceDynamicTypable(ctx context.Context, schemaPath path.Path, valuable basetypes.DynamicValuable) (basetypes.DynamicTypable, diag.Diagnostics) {
	typable, ok := valuable.Type(ctx).(basetypes.DynamicTypable)

	// Type() of a Valuable should always be a Typable to recreate the Valuable,
	// but if for some reason it is not, raise an implementation error instead
	// of a panic.
	if !ok {
		return nil, diag.Diagnostics{
			attributePlanModificationTypableError(schemaPath, valuable),
		}
	}

	return typable, nil
}
