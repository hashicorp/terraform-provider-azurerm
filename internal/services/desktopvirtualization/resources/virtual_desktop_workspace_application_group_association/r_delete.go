// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_workspace_application_group_association

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopWorkspaceApplicationGroupAssociationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceApplicationGroupAssociationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Workspace.WorkspaceName, workspaceResourceType)
	defer locks.UnlockByName(id.Workspace.WorkspaceName, workspaceResourceType)

	locks.ByName(id.ApplicationGroup.ApplicationGroupName, applicationGroupType)
	defer locks.UnlockByName(id.ApplicationGroup.ApplicationGroupName, applicationGroupType)

	existing, err := client.Get(ctx, id.Workspace)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("%s was not found", id.Workspace)
		}

		return fmt.Errorf("retrieving %s: %+v", id.Workspace, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id.Workspace)
	}
	model := *existing.Model

	applicationGroupReferences := []string{}
	applicationGroupId := id.ApplicationGroup.ID()
	if props := model.Properties; props != nil && props.ApplicationGroupReferences != nil {
		for _, referenceId := range *props.ApplicationGroupReferences {
			if strings.EqualFold(referenceId, applicationGroupId) {
				continue
			}

			applicationGroupReferences = append(applicationGroupReferences, referenceId)
		}
	}

	payload := workspace.WorkspacePatch{
		Properties: &workspace.WorkspacePatchProperties{
			ApplicationGroupReferences: &applicationGroupReferences,
		},
		Tags: model.Tags,
	}
	if _, err = client.Update(ctx, id.Workspace, payload); err != nil {
		return fmt.Errorf("removing association between %s and %s: %+v", id.Workspace, id.ApplicationGroup, err)
	}

	return nil
}
