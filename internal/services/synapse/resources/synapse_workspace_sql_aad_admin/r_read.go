// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_sql_aad_admin

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceSqlAADAdminRead(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).WorkspaceSQLAadAdminsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.WorkspaceSqlAADAdminID(d.Id())
	if err != nil {
		return err
	}

	aadAdmin, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		if response.WasNotFound(aadAdmin.Response.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	workspaceID := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

	d.Set("synapse_workspace_id", workspaceID.ID())
	d.Set("login", aadAdmin.Login)
	d.Set("object_id", aadAdmin.Sid)
	d.Set("tenant_id", aadAdmin.TenantID)

	return nil
}
