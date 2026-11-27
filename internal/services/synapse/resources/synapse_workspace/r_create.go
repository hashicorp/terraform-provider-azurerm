// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).WorkspaceClient
	subscriptionId := account.From(meta).SubscriptionId
	identitySQLControlClient := synapseClient.From(meta).WorkspaceManagedIdentitySQLControlSettingsClient
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := workspaces.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroupName, id.WorkspaceName)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_workspace", id.ID())
		}
	}

	managedVirtualNetwork := ""
	if d.Get("managed_virtual_network_enabled").(bool) {
		managedVirtualNetwork = "default"
	}

	publicNetworkAccess := synapse.WorkspacePublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = synapse.WorkspacePublicNetworkAccessDisabled
	}

	workspaceInfo := synapse.Workspace{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		WorkspaceProperties: &synapse.WorkspaceProperties{
			DefaultDataLakeStorage:           expandArmWorkspaceDataLakeStorageAccountDetails(d.Get("storage_data_lake_gen2_filesystem_id").(string)),
			ManagedVirtualNetwork:            pointer.To(managedVirtualNetwork),
			PublicNetworkAccess:              publicNetworkAccess,
			SQLAdministratorLogin:            pointer.To(d.Get("sql_administrator_login").(string)),
			SQLAdministratorLoginPassword:    pointer.To(d.Get("sql_administrator_login_password").(string)),
			ManagedResourceGroupName:         pointer.To(d.Get("managed_resource_group_name").(string)),
			WorkspaceRepositoryConfiguration: expandWorkspaceRepositoryConfiguration(d),
			Encryption:                       expandEncryptionDetails(d),
			AzureADOnlyAuthentication:        pointer.To(d.Get("azuread_authentication_only").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	expandedIdentity, err := expandIdentity(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	workspaceInfo.Identity = expandedIdentity

	if purviewId, ok := d.GetOk("purview_id"); ok {
		workspaceInfo.PurviewConfiguration = &synapse.PurviewConfiguration{
			PurviewResourceID: pointer.To(purviewId.(string)),
		}
	}

	if computeSubnetId, ok := d.GetOk("compute_subnet_id"); ok {
		workspaceInfo.VirtualNetworkProfile = &synapse.VirtualNetworkProfile{
			ComputeSubnetID: pointer.To(computeSubnetId.(string)),
		}
	}

	dataExfiltrationProtectionEnabled := d.Get("data_exfiltration_protection_enabled").(bool)

	if dataExfiltrationProtectionEnabled {
		workspaceInfo.ManagedVirtualNetworkSettings = &synapse.ManagedVirtualNetworkSettings{
			PreventDataExfiltration: pointer.To(dataExfiltrationProtectionEnabled),
		}
	}

	if allowedLinkingTenantIds, ok := d.GetOk("linking_allowed_for_aad_tenant_ids"); ok {
		if workspaceInfo.ManagedVirtualNetworkSettings == nil {
			workspaceInfo.ManagedVirtualNetworkSettings = &synapse.ManagedVirtualNetworkSettings{}
		}
		workspaceInfo.ManagedVirtualNetworkSettings.AllowedAadTenantIdsForLinking = pluginsdk.ExpandStringSlice(allowedLinkingTenantIds.([]any))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.WorkspaceName, workspaceInfo)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	if err := waitSynapseWorkspaceCMKState(ctx, client, &id); err != nil {
		return fmt.Errorf("failed waiting for updating %s: %+v", id, err)
	}

	sqlControlSettings := expandIdentityControlSQLSettings(d.Get("sql_identity_control_enabled").(bool))
	future2, err := identitySQLControlClient.CreateOrUpdate(ctx, id.ResourceGroupName, id.WorkspaceName, *sqlControlSettings)
	if err != nil {
		return fmt.Errorf("configuring Sql Identity Control for %s: %+v", id, err)
	}
	if err = future2.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for configuration of Sql Identity Control for %s: %+v", id, err)
	}

	if err := waitSynapseWorkspaceProvisioningState(ctx, client, &id); err != nil {
		return fmt.Errorf("failed waiting for updating %s: %+v", id, err)
	}

	return resourceSynapseWorkspaceRead(d, meta)
}
