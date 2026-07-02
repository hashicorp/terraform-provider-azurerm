// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package aznr002

import (
	"context"
	"fmt"

	"testdata/src/mockpkg/pluginsdk"
	"testdata/src/mockpkg/sdk"
)

// Test Case: Write-only attribute handled via pluginsdk.GetWriteOnly
type WriteOnlyResourceModel struct {
	Name             string `tfschema:"name"`
	ConnectionString string `tfschema:"connection_string"`
	ConnectionWO     string `tfschema:"connection_string_wo"`
}

type WriteOnlyResource struct{}

var _ sdk.ResourceWithUpdate = WriteOnlyResource{}

func (r WriteOnlyResource) ResourceType() string {
	return "azurerm_write_only_resource"
}

func (r WriteOnlyResource) ModelObject() interface{} {
	return &WriteOnlyResourceModel{}
}

func (r WriteOnlyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"connection_string": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"connection_string_wo": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r WriteOnlyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WriteOnlyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r WriteOnlyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r WriteOnlyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r WriteOnlyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config WriteOnlyResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				// Update connection_string
			}

			// Write-only attribute accessed via pluginsdk.GetWriteOnly - should be detected
			woVal, err := pluginsdk.GetWriteOnly(metadata.ResourceData, "connection_string_wo", nil)
			if err != nil {
				return err
			}
			_ = woVal

			return nil
		},
	}
}
