// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_sql_aad_admin

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceSqlAADAdminDelete(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).WorkspaceSQLAadAdminsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.WorkspaceSqlAADAdminID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		return fmt.Errorf("setting empty Synapse Workspace %q Sql AAD Admin (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on setting empty Synapse Workspace %q Sql AAD Admin (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroup, err)
	}

	return nil
}
