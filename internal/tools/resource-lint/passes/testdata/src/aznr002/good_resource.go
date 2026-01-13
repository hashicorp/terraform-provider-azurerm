package aznr002

import (
	"context"
	"fmt"

	"testdata/src/mockpkg/pluginsdk"
	"testdata/src/mockpkg/sdk"
)

// Test Case 1: All updatable properties properly handled with HasChange
type GoodResourceModel struct {
	Name        string `tfschema:"name"`
	Location    string `tfschema:"location"`
	DisplayName string `tfschema:"display_name"`
	Description string `tfschema:"description"`
	Enabled     bool   `tfschema:"enabled"`
}

type GoodResource struct{}

var _ sdk.ResourceWithUpdate = GoodResource{}

func (r GoodResource) ResourceType() string {
	return "azurerm_good_resource"
}

func (r GoodResource) ModelObject() interface{} {
	return &GoodResourceModel{}
}

func (r GoodResource) Arguments() map[string]*pluginsdk.Schema {
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
		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}

func (r GoodResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r GoodResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r GoodResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r GoodResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r GoodResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config GoodResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// All updatable properties have HasChange checks - this is good!
			if metadata.ResourceData.HasChange("location") {
				// Update location
			}

			if metadata.ResourceData.HasChange("display_name") {
				// Update display_name
			}

			if metadata.ResourceData.HasChange("description") {
				// Update description
			}

			if metadata.ResourceData.HasChange("enabled") {
				// Update enabled
			}

			return nil
		},
	}
}
