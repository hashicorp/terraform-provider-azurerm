// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package playwright_workspace

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2025-09-01/playwrightworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) flatten(metadata sdk.ResourceMetaData, id *playwrightworkspaces.PlaywrightWorkspaceId, model *playwrightworkspaces.PlaywrightWorkspace) error {
	state := PlaywrightWorkspaceModel{
		Name:              id.PlaywrightWorkspaceName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)

		if properties := model.Properties; properties != nil {
			state.DataplaneUri = pointer.From(properties.DataplaneUri)
			state.Uuid = pointer.From(properties.WorkspaceId)
		}

		state.Tags = pointer.From(model.Tags)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
