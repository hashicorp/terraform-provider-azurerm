// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_workspace

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	desktopvirtualizationClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopWorkspaceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := desktopvirtualizationClient.From(meta).WorkspacesClient

	id, err := workspace.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.WorkspaceName, workspaceResourceType)
	defer locks.UnlockByName(id.WorkspaceName, workspaceResourceType)

	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()
	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
