// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/jobdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StorageMoverJobDefinitionResourceModel struct {
	Name                  string                  `tfschema:"name"`
	StorageMoverProjectId string                  `tfschema:"storage_mover_project_id"`
	SourceName            string                  `tfschema:"source_name"`
	TargetName            string                  `tfschema:"target_name"`
	CopyMode              jobdefinitions.CopyMode `tfschema:"copy_mode"`
	SourceSubpath         string                  `tfschema:"source_sub_path"`
	TargetSubpath         string                  `tfschema:"target_sub_path"`
	AgentName             string                  `tfschema:"agent_name"`
	Description           string                  `tfschema:"description"`
}

type StorageMoverJobDefinitionResource struct{}

var _ sdk.ResourceWithUpdate = StorageMoverJobDefinitionResource{}

func (r StorageMoverJobDefinitionResource) ResourceType() string {
	return "azurerm_storage_mover_job_definition"
}

func (r StorageMoverJobDefinitionResource) ModelObject() interface{} {
	return &StorageMoverJobDefinitionResourceModel{}
}

func (r StorageMoverJobDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return jobdefinitions.ValidateJobDefinitionID
}

func (r StorageMoverJobDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9a-zA-Z][-_0-9a-zA-Z]{0,63}$`),
				`The name must be between 1 and 64 characters in length, begin with a letter or number, and may contain letters, numbers, dashes and underscore.`,
			),
		},

		"storage_mover_project_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: projects.ValidateProjectID,
		},

		"source_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"copy_mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(jobdefinitions.CopyModeMirror),
				string(jobdefinitions.CopyModeAdditive),
			}, false),
		},

		"source_sub_path": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_sub_path": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"agent_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r StorageMoverJobDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StorageMoverJobDefinitionResource) Create() sdk.ResourceFunc {
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
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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
			return nil
		},
	}
}

func (r StorageMoverJobDefinitionResource) Update() sdk.ResourceFunc {
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

func (r StorageMoverJobDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.JobDefinitionsClient

			id, err := jobdefinitions.ParseJobDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := StorageMoverJobDefinitionResourceModel{
				Name:                  id.JobDefinitionName,
				StorageMoverProjectId: projects.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName, id.ProjectName).ID(),
			}

			if v := resp.Model; v != nil {
				state.AgentName = pointer.From(v.Properties.AgentName)

				state.CopyMode = v.Properties.CopyMode

				state.Description = pointer.From(v.Properties.Description)

				state.SourceName = v.Properties.SourceName

				state.SourceSubpath = pointer.From(v.Properties.SourceSubpath)

				state.TargetName = v.Properties.TargetName

				state.TargetSubpath = pointer.From(v.Properties.TargetSubpath)
			}
			return metadata.Encode(&state)
		},
	}
}

func (r StorageMoverJobDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.JobDefinitionsClient

			id, err := jobdefinitions.ParseJobDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
