// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/resourceanchors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = ResourceAnchorDataSource{}

type ResourceAnchorDataSource struct{}

type ResourceAnchorDataSourceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`

	Location            string            `tfschema:"location"`
	Tags                map[string]string `tfschema:"tags"`
	LinkedCompartmentID string            `tfschema:"linked_compartment_id"`
}

func (ResourceAnchorDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (ResourceAnchorDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// Common fields
		"location": commonschema.LocationComputed(),
		"tags":     commonschema.TagsDataSource(),

		"linked_compartment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (ResourceAnchorDataSource) ModelObject() interface{} {
	return &ResourceAnchorDataSourceModel{}
}

func (s ResourceAnchorDataSource) ResourceType() string {
	return "azurerm_oracle_resource_anchor"
}

func (ResourceAnchorDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient09.ResourceAnchors
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ResourceAnchorDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := resourceanchors.NewResourceAnchorID(subscriptionId, model.ResourceGroupName, model.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("resource Anchor %s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ResourceAnchorDataSourceModel{
				Name:              id.ResourceAnchorName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = model.Location
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.LinkedCompartmentID = pointer.From(props.LinkedCompartmentId)
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
