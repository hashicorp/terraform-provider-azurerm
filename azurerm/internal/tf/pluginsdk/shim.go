package pluginsdk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// CustomizeDiffShim is a shim around the Terraform Plugin SDK
// which allows us to switch out the version of the Plugin SDK being used
// without breaking open PR's
func CustomizeDiffShim(diffFunc CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(diff *schema.ResourceDiff, i interface{}) error {
		return diffFunc(context.TODO(), diff, i)
	}
}

// ValueChangeConditionShim is a shim around the Terraform Plugin SDK
// which allows us to switch out the version of the Plugin SDK being used
// without breaking open PR's
func ValueChangeConditionShim(shimFunc ValueChangeConditionFunc) customdiff.ValueChangeConditionFunc {
	return func(old, new, meta interface{}) bool {
		return shimFunc(context.TODO(), old, new, meta)
	}
}
