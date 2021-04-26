package location

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

// Schema returns the default Schema which should be used for Location fields
// where these are Required and Cannot be Changed
func Schema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:             pluginsdk.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateFunc:     EnhancedValidate,
		StateFunc:        StateFunc,
		DiffSuppressFunc: DiffSuppressFunc,
	}
}

// SchemaOptional returns the Schema for a Location field where this can be optionally specified
func SchemaOptional() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:             pluginsdk.TypeString,
		Optional:         true,
		ForceNew:         true,
		StateFunc:        StateFunc,
		DiffSuppressFunc: DiffSuppressFunc,
	}
}

func SchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Computed: true,
	}
}

// Schema returns the Schema which should be used for Location fields
// where these are Required and can be changed
func SchemaWithoutForceNew() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:             pluginsdk.TypeString,
		Required:         true,
		ValidateFunc:     EnhancedValidate,
		StateFunc:        StateFunc,
		DiffSuppressFunc: DiffSuppressFunc,
	}
}

func DiffSuppressFunc(_, old, new string, _ *pluginsdk.ResourceData) bool {
	return Normalize(old) == Normalize(new)
}

func HashCode(location interface{}) int {
	loc := location.(string)
	return pluginsdk.HashString(Normalize(loc))
}

func StateFunc(location interface{}) string {
	input := location.(string)
	return Normalize(input)
}
