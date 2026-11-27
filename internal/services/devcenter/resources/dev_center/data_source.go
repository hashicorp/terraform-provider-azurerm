// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dev_center

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	devcenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = DataSource{}

type DataSource struct{}

type DevCenterDataSourceModel struct {
	DevCenterUri      string                                     `tfschema:"dev_center_uri"`
	Identity          []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location          string                                     `tfschema:"location"`
	Name              string                                     `tfschema:"name"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	Tags              map[string]any                             `tfschema:"tags"`
}

func (DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Dev Center.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dev_center_uri": {
			Description: "The URI of the Dev Center.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"location": commonschema.LocationComputed(),

		"tags": commonschema.TagsDataSource(),
	}
}

func (DataSource) ModelObject() any {
	return &DevCenterDataSourceModel{}
}

func (DataSource) ResourceType() string {
	return "azurerm_dev_center"
}

func (DataSource) Description() string {
	return "Use this data source to access information about an existing Dev Center."
}

func (DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := devcenterClient.From(metadata.Client).V20250201.DevCenters
			subscriptionId := account.From(metadata.Client).SubscriptionId

			var state DevCenterDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := devcenters.NewDevCenterID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.Name = id.DevCenterName
				state.ResourceGroupName = id.ResourceGroupName

				if err := mapDevCenterToDevCenterResourceSchema(*model, &state); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func mapDevCenterToDevCenterResourceSchema(input devcenters.DevCenter, output *DevCenterDataSourceModel) error {
	identity, err := identity.FlattenSystemAndUserAssignedMapToModel(input.Identity)
	if err != nil {
		return fmt.Errorf("flattening SystemAndUserAssigned Identity: %+v", err)
	}
	output.Identity = *identity

	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &devcenters.DevCenterProperties{}
	}

	output.DevCenterUri = pointer.From(input.Properties.DevCenterUri)

	return nil
}
