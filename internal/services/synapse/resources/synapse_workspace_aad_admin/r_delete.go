// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_aad_admin

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceAADAdminDelete(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).WorkspaceAadAdminsClient
	workspaceClient := synapseClient.From(meta).WorkspaceClient
	azureADOnlyAuthenticationsClient := synapseClient.From(meta).WorkspaceAzureADOnlyAuthenticationsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.WorkspaceAADAdminID(d.Id())
	if err != nil {
		return err
	}

	workspace, err := workspaceClient.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		if response.WasNotFound(workspace.Response.Response) {
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
