package aznr002

import (
	"context"
	"fmt"

	"testdata/src/mockpkg/pluginsdk"
	"testdata/src/mockpkg/sdk"
)

// Test Case 2: Missing HasChange for some updatable properties - should report error
type BadResourceModel struct {
	Name        string `tfschema:"name"`
	Location    string `tfschema:"location"`
	DisplayName string `tfschema:"display_name"`
	Description string `tfschema:"description"`
	Enabled     bool   `tfschema:"enabled"`
}

type BadResource struct{}

var _ sdk.ResourceWithUpdate = BadResource{}

func (r BadResource) ResourceType() string {
	return "azurerm_bad_resource"
}

func (r BadResource) ModelObject() interface{} {
	return &BadResourceModel{}
}

func (r BadResource) Arguments() map[string]*pluginsdk.Schema {
	sche := map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"description": descriptionSchema(),
		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
	return sche
}

func (r BadResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r BadResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r BadResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r BadResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r BadResource) Update() sdk.ResourceFunc { // want "AZNR002: resource has updatable properties not handled in Update function: `description, enabled`. If they are non-updatable, mark them as ForceNew: true in Arguments\\(\\) schema"
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config BadResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Only checking some properties - missing description and enabled
			if metadata.ResourceData.HasChange("location") {
				// Update location
			}

			if metadata.ResourceData.HasChange("display_name") {
				// Update display_name
			}

			// Missing: description, enabled

			return nil
		},
	}
}

func descriptionSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
	}
}
