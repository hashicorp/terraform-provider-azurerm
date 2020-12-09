package azure

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// shared schema
func MergeSchema(a map[string]*schema.Schema, b map[string]*schema.Schema) map[string]*schema.Schema {
	// TODO: Deprecate and remove this

	s := map[string]*schema.Schema{}

	for k, v := range a {
		s[k] = v
	}

	for k, v := range b {
		s[k] = v
	}

	return s
}
