package azure

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// shared schema
func MergeSchema(a map[string]*pluginsdk.Schema, b map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
	// TODO: Deprecate and remove this

	s := map[string]*pluginsdk.Schema{}

	for k, v := range a {
		s[k] = v
	}

	for k, v := range b {
		s[k] = v
	}

	return s
}
