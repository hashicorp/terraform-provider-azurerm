// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_project

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Projects
			schema := DevCenterProjectResourceSchema{}

			id, err := projects.ParseProjectID(metadata.ResourceData.Id())
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
				schema.Name = id.ProjectName
				schema.ResourceGroupName = id.ResourceGroupName
				schema.Location = location.Normalize(model.Location)
				schema.Tags = tags.Flatten(model.Tags)

				identity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %v", err)
				}
				schema.Identity = pointer.From(identity)

				if props := model.Properties; props != nil {
					schema.Description = pointer.From(props.Description)
					schema.DevCenterUri = pointer.From(props.DevCenterUri)
					schema.MaximumDevBoxesPerUser = pointer.From(props.MaxDevBoxesPerUser)
					if devCenterId := pointer.From(props.DevCenterId); devCenterId != "" {
						parsedDevCenterId, err := devcenters.ParseDevCenterIDInsensitively(devCenterId)
						if err != nil {
							return err
						}
						schema.DevCenterId = parsedDevCenterId.ID()
					}
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

type DevCenterProjectResourceSchema struct {
	Description            string                                     `tfschema:"description"`
	DevCenterId            string                                     `tfschema:"dev_center_id"`
	DevCenterUri           string                                     `tfschema:"dev_center_uri"`
	Identity               []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location               string                                     `tfschema:"location"`
	MaximumDevBoxesPerUser int64                                      `tfschema:"maximum_dev_boxes_per_user"`
	Name                   string                                     `tfschema:"name"`
	ResourceGroupName      string                                     `tfschema:"resource_group_name"`
	Tags                   map[string]any                             `tfschema:"tags"`
}
