// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package customdiff

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/internal/logging"
)

// ForceNewIf returns a CustomizeDiffFunc that flags the given key as
// requiring a new resource if the given condition function returns true.
//
// The return value of the condition function is ignored if the old and new
// values of the field compare equal, since no attribute diff is generated in
// that case.
//
// This function is best effort and will generate a warning log on any errors.
func ForceNewIf(key string, f ResourceConditionFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		if f(ctx, d, meta) {
			// To prevent backwards compatibility issues, this logic only
			// generates a warning log instead of returning the error to
			// the provider and ultimately the practitioner. Providers may
			// not be aware of all situations in which the key may not be
			// present in the data, such as during resource creation, so any
			// further changes here should take that into account by
			// documenting how to prevent the error.
			if err := d.ForceNew(key); err != nil {
				logging.HelperSchemaWarn(ctx, "unable to require attribute replacement", map[string]interface{}{
					logging.KeyAttributePath: key,
					logging.KeyError:         err,
				})
			}
		}
		return nil
	}
}

// ForceNewIfChange returns a CustomizeDiffFunc that flags the given key as
// requiring a new resource if the given condition function returns true.
//
// The return value of the condition function is ignored if the old and new
// values compare equal, since no attribute diff is generated in that case.
//
// This function is similar to ForceNewIf but provides the condition function
// only the old and new values of the given key, which leads to more compact
// and explicit code in the common case where the decision can be made with
// only the specific field value.
//
// This function is best effort and will generate a warning log on any errors.
func ForceNewIfChange(key string, f ValueChangeConditionFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		oldValue, newValue := d.GetChange(key)
		if f(ctx, oldValue, newValue, meta) {
			// To prevent backwards compatibility issues, this logic only
			// generates a warning log instead of returning the error to
			// the provider and ultimately the practitioner. Providers may
			// not be aware of all situations in which the key may not be
			// present in the data, such as during resource creation, so any
			// further changes here should take that into account by
			// documenting how to prevent the error.
			if err := d.ForceNew(key); err != nil {
				logging.HelperSchemaWarn(ctx, "unable to require attribute replacement", map[string]interface{}{
					logging.KeyAttributePath: key,
					logging.KeyError:         err,
				})
			}
		}
		return nil
	}
}
