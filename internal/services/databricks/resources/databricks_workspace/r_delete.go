// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabricksWorkspaceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	deleteOperationOptions := workspaces.DefaultDeleteOperationOptions()
	if meta.(*clients.Client).Features.DatabricksWorkspace.ForceDelete {
		deleteOperationOptions.ForceDeletion = pointer.To(true)
	}

	if err = client.DeleteThenPoll(ctx, *id, deleteOperationOptions); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
