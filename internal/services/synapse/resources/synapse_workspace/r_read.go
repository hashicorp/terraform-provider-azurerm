// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	identitySQLControlClient := meta.(*clients.Client).Synapse.WorkspaceManagedIdentitySQLControlSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.WorkspaceName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	sqlControlSettings, err := identitySQLControlClient.Get(ctx, id.ResourceGroupName, id.WorkspaceName)
	if err != nil {
		return fmt.Errorf("retrieving Sql Identity Control for %s: %+v", *id, err)
	}

	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.NormalizeNilable(resp.Location))

	flattenIdentities, err := flattenIdentity(resp.Identity)
	if err != nil {
		return err
	}

	if err := d.Set("identity", flattenIdentities); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	if props := resp.WorkspaceProperties; props != nil {
		managedVirtualNetworkEnabled := false
		if props.ManagedVirtualNetwork != nil && strings.EqualFold(*props.ManagedVirtualNetwork, "default") {
			managedVirtualNetworkEnabled = true
			if props.ManagedVirtualNetworkSettings != nil {
				d.Set("data_exfiltration_protection_enabled", props.ManagedVirtualNetworkSettings.PreventDataExfiltration)
				d.Set("linking_allowed_for_aad_tenant_ids", pluginsdk.FlattenSlice(props.ManagedVirtualNetworkSettings.AllowedAadTenantIdsForLinking))
			}
		}
		d.Set("managed_virtual_network_enabled", managedVirtualNetworkEnabled)
		d.Set("storage_data_lake_gen2_filesystem_id", flattenArmWorkspaceDataLakeStorageAccountDetails(props.DefaultDataLakeStorage))
		d.Set("sql_administrator_login", props.SQLAdministratorLogin)
		d.Set("managed_resource_group_name", props.ManagedResourceGroupName)
		d.Set("connectivity_endpoints", pluginsdk.FlattenMapStringPtrString(props.ConnectivityEndpoints))
		d.Set("public_network_access_enabled", resp.PublicNetworkAccess == synapse.WorkspacePublicNetworkAccessEnabled)
		d.Set("azuread_authentication_only", props.AzureADOnlyAuthentication)
		cmk, err := flattenEncryptionDetails(props.Encryption)
		if err != nil {
			return fmt.Errorf("flattening `customer_managed_key`: %+v", err)
		}
		if err := d.Set("customer_managed_key", cmk); err != nil {
			return fmt.Errorf("setting `customer_managed_key`: %+v", err)
		}

		repoType, repo := flattenWorkspaceRepositoryConfiguration(props.WorkspaceRepositoryConfiguration)
		switch repoType {
		case workspaceVSTSConfiguration:
			if err := d.Set("azure_devops_repo", repo); err != nil {
				return fmt.Errorf("setting `azure_devops_repo`: %+v", err)
			}
		case workspaceGitHubConfiguration:
			if err := d.Set("github_repo", repo); err != nil {
				return fmt.Errorf("setting `github_repo`: %+v", err)
			}
		}

		if props.VirtualNetworkProfile != nil {
			d.Set("compute_subnet_id", props.VirtualNetworkProfile.ComputeSubnetID)
		}
		if props.PurviewConfiguration != nil {
			d.Set("purview_id", props.PurviewConfiguration.PurviewResourceID)
		}
	}

	if err := d.Set("sql_identity_control_enabled", flattenIdentityControlSQLSettings(sqlControlSettings)); err != nil {
		return fmt.Errorf("setting `sql_identity_control_enabled`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

const (
	workspaceVSTSConfiguration   = "WorkspaceVSTSConfiguration"
	workspaceGitHubConfiguration = "WorkspaceGitHubConfiguration"
)

func waitSynapseWorkspaceCMKState(ctx context.Context, client *synapse.WorkspacesClient, id *workspaces.WorkspaceId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			"Updating",
			"ActivatingWorkspace",
		},
		Target: []string{
			"Succeeded",
			"Consistent",
			"AwaitingUserAction",
		},
		Refresh:                   synapseWorkspaceCMKUpdateStateRefreshFunc(ctx, client, id),
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("failed waiting for updating %s: %+v", id, err)
	}
	return nil
}

func synapseWorkspaceCMKUpdateStateRefreshFunc(ctx context.Context, client *synapse.WorkspacesClient, id *workspaces.WorkspaceId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id.ResourceGroupName, id.WorkspaceName)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if res.Encryption != nil && res.Encryption.Cmk != nil {
			return res, *res.Encryption.Cmk.Status, nil
		}
		return res, "Succeeded", nil
	}
}

func waitSynapseWorkspaceProvisioningState(ctx context.Context, client *synapse.WorkspacesClient, id *workspaces.WorkspaceId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			"Provisioning",
		},
		Target: []string{
			"Succeeded",
		},
		Refresh:                   synapseWorkspaceProvisioningStateRefreshFunc(ctx, client, id),
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("failed waiting for updating %s: %+v", id, err)
	}
	return nil
}

func synapseWorkspaceProvisioningStateRefreshFunc(ctx context.Context, client *synapse.WorkspacesClient, id *workspaces.WorkspaceId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id.ResourceGroupName, id.WorkspaceName)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}
		return res, *res.ProvisioningState, nil
	}
}
