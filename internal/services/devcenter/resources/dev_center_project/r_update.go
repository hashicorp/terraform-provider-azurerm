// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_project

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Projects

			id, err := projects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config DevCenterProjectResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			var payload projects.ProjectUpdate

			if metadata.ResourceData.HasChanges("tags") {
				payload.Tags = tags.Expand(config.Tags)
			}

			if metadata.ResourceData.HasChange("identity") {
				identity, err := identity.ExpandSystemAndUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return err
				}
				payload.Identity = identity
			}

			if payload.Properties == nil {
				payload.Properties = &projects.ProjectUpdateProperties{}
			}

			if metadata.ResourceData.HasChanges("maximum_dev_boxes_per_user") {
				payload.Properties.MaxDevBoxesPerUser = pointer.To(config.MaximumDevBoxesPerUser)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
