package pluginsdk

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This file is intended to provide a transition from Plugin SDKv1 to Plugin SDKv2
// without introducing a merge conflict into every PR.

type (
	BasicMapReader         = schema.BasicMapReader
	MapFieldReader         = schema.MapFieldReader
	MapFieldWriter         = schema.MapFieldWriter
	Resource               = schema.Resource
	ResourceData           = schema.ResourceData
	ResourceDiff           = schema.ResourceDiff
	SchemaDiffSuppressFunc = schema.SchemaDiffSuppressFunc
	StateUpgrader          = schema.StateUpgrader
	SchemaValidateFunc     = func(interface{}, string) ([]string, []error)
	ValueType              = schema.ValueType
)

type (
	StateChangeConf  = resource.StateChangeConf
	StateRefreshFunc = resource.StateRefreshFunc
)

type (
	CreateFunc = schema.CreateFunc //nolint:SA1019
	DeleteFunc = schema.DeleteFunc
	ExistsFunc = schema.ExistsFunc
	ReadFunc   = schema.ReadFunc
	UpdateFunc = schema.UpdateFunc
)
