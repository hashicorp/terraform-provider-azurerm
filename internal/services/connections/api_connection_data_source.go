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
	Name               string                   `tfschema:"name"`
	ResourceGroupName  string                   `tfschema:"resource_group_name"`
	Location           string                   `tfschema:"location"`
	ManagedApiId       string                   `tfschema:"managed_api_id"`
	DisplayName        string                   `tfschema:"display_name"`
	Kind               string                   `tfschema:"kind"`
	ParameterValues    map[string]string        `tfschema:"parameter_values"`
	ParameterValueType string                   `tfschema:"parameter_value_type"`
	ParameterValueSet  []ParameterValueSetModel `tfschema:"parameter_value_set"`
	Tags               map[string]string        `tfschema:"tags"`
}

type ParameterValueSetModel struct {
	Name   string            `tfschema:"name"`
	Values map[string]string `tfschema:"values"`
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

		"kind": {
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

		"parameter_value_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"parameter_value_set": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"values": {
						Type:     pluginsdk.TypeMap,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
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
				state.Kind = pointer.From(model.Kind)

				if props := model.Properties; props != nil {
					state.DisplayName = pointer.From(props.DisplayName)

					if props.Api != nil {
						state.ManagedApiId = pointer.From(props.Api.Id)
					}

					// In version 2016-06-01 the API doesn't return `ParameterValues`.
					// The non-secret parameters are returned in `NonSecretParameterValues` instead.
					state.ParameterValues = flattenParameterValues(pointer.From(props.NonSecretParameterValues))
					state.ParameterValueType = pointer.From(props.ParameterValueType)
					state.ParameterValueSet = flattenParameterValueSetForDataSource(props.ParameterValueSet)
				}

				state.Tags = pointer.From(model.Tags)
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}

func flattenParameterValueSetForDataSource(input *connections.ParameterValueSet) []ParameterValueSetModel {
	if input == nil {
		return []ParameterValueSetModel{}
	}

	return []ParameterValueSetModel{
		{
			Name:   pointer.From(input.Name),
			Values: flattenParameterValueSetValues(input.Values),
		},
	}
}
