// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	connectionsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/connections/client"
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

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

func (r DataSource) ResourceType() string {
	return "azurerm_api_connection"
}

func (r DataSource) Description() string {
	return "Use this data source to access information about an existing API Connection."
}

func (r DataSource) ModelObject() any {
	return &ApiConnectionDataSourceModel{}
}

func (r DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return connections.ValidateConnectionID
}

func (r DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the API Connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"managed_api_id": {
			Description: "The ID of the Managed API that this connection is linked to.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"display_name": {
			Description: "The display name of the API Connection.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"parameter_values": {
			Description: "A mapping of parameter names to their values for the API Connection.",
			Type:        pluginsdk.TypeMap,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := connectionsClient.From(metadata.Client).ConnectionsClient
			subscriptionId := account.From(metadata.Client).SubscriptionId

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
					state.DisplayName = pointer.From(props.DisplayName)

					if props.Api != nil {
						state.ManagedApiId = pointer.From(props.Api.Id)
					}

					// In version 2016-06-01 the API doesn't return `ParameterValues`.
					// The non-secret parameters are returned in `NonSecretParameterValues` instead.
					state.ParameterValues = flattenParameterValues(pointer.From(props.NonSecretParameterValues))
				}

				state.Tags = pointer.From(model.Tags)
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
