package pluginsdk

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

type Set = schema.Set
type SchemaSetFunc = schema.SchemaSetFunc

// NewSet is a convenience method for creating a new set with the given
// items.
func NewSet(f SchemaSetFunc, items []interface{}) *Set {
	return schema.NewSet(f, items)
}

// HashResource hashes complex structures that are described using
// a *Resource. This is the default set implementation used when a set's
// element type is a full resource.
func HashResource(resource *Resource) SchemaSetFunc {
	return schema.HashResource(resource)
}
