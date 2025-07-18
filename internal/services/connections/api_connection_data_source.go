// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connections

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiConnectionDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	ManagedApiId      string            `tfschema:"managed_api_id"`
	DisplayName       string            `tfschema:"display_name"`
	ParameterValues   map[string]string `tfschema:"parameter_values"`
	Tags              map[string]string `tfschema:"tags"`
}

type ApiConnectionDataSource struct{}

var _ sdk.DataSource = ApiConnectionDataSource{}

func (r ApiConnectionDataSource) ResourceType() string {
	return "azurerm_api_connection"
}

func (r ApiConnectionDataSource) ModelObject() interface{} {
	return &ApiConnectionDataSourceModel{}
}

func (r ApiConnectionDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return connections.ValidateConnectionID
}

func (r ApiConnectionDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r ApiConnectionDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"managed_api_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"parameter_values": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r ApiConnectionDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Connections.ConnectionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ApiConnectionDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := connections.NewConnectionID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.Name = id.ConnectionName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Location = location.NormalizeNilable(model.Location)

				if props := model.Properties; props != nil {
					if props.DisplayName != nil {
						state.DisplayName = *props.DisplayName
					}

					if props.Api != nil && props.Api.Id != nil {
						state.ManagedApiId = *props.Api.Id
					}

					// In version 2016-06-01 the API doesn't return `ParameterValues`.
					// The non-secret parameters are returned in `NonSecretParameterValues` instead.
					state.ParameterValues = flattenParameterValues(pointer.From(props.NonSecretParameterValues))
				}

				if model.Tags != nil {
					state.Tags = *model.Tags
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
