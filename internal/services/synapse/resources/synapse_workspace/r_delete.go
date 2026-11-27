// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroupName, id.WorkspaceName)
	if err != nil {
		return fmt.Errorf("deleting Synapse Workspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
	}

	// sometimes the waitForCompletion rest api will return 404
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for Synapse Workspace %q (Resource Group %q) to be deleted: %+v", id.WorkspaceName, id.ResourceGroupName, err)
		}
	}

	return nil
}
