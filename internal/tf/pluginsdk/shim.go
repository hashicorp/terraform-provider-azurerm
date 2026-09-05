// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CustomizeDiffShim is a shim around the Terraform Plugin SDK
// which allows us to switch out the version of the Plugin SDK being used
// without breaking open PR's
func CustomizeDiffShim(diffFunc CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
		return diffFunc(ctx, diff, i)
	}
}
