package pluginsdk

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

type Set = schema.Set
type SchemaSetFunc = schema.SchemaSetFunc

func NewSet(f SchemaSetFunc, items []interface{}) *Set {
	return schema.NewSet(f, items)
}
