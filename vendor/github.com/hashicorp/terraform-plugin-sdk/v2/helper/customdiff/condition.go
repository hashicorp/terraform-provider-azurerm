// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package customdiff

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceConditionFunc is a function type that makes a boolean decision based
// on an entire resource diff.
type ResourceConditionFunc func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool

// ValueChangeConditionFunc is a function type that makes a boolean decision
// by comparing two values.
type ValueChangeConditionFunc func(ctx context.Context, oldValue, newValue, meta interface{}) bool

// ValueConditionFunc is a function type that makes a boolean decision based
// on a given value.
type ValueConditionFunc func(ctx context.Context, value, meta interface{}) bool

// If returns a CustomizeDiffFunc that calls the given condition
// function and then calls the given CustomizeDiffFunc only if the condition
// function returns true.
//
// This can be used to include conditional customizations when composing
// customizations using All and Sequence, but should generally be used only in
// simple scenarios. Prefer directly writing a CustomizeDiffFunc containing
// a conditional branch if the given CustomizeDiffFunc is already a
// locally-defined function, since this avoids obscuring the control flow.
func If(cond ResourceConditionFunc, f schema.CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		if cond(ctx, d, meta) {
			return f(ctx, d, meta)
		}
		return nil
	}
}

// IfValueChange returns a CustomizeDiffFunc that calls the given condition
// function with the old and new values of the given key and then calls the
// given CustomizeDiffFunc only if the condition function returns true.
func IfValueChange(key string, cond ValueChangeConditionFunc, f schema.CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		oldValue, newValue := d.GetChange(key)
		if cond(ctx, oldValue, newValue, meta) {
			return f(ctx, d, meta)
		}
		return nil
	}
}

// IfValue returns a CustomizeDiffFunc that calls the given condition
// function with the new values of the given key and then calls the
// given CustomizeDiffFunc only if the condition function returns true.
func IfValue(key string, cond ValueConditionFunc, f schema.CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		if cond(ctx, d.Get(key), meta) {
			return f(ctx, d, meta)
		}
		return nil
	}
}
