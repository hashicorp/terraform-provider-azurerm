package pluginsdk

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type Schema = schema.Schema

const (
	TypeInvalid = schema.TypeInvalid
	TypeBool    = schema.TypeBool
	TypeInt     = schema.TypeInt
	TypeFloat   = schema.TypeFloat
	TypeString  = schema.TypeString
	TypeList    = schema.TypeList
	TypeMap     = schema.TypeMap
	TypeSet     = schema.TypeSet
)

const (
	SchemaConfigModeAuto  = schema.SchemaConfigModeAuto
	SchemaConfigModeAttr  = schema.SchemaConfigModeAttr
	SchemaConfigModeBlock = schema.SchemaConfigModeBlock
)
