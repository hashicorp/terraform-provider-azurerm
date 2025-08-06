// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSynapseWorkspaceSqlAADAdmin() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseWorkspaceSqlAADAdminCreateUpdate,
		Read:   resourceSynapseWorkspaceSqlAADAdminRead,
		Update: resourceSynapseWorkspaceSqlAADAdminCreateUpdate,
		Delete: resourceSynapseWorkspaceSqlAADAdminDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WorkspaceSqlAADAdminID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"login": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"object_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceSynapseWorkspaceSqlAADAdminCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceSQLAadAdminsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}
	workspaceName := workspaceId.Name
	workspaceResourceGroup := workspaceId.ResourceGroup

	aadAdmin := &synapse.WorkspaceAadAdminInfo{
		AadAdminProperties: &synapse.AadAdminProperties{
			TenantID:          utils.String(d.Get("tenant_id").(string)),
			Login:             utils.String(d.Get("login").(string)),
			AdministratorType: utils.String("ActiveDirectory"),
			Sid:               utils.String(d.Get("object_id").(string)),
		},
	}

	workspaceAadAdminsCreateOrUpdateFuture, err := client.CreateOrUpdate(ctx, workspaceResourceGroup, workspaceName, *aadAdmin)
	if err != nil {
		return fmt.Errorf("updating Synapse Workspace %q Sql AAD Admin (Resource Group %q): %+v", workspaceName, workspaceResourceGroup, err)
	}

	if err = workspaceAadAdminsCreateOrUpdateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on updating for Synapse Workspace %q Sql AAD Admin (Resource Group %q): %+v", workspaceName, workspaceResourceGroup, err)
	}

	id := parse.NewWorkspaceSqlAADAdminID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, "activeDirectory")
	d.SetId(id.ID())

	return resourceSynapseWorkspaceSqlAADAdminRead(d, meta)
}

func resourceSynapseWorkspaceSqlAADAdminRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceSQLAadAdminsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceSqlAADAdminID(d.Id())
	if err != nil {
		return err
	}

	aadAdmin, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		if utils.ResponseWasNotFound(aadAdmin.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	workspaceID := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

	d.Set("synapse_workspace_id", workspaceID.ID())
	d.Set("login", aadAdmin.AadAdminProperties.Login)
	d.Set("object_id", aadAdmin.AadAdminProperties.Sid)
	d.Set("tenant_id", aadAdmin.AadAdminProperties.TenantID)

	return nil
}

func resourceSynapseWorkspaceSqlAADAdminDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceSQLAadAdminsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
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
