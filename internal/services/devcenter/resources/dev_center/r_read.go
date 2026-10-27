// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevCenters
			schema := DevCenterResourceSchema{}

			id, err := devcenters.ParseDevCenterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.DevCenterName
				schema.ResourceGroupName = id.ResourceGroupName
				if err := r.mapDevCenterToDevCenterResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

type DevCenterResourceSchema struct {
	DevCenterUri                  string                                     `tfschema:"dev_center_uri"`
	Identity                      []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location                      string                                     `tfschema:"location"`
	Name                          string                                     `tfschema:"name"`
	ResourceGroupName             string                                     `tfschema:"resource_group_name"`
	ProjectCatalogItemSyncEnabled bool                                       `tfschema:"project_catalog_item_sync_enabled"`
	Tags                          map[string]any                             `tfschema:"tags"`
}

func (r Resource) mapDevCenterToDevCenterResourceSchema(input devcenters.DevCenter, output *DevCenterResourceSchema) error {
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

	if v := input.Properties.ProjectCatalogSettings; v != nil {
		output.ProjectCatalogItemSyncEnabled = pointer.From(v.CatalogItemSyncEnableStatus) == devcenters.CatalogItemSyncEnableStatusEnabled
	}

	return nil
}
