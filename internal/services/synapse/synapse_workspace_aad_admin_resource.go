// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSynapseWorkspaceAADAdmin() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseWorkspaceAADAdminCreateUpdate,
		Read:   resourceSynapseWorkspaceAADAdminRead,
		Update: resourceSynapseWorkspaceAADAdminCreateUpdate,
		Delete: resourceSynapseWorkspaceAADAdminDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WorkspaceAADAdminID(id)
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

func resourceSynapseWorkspaceAADAdminCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceAadAdminsClient
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
		return fmt.Errorf("updating Synapse Workspace %q AAD Admin (Resource Group %q): %+v", workspaceName, workspaceResourceGroup, err)
	}

	if err = workspaceAadAdminsCreateOrUpdateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on updating for Synapse Workspace %q AAD Admin (Resource Group %q): %+v", workspaceName, workspaceResourceGroup, err)
	}

	id := parse.NewWorkspaceAADAdminID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, "activeDirectory")
	d.SetId(id.ID())

	return resourceSynapseWorkspaceAADAdminRead(d, meta)
}

func resourceSynapseWorkspaceAADAdminRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceAadAdminsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceAADAdminID(d.Id())
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

func resourceSynapseWorkspaceAADAdminDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceAadAdminsClient
	workspaceClient := meta.(*clients.Client).Synapse.WorkspaceClient
	azureADOnlyAuthenticationsClient := meta.(*clients.Client).Synapse.WorkspaceAzureADOnlyAuthenticationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceAADAdminID(d.Id())
	if err != nil {
		return err
	}

	workspace, err := workspaceClient.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		if utils.ResponseWasNotFound(workspace.Response) {
			return fmt.Errorf("retrieving %q: %+v", id, err)
		}
	}

	if pointer.From(workspace.AzureADOnlyAuthentication) {
		// Remove Azure Active Azure Admin, API requires AzureADOnlyAuthentication to be disabled first, see below error message from API.
		// "User tried to delete managed server Azure Active Azure admin when AzureADOnlyAuthentication is set, please use azureADOnlyAuthentications API first."
		future, err := azureADOnlyAuthenticationsClient.Create(ctx, id.ResourceGroup, id.WorkspaceName, synapse.AzureADOnlyAuthentication{
			AzureADOnlyAuthenticationProperties: &synapse.AzureADOnlyAuthenticationProperties{
				AzureADOnlyAuthentication: pointer.To(false),
			},
		})
		if err != nil {
			return fmt.Errorf("updating azuread_authentication_only for %s: %+v", *workspace.ID, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for azuread_authentication_only to finish updating for %s: %+v", *workspace.ID, err)
		}
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		return fmt.Errorf("setting empty Synapse Workspace %q AAD Admin (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on setting empty Synapse Workspace %q AAD Admin (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroup, err)
	}

	return nil
}
