package aznr002

import (
	"context"
	"fmt"

	"testdata/src/mockpkg/pluginsdk"
	"testdata/src/mockpkg/sdk"
)

// Test Case 5: Using HasChanges (plural) to check multiple properties at once
type HasChangesResourceModel struct {
	Name        string `tfschema:"name"`
	Location    string `tfschema:"location"`
	DisplayName string `tfschema:"display_name"`
	Description string `tfschema:"description"`
}

type HasChangesResource struct{}

var _ sdk.ResourceWithUpdate = HasChangesResource{}

func (r HasChangesResource) ResourceType() string {
	return "azurerm_has_changes_resource"
}

func (r HasChangesResource) ModelObject() interface{} {
	return &HasChangesResourceModel{}
}

func (r HasChangesResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r HasChangesResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r HasChangesResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r HasChangesResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r HasChangesResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r HasChangesResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config HasChangesResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Using HasChanges to check multiple properties at once
			if metadata.ResourceData.HasChanges("location", "display_name", "description") {
				// Update all properties
			}

			return nil
		},
	}
}
