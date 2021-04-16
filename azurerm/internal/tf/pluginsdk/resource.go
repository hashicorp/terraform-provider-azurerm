package pluginsdk

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// This file is intended to provide a transition from Plugin SDKv1 to Plugin SDKv2
// without introducing a merge conflict into every PR.

type SchemaDiffSuppressFunc = schema.SchemaDiffSuppressFunc
type Resource = schema.Resource
type ResourceData = schema.ResourceData
type ResourceDiff = schema.ResourceDiff
type ResourceImporter = schema.ResourceImporter
type SchemaValidateFunc schema.SchemaValidateFunc
type StateUpgrader = schema.StateUpgrader

// ImportStatePassthrough is an implementation of StateFunc that can be
// used to simply pass the ID directly through. This should be used only
// in the case that an ID-only refresh is possible.
func ImportStatePassthrough(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return schema.ImportStatePassthrough(d, m)
}
