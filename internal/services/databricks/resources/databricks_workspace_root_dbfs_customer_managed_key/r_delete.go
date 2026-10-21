// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace_root_dbfs_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func databricksWorkspaceRootDbfsCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.WorkspaceName, "azurerm_databricks_workspace")
	defer locks.UnlockByName(id.WorkspaceName, "azurerm_databricks_workspace")

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `Model` was nil", id)
	}

	if existing.Model.Properties.Parameters == nil {
		return fmt.Errorf("retrieving %s: `Parameters` was nil", id)
	}

	existing.Model.Properties.Parameters.Encryption = &workspaces.WorkspaceEncryptionParameter{
		Value: &workspaces.Encryption{
			KeySource: pointer.To(workspaces.KeySourceDefault),
		},
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("removing Root DBFS Customer Managed Key from %s: %+v", *id, err)
	}

	return nil
}
