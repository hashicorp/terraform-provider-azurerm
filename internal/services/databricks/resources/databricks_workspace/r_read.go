// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabricksWorkspaceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if sku := model.Sku; sku != nil {
			d.Set("sku", sku.Name)
		}

		managedResourceGroupID, err := commonids.ParseResourceGroupIDInsensitively(pointer.From(model.Properties.ManagedResourceGroupId))
		if err != nil {
			return err
		}

		if !features.SixPointOh() {
			d.Set("managed_resource_group_id", model.Properties.ManagedResourceGroupId)
		} else {
			d.Set("managed_resource_group_id", managedResourceGroupID.ID())
		}

		d.Set("managed_resource_group_name", managedResourceGroupID.ResourceGroupName)

		if defaultStorageFirewall := model.Properties.DefaultStorageFirewall; defaultStorageFirewall != nil {
			d.Set("default_storage_firewall_enabled", *defaultStorageFirewall != workspaces.DefaultStorageFirewallDisabled)
			if model.Properties.AccessConnector != nil {
				d.Set("access_connector_id", model.Properties.AccessConnector.Id)
			}
		}

		if publicNetworkAccess := model.Properties.PublicNetworkAccess; publicNetworkAccess != nil {
			d.Set("public_network_access_enabled", *publicNetworkAccess != workspaces.PublicNetworkAccessDisabled)
			if *publicNetworkAccess == workspaces.PublicNetworkAccessDisabled {
				if model.Properties.RequiredNsgRules != nil {
					d.Set("network_security_group_rules_required", string(*model.Properties.RequiredNsgRules))
				}
			}
		}

		var cmkEnabled, infraEnabled bool
		if parameters := model.Properties.Parameters; parameters != nil {
			if parameters.PrepareEncryption != nil {
				cmkEnabled = parameters.PrepareEncryption.Value
			}
			d.Set("customer_managed_key_enabled", cmkEnabled)

			if parameters.RequireInfrastructureEncryption != nil {
				infraEnabled = parameters.RequireInfrastructureEncryption.Value
			}
			d.Set("infrastructure_encryption_enabled", infraEnabled)

			// The subnet associations only exist in the statefile, so we need to do a Get before we Set
			// with what has come back from the Azure response...
			_, pubSubAssoc, priSubAssoc := expandWorkspaceCustomParameters(d.Get("custom_parameters").([]any), cmkEnabled, infraEnabled, "", "")

			custom, backendPoolReadId := flattenWorkspaceCustomParameters(parameters, pubSubAssoc, priSubAssoc)
			if err := d.Set("custom_parameters", custom); err != nil {
				return fmt.Errorf("setting `custom_parameters`: %+v", err)
			}

			d.Set("load_balancer_backend_address_pool_id", backendPoolReadId)
		}

		if err := d.Set("storage_account_identity", flattenWorkspaceManagedIdentity(model.Properties.StorageAccountIdentity)); err != nil {
			return fmt.Errorf("setting `storage_account_identity`: %+v", err)
		}

		if err := d.Set("managed_disk_identity", flattenWorkspaceManagedIdentity(model.Properties.ManagedDiskIdentity)); err != nil {
			return fmt.Errorf("setting `managed_disk_identity`: %+v", err)
		}

		d.Set("workspace_url", pointer.From(model.Properties.WorkspaceURL))
		d.Set("workspace_id", pointer.From(model.Properties.WorkspaceId))

		// customer managed key for managed services
		var servicesKeyId string
		if encryption := model.Properties.Encryption; encryption != nil {
			if encryptionProps := encryption.Entities.ManagedServices; encryptionProps != nil {
				if encryptionProps.KeyVaultProperties.KeyVaultUri != "" {
					key, err := keyvault.NewNestedItemID(encryptionProps.KeyVaultProperties.KeyVaultUri, keyvault.NestedItemTypeKey, encryptionProps.KeyVaultProperties.KeyName, encryptionProps.KeyVaultProperties.KeyVersion)
					if err == nil {
						servicesKeyId = key.ID()
					}
				}
			}
		}

		// customer managed key for managed disk
		var diskKeyId string
		var encryptDiskRotationEnabled bool
		if encryption := model.Properties.Encryption; encryption != nil {
			if encryptionProps := encryption.Entities.ManagedDisk; encryptionProps != nil {
				if encryptionProps.KeyVaultProperties.KeyVaultUri != "" {
					key, err := keyvault.NewNestedItemID(encryptionProps.KeyVaultProperties.KeyVaultUri, keyvault.NestedItemTypeKey, encryptionProps.KeyVaultProperties.KeyName, encryptionProps.KeyVaultProperties.KeyVersion)
					if err == nil {
						diskKeyId = key.ID()
					}
				}

				encryptDiskRotationEnabled = *encryptionProps.RotationToLatestKeyVersionEnabled
			}
		}

		d.Set("enhanced_security_compliance", flattenWorkspaceEnhancedSecurity(model.Properties.EnhancedSecurityCompliance))
		d.Set("disk_encryption_set_id", pointer.From(model.Properties.DiskEncryptionSetId))

		// Always set these even if they are empty to keep the state file
		// consistent with the configuration file...
		d.Set("managed_services_cmk_key_vault_key_id", servicesKeyId)
		d.Set("managed_disk_cmk_key_vault_key_id", diskKeyId)
		d.Set("managed_disk_cmk_rotation_to_latest_version_enabled", encryptDiskRotationEnabled)

		if !features.SixPointOh() {
			d.Set("managed_services_cmk_key_vault_id", d.Get("managed_services_cmk_key_vault_id").(string))
			d.Set("managed_disk_cmk_key_vault_id", d.Get("managed_disk_cmk_key_vault_id").(string))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func workspaceCustomParametersString() []string {
	return []string{
		"custom_parameters.0.machine_learning_workspace_id", "custom_parameters.0.no_public_ip",
		"custom_parameters.0.public_subnet_name", "custom_parameters.0.private_subnet_name", "custom_parameters.0.virtual_network_id",
		"custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id",
		"custom_parameters.0.nat_gateway_name", "custom_parameters.0.public_ip_name", "custom_parameters.0.storage_account_name", "custom_parameters.0.storage_account_sku_name",
		"custom_parameters.0.vnet_address_prefix",
	}
}

func checkSubnetDelegations(ctx context.Context, client *subnets.SubnetsClient, vnetID, publicSubnetName, privateSubnetName string) error {
	requiredDelegationService := "Microsoft.Databricks/workspaces"

	if vnetID == "" || (publicSubnetName == "" && privateSubnetName == "") {
		return nil
	}

	id, err := commonids.ParseVirtualNetworkID(vnetID)
	if err != nil {
		return err
	}

	if publicSubnetName != "" {
		subnetID := commonids.NewSubnetID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName, publicSubnetName)
		resp, err := client.Get(ctx, subnetID, subnets.DefaultGetOperationOptions())
		if err != nil || resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("failed to check public subnet delegation for %s: %s", publicSubnetName, err)
		}
		if resp.Model.Properties.Delegations == nil {
			return fmt.Errorf("required public subnet delegation to %s on %s not found", requiredDelegationService, publicSubnetName)
		}

		if delegations := resp.Model.Properties.Delegations; delegations == nil {
			return fmt.Errorf("required public subnet delegation to %s on %s not found", requiredDelegationService, publicSubnetName)
		} else {
			found := false
			for _, v := range *delegations {
				if v.Properties == nil {
					continue
				}
				if pointer.From(v.Properties.ServiceName) == requiredDelegationService {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("required public subnet delegation to %s on %s not found", requiredDelegationService, publicSubnetName)
			}
		}
	}

	if privateSubnetName != "" {
		subnetID := commonids.NewSubnetID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName, privateSubnetName)
		resp, err := client.Get(ctx, subnetID, subnets.DefaultGetOperationOptions())
		if err != nil || resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("failed to check private subnet delegation for %s: %s", privateSubnetName, err)
		}

		if delegations := resp.Model.Properties.Delegations; delegations == nil {
			return fmt.Errorf("required private subnet delegation to %s on %s not found", requiredDelegationService, privateSubnetName)
		} else {
			found := false
			for _, v := range *delegations {
				if v.Properties == nil {
					continue
				}
				if pointer.From(v.Properties.ServiceName) == requiredDelegationService {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("required private subnet delegation to %s on %s not found", requiredDelegationService, privateSubnetName)
			}
		}
	}

	return nil
}
