package compute

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// nolint: deadcode unused
func adminPasswordDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	// this is not the greatest hack in the world, this is just a tribute.
	if old == "ignored-as-imported" || new == "ignored-as-imported" {
		return true
	}

	return false
}
