// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_aad_admin

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceAADAdminCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.WorkspaceAadAdminsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	workspaceId, err := workspaces.ParseWorkspaceIDInsensitively(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}
	workspaceName := workspaceId.WorkspaceName
	workspaceResourceGroup := workspaceId.ResourceGroupName

	aadAdmin := &synapse.WorkspaceAadAdminInfo{
		AadAdminProperties: &synapse.AadAdminProperties{
			TenantID:          pointer.To(d.Get("tenant_id").(string)),
			Login:             pointer.To(d.Get("login").(string)),
			AdministratorType: pointer.To("ActiveDirectory"),
			Sid:               pointer.To(d.Get("object_id").(string)),
		},
	}

	workspaceAadAdminsCreateOrUpdateFuture, err := client.CreateOrUpdate(ctx, workspaceResourceGroup, workspaceName, *aadAdmin)
	if err != nil {
		return fmt.Errorf("updating Synapse Workspace %q AAD Admin (Resource Group %q): %+v", workspaceName, workspaceResourceGroup, err)
	}

	id := parse.NewWorkspaceAADAdminID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, "activeDirectory")
	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	if err = workspaceAadAdminsCreateOrUpdateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on updating for Synapse Workspace %q AAD Admin (Resource Group %q): %+v", workspaceName, workspaceResourceGroup, err)
	}

	return resourceSynapseWorkspaceAADAdminRead(d, meta)
}
