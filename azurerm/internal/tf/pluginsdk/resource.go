package pluginsdk

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// This file is intended to provide a transition from Plugin SDKv1 to Plugin SDKv2
// without introducing a merge conflict into every PR.

type BasicMapReader = schema.BasicMapReader
type MapFieldReader = schema.MapFieldReader
type MapFieldWriter = schema.MapFieldWriter
type Resource = schema.Resource
type ResourceData = schema.ResourceData
type ResourceDiff = schema.ResourceDiff
type SchemaDiffSuppressFunc = schema.SchemaDiffSuppressFunc
type StateUpgrader = schema.StateUpgrader
type SchemaValidateFunc = schema.SchemaValidateFunc
type ValueType = schema.ValueType

// ImportStatePassthrough is an implementation of StateFunc that can be
// used to simply pass the ID directly through. This should be used only
// in the case that an ID-only refresh is possible.
func ImportStatePassthrough(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return schema.ImportStatePassthrough(d, m)
}

type StateChangeConf = resource.StateChangeConf
type StateRefreshFunc = resource.StateRefreshFunc

type CreateFunc = schema.CreateFunc
type DeleteFunc = schema.DeleteFunc
type ExistsFunc = schema.ExistsFunc
type ReadFunc = schema.ReadFunc
type UpdateFunc = schema.UpdateFunc
