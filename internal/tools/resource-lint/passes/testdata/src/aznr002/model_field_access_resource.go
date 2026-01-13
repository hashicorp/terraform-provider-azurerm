package aznr002

import (
	"context"
	"fmt"

	"testdata/src/mockpkg/pluginsdk"
	"testdata/src/mockpkg/sdk"
)

// Test Case 4: Using model field access (Pattern 4)
type ModelFieldAccessResourceModel struct {
	Name        string            `tfschema:"name"`
	DisplayName string            `tfschema:"display_name"`
	Description string            `tfschema:"description"`
	Tags        map[string]string `tfschema:"tags"`
}

type ModelFieldAccessResource struct{}

var _ sdk.ResourceWithUpdate = ModelFieldAccessResource{}

func (r ModelFieldAccessResource) ResourceType() string {
	return "azurerm_model_field_access_resource"
}

func (r ModelFieldAccessResource) ModelObject() interface{} {
	return &ModelFieldAccessResourceModel{}
}

func (r ModelFieldAccessResource) Arguments() map[string]*pluginsdk.Schema {
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
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r ModelFieldAccessResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ModelFieldAccessResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r ModelFieldAccessResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r ModelFieldAccessResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r ModelFieldAccessResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config ModelFieldAccessResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Using model field access - should be detected by Pattern 4
			displayName := config.DisplayName
			description := config.Description
			tags := config.Tags

			// Do something with the fields
			_ = displayName
			_ = description
			_ = tags

			return nil
		},
	}
}
