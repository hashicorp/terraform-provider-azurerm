// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FunctionAppFunctionKeysDataSource struct{}

var _ sdk.DataSource = FunctionAppFunctionKeysDataSource{}

type FunctionAppFunctionKeysDataSourceModel struct {
	FunctionAppID string            `tfschema:"function_app_id"`
	FunctionName  string            `tfschema:"function_name"`
	DefaultKey    string            `tfschema:"default_key"`
	Keys          map[string]string `tfschema:"keys"`
}

func (d FunctionAppFunctionKeysDataSource) ResourceType() string {
	return "azurerm_function_app_function_keys"
}

func (d FunctionAppFunctionKeysDataSource) ModelObject() any {
	return &FunctionAppFunctionKeysDataSourceModel{}
}

func (d FunctionAppFunctionKeysDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"function_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateFunctionAppID,
			Description:  "The ID of the Function App in which the Function resides.",
		},

		"function_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.FunctionAppFunctionName,
			Description:  "The name of the Function.",
		},
	}
}

func (d FunctionAppFunctionKeysDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"default_key": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "The value of the `default` Function key.",
		},

		"keys": {
			Type:      pluginsdk.TypeMap,
			Computed:  true,
			Sensitive: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A map of all Function keys, where the map key is the name of the Function key and the value is the key itself.",
		},
	}
}

func (d FunctionAppFunctionKeysDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			var state FunctionAppFunctionKeysDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			appId, err := commonids.ParseFunctionAppID(state.FunctionAppID)
			if err != nil {
				return err
			}

			id := webapps.NewFunctionID(appId.SubscriptionId, appId.ResourceGroupName, appId.SiteName, state.FunctionName)

			keys, err := listFunctionAppFunctionKeys(ctx, client, id)
			if err != nil {
				return err
			}

			state.Keys = keys
			state.DefaultKey = keys[functionAppFunctionDefaultKeyName]

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
