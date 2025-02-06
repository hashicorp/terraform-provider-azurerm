// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package customdiff

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// All returns a CustomizeDiffFunc that runs all of the given
// CustomizeDiffFuncs and returns all of the errors produced.
//
// If one function produces an error, functions after it are still run.
// If this is not desirable, use function Sequence instead.
//
// If multiple functions returns errors, the result is a multierror.
//
// For example:
//
//	&schema.Resource{
//	    // ...
//	    CustomizeDiff: customdiff.All(
//	        customdiff.ValidateChange("size", func (ctx context.Context, old, new, meta interface{}) error {
//	            // If we are increasing "size" then the new value must be
//	            // a multiple of the old value.
//	            if new.(int) <= old.(int) {
//	                return nil
//	            }
//	            if (new.(int) % old.(int)) != 0 {
//	                return fmt.Errorf("new size value must be an integer multiple of old value %d", old.(int))
//	            }
//	            return nil
//	        }),
//	        customdiff.ForceNewIfChange("size", func (ctx context.Context, old, new, meta interface{}) bool {
//	            // "size" can only increase in-place, so we must create a new resource
//	            // if it is decreased.
//	            return new.(int) < old.(int)
//	        }),
//	        customdiff.ComputedIf("version_id", func (ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
//	            // Any change to "content" causes a new "version_id" to be allocated.
//	            return d.HasChange("content")
//	        }),
//	    ),
//	}
func All(funcs ...schema.CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		var errs []error
		for _, f := range funcs {
			thisErr := f(ctx, d, meta)
			if thisErr != nil {
				errs = append(errs, thisErr)
			}
		}
		return errors.Join(errs...)
	}
}

// Sequence returns a CustomizeDiffFunc that runs all of the given
// CustomizeDiffFuncs in sequence, stopping at the first one that returns
// an error and returning that error.
//
// If all functions succeed, the combined function also succeeds.
func Sequence(funcs ...schema.CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		for _, f := range funcs {
			err := f(ctx, d, meta)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
