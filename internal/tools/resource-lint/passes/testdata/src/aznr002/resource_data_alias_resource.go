package aznr002

import (
	"context"
	"fmt"

	"testdata/src/mockpkg/pluginsdk"
	"testdata/src/mockpkg/sdk"
)

// Test Case 7: Using ResourceData alias (rd := metadata.ResourceData)
type ResourceDataAliasModel struct {
	Name        string `tfschema:"name"`
	DisplayName string `tfschema:"display_name"`
	Description string `tfschema:"description"`
}

type ResourceDataAliasResource struct{}

var _ sdk.ResourceWithUpdate = ResourceDataAliasResource{}

func (r ResourceDataAliasResource) ResourceType() string {
	return "azurerm_resource_data_alias_resource"
}

func (r ResourceDataAliasResource) ModelObject() interface{} {
	return &ResourceDataAliasModel{}
}

func (r ResourceDataAliasResource) Arguments() map[string]*pluginsdk.Schema {
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
	}
}

func (r ResourceDataAliasResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ResourceDataAliasResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r ResourceDataAliasResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r ResourceDataAliasResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r ResourceDataAliasResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config ResourceDataAliasModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Using alias for ResourceData
			rd := metadata.ResourceData

			if rd.HasChange("display_name") {
				// Update display_name
			}

			if rd.HasChange("description") {
				// Update description
			}

			return nil
		},
	}
}
