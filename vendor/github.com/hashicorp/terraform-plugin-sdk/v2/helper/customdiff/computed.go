package customdiff

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ComputedIf returns a CustomizeDiffFunc that sets the given key's new value
// as computed if the given condition function returns true.
func ComputedIf(key string, f ResourceConditionFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		if f(ctx, d, meta) {
			if err := d.SetNewComputed(key); err != nil {
				return fmt.Errorf("unable to set %q to unknown: %w", key, err)
			}
		}
		return nil
	}
}
