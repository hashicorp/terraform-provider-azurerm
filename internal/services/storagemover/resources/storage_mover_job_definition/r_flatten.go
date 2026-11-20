// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_job_definition

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/jobdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *jobdefinitions.JobDefinitionId, model *jobdefinitions.JobDefinition) error {
	state := StorageMoverJobDefinitionResourceModel{
		Name:                  id.JobDefinitionName,
		StorageMoverProjectId: projects.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName, id.ProjectName).ID(),
	}

	if model != nil {
		state.AgentName = pointer.From(model.Properties.AgentName)
		state.CopyMode = model.Properties.CopyMode
		state.Description = pointer.From(model.Properties.Description)
		state.SourceName = model.Properties.SourceName
		state.SourceSubpath = pointer.From(model.Properties.SourceSubpath)
		state.TargetName = model.Properties.TargetName
		state.TargetSubpath = pointer.From(model.Properties.TargetSubpath)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
