// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package customdiff

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/internal/logging"
)

// ComputedIf returns a CustomizeDiffFunc that sets the given key's new value
// as computed if the given condition function returns true.
//
// This function is best effort and will generate a warning log on any errors.
func ComputedIf(key string, f ResourceConditionFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		if f(ctx, d, meta) {
			// To prevent backwards compatibility issues, this logic only
			// generates a warning log instead of returning the error to
			// the provider and ultimately the practitioner. Providers may
			// not be aware of all situations in which the key may not be
			// present in the data, such as during resource creation, so any
			// further changes here should take that into account by
			// documenting how to prevent the error.
			if err := d.SetNewComputed(key); err != nil {
				logging.HelperSchemaWarn(ctx, "unable to set attribute value to unknown", map[string]interface{}{
					logging.KeyAttributePath: key,
					logging.KeyError:         err,
				})
			}
		}
		return nil
	}
}
