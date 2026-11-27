// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_job_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/jobdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model StorageMoverJobDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.StorageMover.JobDefinitionsClient
			projectId, err := projects.ParseProjectID(model.StorageMoverProjectId)
			if err != nil {
				return err
			}

			id := jobdefinitions.NewJobDefinitionID(projectId.SubscriptionId, projectId.ResourceGroupName, projectId.StorageMoverName, projectId.ProjectName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := jobdefinitions.JobDefinition{
				Properties: jobdefinitions.JobDefinitionProperties{
					CopyMode:   model.CopyMode,
					SourceName: model.SourceName,
					TargetName: model.TargetName,
				},
			}

			if model.AgentName != "" {
				properties.Properties.AgentName = &model.AgentName
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if model.TargetSubpath != "" {
				properties.Properties.TargetSubpath = &model.TargetSubpath
			}

			if model.SourceSubpath != "" {
				properties.Properties.SourceSubpath = &model.SourceSubpath
			}

			if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}
