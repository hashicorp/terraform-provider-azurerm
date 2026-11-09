// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_workspace_application_group_association

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopWorkspaceApplicationGroupAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceApplicationGroupAssociationID(d.Id())
	if err != nil {
		return err
	}

	workspace, err := client.Get(ctx, id.Workspace)
	if err != nil {
		if response.WasNotFound(workspace.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id.Workspace)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id.Workspace, err)
	}
	if model := workspace.Model; model != nil {
		applicationGroupId := id.ApplicationGroup.ID()
		exists := associationExists(model.Properties, applicationGroupId)
		if !exists {
			log.Printf("[DEBUG] Association between %s and %s was not found - removing from state!", id.Workspace, id.ApplicationGroup)
			d.SetId("")
			return nil
		}

		d.Set("workspace_id", id.Workspace.ID())
		d.Set("application_group_id", applicationGroupId)
	}

	return nil
}

func associationExists(props *workspace.WorkspaceProperties, applicationGroupId string) bool {
	if props == nil || props.ApplicationGroupReferences == nil {
		return false
	}

	for _, id := range *props.ApplicationGroupReferences {
		if strings.EqualFold(id, applicationGroupId) {
			return true
		}
	}

	return false
}
