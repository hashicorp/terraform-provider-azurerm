package location

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func Schema() *pluginsdk.Schema {
	return commonschema.Location()
}

func SchemaComputed() *pluginsdk.Schema {
	return commonschema.LocationComputed()
}

func DiffSuppressFunc(v, old, new string, d *pluginsdk.ResourceData) bool {
	return location.DiffSuppressFunc(v, old, new, d)
}

func HashCode(location interface{}) int {
	// NOTE: this is intentionally not present upstream as the only usage is deprecated
	// and so this can be removed in 3.0
	loc := location.(string)
	return pluginsdk.HashString(Normalize(loc))
}

func StateFunc(input interface{}) string {
	return location.StateFunc(input)
}
