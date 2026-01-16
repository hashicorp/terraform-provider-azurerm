package aznr002

import (
	"context"
	"fmt"

	"testdata/src/mockpkg/pluginsdk"
	"testdata/src/mockpkg/sdk"
)

// Test Case: Helper function with model argument should be traced
// Properties accessed in helper should be considered handled
type HelperTracedResourceModel struct {
	Name        string `tfschema:"name"`
	DisplayName string `tfschema:"display_name"`
	Description string `tfschema:"description"`
	Enabled     bool   `tfschema:"enabled"`
	Tags        string `tfschema:"tags"`
}

type HelperTracedResource struct{}

var _ sdk.ResourceWithUpdate = HelperTracedResource{}

func (r HelperTracedResource) ResourceType() string {
	return "azurerm_helper_traced_resource"
}

func (r HelperTracedResource) ModelObject() interface{} {
	return &HelperTracedResourceModel{}
}

func (r HelperTracedResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"tags": { // want `AZNR002`
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r HelperTracedResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r HelperTracedResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r HelperTracedResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r HelperTracedResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r HelperTracedResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config HelperTracedResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Call helper at top level - should trace into it
			props := expandHelperTracedProperties(config)
			_ = props

			// Tags is NOT handled in helper - should report error
			return nil
		},
	}
}

// Helper function that accesses model fields
// These accesses should be detected by the tracer
func expandHelperTracedProperties(model HelperTracedResourceModel) map[string]interface{} {
	result := make(map[string]interface{})

	// Access model fields - these should be detected as handled
	result["displayName"] = model.DisplayName
	result["description"] = model.Description
	result["enabled"] = model.Enabled

	// Note: model.Tags is NOT accessed - should trigger AZNR002 error

	return result
}
