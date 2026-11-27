// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).WorkspaceClient
	azureADOnlyAuthenticationsClient := synapseClient.From(meta).WorkspaceAzureADOnlyAuthenticationsClient
	identitySQLControlClient := synapseClient.From(meta).WorkspaceManagedIdentitySQLControlSettingsClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	// lintignore:R019 // deliberate subset: only the fields feeding WorkspacePatchInfo; aad_admin, sql_aad_admin and identity settings are updated via separate API calls below
	if d.HasChanges(
		"tags",
		"sql_administrator_login_password",
		"github_repo",
		"azure_devops_repo",
		"customer_managed_key",
		"public_network_access_enabled",
	) {
		publicNetworkAccess := synapse.WorkspacePublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = synapse.WorkspacePublicNetworkAccessDisabled
		}
		workspacePatchInfo := synapse.WorkspacePatchInfo{
			Tags: tags.Expand(d.Get("tags").(map[string]any)),
			WorkspacePatchProperties: &synapse.WorkspacePatchProperties{
				SQLAdministratorLoginPassword:    pointer.To(d.Get("sql_administrator_login_password").(string)),
				WorkspaceRepositoryConfiguration: expandWorkspaceRepositoryConfiguration(d),
				Encryption:                       expandEncryptionDetails(d),
				PublicNetworkAccess:              publicNetworkAccess,
			},
		}

		if allowedLinkingTenantIds, ok := d.GetOk("linking_allowed_for_aad_tenant_ids"); ok {
			if workspacePatchInfo.ManagedVirtualNetworkSettings == nil {
				workspacePatchInfo.ManagedVirtualNetworkSettings = &synapse.ManagedVirtualNetworkSettings{}
			}
			workspacePatchInfo.ManagedVirtualNetworkSettings.AllowedAadTenantIdsForLinking = pluginsdk.ExpandStringSlice(allowedLinkingTenantIds.([]any))
		}

		if purviewId, ok := d.GetOk("purview_id"); ok {
			workspacePatchInfo.PurviewConfiguration = &synapse.PurviewConfiguration{
				PurviewResourceID: pointer.To(purviewId.(string)),
			}
		}

		if err := waitSynapseWorkspaceProvisioningState(ctx, client, id); err != nil {
			return fmt.Errorf("failed waiting for updating %s: %+v", id, err)
		}

		future, err := client.Update(ctx, id.ResourceGroupName, id.WorkspaceName, workspacePatchInfo)
		if err != nil {
			return fmt.Errorf("updating Synapse Workspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on updating future for Synapse Workspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
		}

		if err := waitSynapseWorkspaceCMKState(ctx, client, id); err != nil {
			return fmt.Errorf("failed waiting for updating %s: %+v", id, err)
		}
	}

	if d.HasChange("azuread_authentication_only") {
		future, err := azureADOnlyAuthenticationsClient.Create(ctx, id.ResourceGroupName, id.WorkspaceName, synapse.AzureADOnlyAuthentication{
			AzureADOnlyAuthenticationProperties: &synapse.AzureADOnlyAuthenticationProperties{
				AzureADOnlyAuthentication: pointer.To(d.Get("azuread_authentication_only").(bool)),
			},
		})
		if err != nil {
			return fmt.Errorf("updating azuread_authentication_only for %s: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for azuread_authentication_only to finish updating for %s: %+v", id, err)
		}
	}

	if d.HasChange("sql_identity_control_enabled") {
		sqlControlSettings := expandIdentityControlSQLSettings(d.Get("sql_identity_control_enabled").(bool))
		if err := waitSynapseWorkspaceProvisioningState(ctx, client, id); err != nil {
			return fmt.Errorf("failed waiting for updating %s: %+v", id, err)
		}
		future, err := identitySQLControlClient.CreateOrUpdate(ctx, id.ResourceGroupName, id.WorkspaceName, *sqlControlSettings)
		if err != nil {
			return fmt.Errorf("updating workspace identity control for SQL pool: %+v", err)
		}
		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update workspace identity control for SQL pool of %q: %+v", id, err)
		}
	}

	if err := waitSynapseWorkspaceProvisioningState(ctx, client, id); err != nil {
		return fmt.Errorf("failed waiting for updating %s: %+v", id, err)
	}

	return resourceSynapseWorkspaceRead(d, meta)
}
