// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_job_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/jobdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.JobDefinitionsClient

			id, err := jobdefinitions.ParseJobDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StorageMoverJobDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("agent_name") {
				if model.AgentName != "" {
					properties.Properties.AgentName = &model.AgentName
				} else {
					properties.Properties.AgentName = nil
				}
			}

			if metadata.ResourceData.HasChange("copy_mode") {
				properties.Properties.CopyMode = model.CopyMode
			}

			if metadata.ResourceData.HasChange("agent_name") {
				properties.Properties.AgentName = pointer.To(model.AgentName)
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Properties.Description = pointer.To(model.Description)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
