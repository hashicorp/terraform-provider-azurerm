package aznr002

import (
	"context"
	"fmt"

	"testdata/src/mockpkg/pluginsdk"
	"testdata/src/mockpkg/sdk"
)

// Test Case 6: Resource with ForceNew properties (should not be reported)
type ForceNewResourceModel struct {
	Name        string `tfschema:"name"`
	Location    string `tfschema:"location"`
	DisplayName string `tfschema:"display_name"`
}

type ForceNewResource struct{}

var _ sdk.ResourceWithUpdate = ForceNewResource{}

func (r ForceNewResource) ResourceType() string {
	return "azurerm_force_new_resource"
}

func (r ForceNewResource) ModelObject() interface{} {
	return &ForceNewResourceModel{}
}

func (r ForceNewResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true, // ForceNew - not updatable
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r ForceNewResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ForceNewResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r ForceNewResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r ForceNewResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r ForceNewResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config ForceNewResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Only need to check display_name since location is ForceNew
			if metadata.ResourceData.HasChange("display_name") {
				// Update display_name
			}

			// No need to check location - it's ForceNew

			return nil
		},
	}
}
