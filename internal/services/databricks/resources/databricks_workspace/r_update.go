// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/accessconnector"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabricksWorkspaceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	acClient := meta.(*clients.Client).DataBricks.AccessConnectorClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model is nil", id)
	}

	model := *existing.Model
	props := model.Properties

	if d.HasChange("sku") {
		if model.Sku == nil {
			model.Sku = &workspaces.Sku{}
		}
		model.Sku.Name = d.Get("sku").(string)
	}

	if d.HasChange("tags") {
		model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("customer_managed_key_enabled") {
		if props.Parameters == nil {
			props.Parameters = &workspaces.WorkspaceCustomParameters{}
		}
		props.Parameters.PrepareEncryption = &workspaces.WorkspaceCustomBooleanParameter{
			Value: d.Get("customer_managed_key_enabled").(bool),
		}
	}

	if d.HasChange("infrastructure_encryption_enabled") {
		if props.Parameters == nil {
			props.Parameters = &workspaces.WorkspaceCustomParameters{}
		}
		props.Parameters.RequireInfrastructureEncryption = &workspaces.WorkspaceCustomBooleanParameter{
			Value: d.Get("infrastructure_encryption_enabled").(bool),
		}
	}

	if d.HasChanges("default_storage_firewall_enabled", "access_connector_id") {
		defaultStorageFirewallEnabled := workspaces.DefaultStorageFirewallDisabled
		defaultStorageFirewallEnabledRaw := d.Get("default_storage_firewall_enabled").(bool)

		if defaultStorageFirewallEnabledRaw {
			defaultStorageFirewallEnabled = workspaces.DefaultStorageFirewallEnabled

			accessConnectorProperties := workspaces.WorkspacePropertiesAccessConnector{}
			accessConnectorIdRaw := d.Get("access_connector_id").(string)
			accessConnectorId, err := accessconnector.ParseAccessConnectorID(accessConnectorIdRaw)
			if err != nil {
				return err
			}

			accessConnector, err := acClient.Get(ctx, *accessConnectorId)
			if err != nil {
				return fmt.Errorf("retrieving Access Connector %s: %+v", accessConnectorId.AccessConnectorName, err)
			}

			if accessConnector.Model.Identity != nil {
				accIdentityId := ""
				for raw := range accessConnector.Model.Identity.IdentityIds {
					identityId, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
					if err != nil {
						return err
					}
					accIdentityId = identityId.ID()
					break
				}

				accessConnectorProperties.Id = *accessConnector.Model.Id
				accessConnectorProperties.IdentityType = workspaces.IdentityType(accessConnector.Model.Identity.Type)
				accessConnectorProperties.UserAssignedIdentityId = &accIdentityId
			}

			props.AccessConnector = &accessConnectorProperties
		}

		props.DefaultStorageFirewall = &defaultStorageFirewallEnabled
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccessRaw := d.Get("public_network_access_enabled").(bool)
		publicNetworkAccess := workspaces.PublicNetworkAccessDisabled
		if publicNetworkAccessRaw {
			publicNetworkAccess = workspaces.PublicNetworkAccessEnabled
		}
		props.PublicNetworkAccess = &publicNetworkAccess
	}

	if d.HasChange("network_security_group_rules_required") {
		props.RequiredNsgRules = pointer.ToEnum[workspaces.RequiredNsgRules](d.Get("network_security_group_rules_required").(string))
	}

	if d.HasChange("custom_parameters") {
		if props.Parameters == nil {
			props.Parameters = &workspaces.WorkspaceCustomParameters{}
		}

		if customParams := d.Get("custom_parameters").([]any); len(customParams) > 0 && customParams[0] != nil {
			config := customParams[0].(map[string]any)
			var pubSubnetAssoc, priSubnetAssoc *string

			pubSub := config["public_subnet_name"].(string)
			priSub := config["private_subnet_name"].(string)

			if v, ok := config["public_subnet_network_security_group_association_id"].(string); ok {
				pubSubnetAssoc = &v
			}

			if v, ok := config["private_subnet_network_security_group_association_id"].(string); ok {
				priSubnetAssoc = &v
			}

			if config["virtual_network_id"].(string) == "" && (pubSub != "" || priSub != "") {
				return fmt.Errorf("`public_subnet_name` and/or `private_subnet_name` cannot be defined if `virtual_network_id` is not set")
			}
			if config["virtual_network_id"].(string) != "" && (pubSub == "" || priSub == "") {
				return fmt.Errorf("`public_subnet_name` and `private_subnet_name` must both have values if `virtual_network_id` is set")
			}
			if pubSub != "" && pubSubnetAssoc == nil {
				return fmt.Errorf("you must define a value for `public_subnet_network_security_group_association_id` if `public_subnet_name` is set")
			}
			if priSub != "" && priSubnetAssoc == nil {
				return fmt.Errorf("you must define a value for `private_subnet_network_security_group_association_id` if `private_subnet_name` is set")
			}

			if v, ok := config["nat_gateway_name"].(string); ok && v != "" {
				props.Parameters.NatGatewayName = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["public_ip_name"].(string); ok && v != "" {
				props.Parameters.PublicIPName = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["storage_account_name"].(string); ok && v != "" {
				props.Parameters.StorageAccountName = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["storage_account_sku_name"].(string); ok && v != "" {
				props.Parameters.StorageAccountSkuName = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["vnet_address_prefix"].(string); ok && v != "" {
				props.Parameters.VnetAddressPrefix = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["machine_learning_workspace_id"].(string); ok && v != "" {
				props.Parameters.AmlWorkspaceId = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["no_public_ip"].(bool); ok {
				props.Parameters.EnableNoPublicIP = &workspaces.WorkspaceNoPublicIPBooleanParameter{
					Value: v,
				}
			}

			if v, ok := config["public_subnet_name"].(string); ok && v != "" {
				props.Parameters.CustomPublicSubnetName = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["private_subnet_name"].(string); ok && v != "" {
				props.Parameters.CustomPrivateSubnetName = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["virtual_network_id"].(string); ok && v != "" {
				props.Parameters.CustomVirtualNetworkId = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}
		}
	}

	if d.HasChanges("managed_services_cmk_key_vault_key_id", "managed_disk_cmk_key_vault_key_id", "managed_disk_cmk_rotation_to_latest_version_enabled") {
		encryption, err := expandDatabricksWorkspaceEncryption(d)
		if err != nil {
			return fmt.Errorf("expanding workspace encryption: %+v", err)
		}
		props.Encryption = encryption
	}

	if d.HasChange("enhanced_security_compliance") {
		props.EnhancedSecurityCompliance = expandWorkspaceEnhancedSecurity(d.Get("enhanced_security_compliance").([]any))
	}

	model.Properties = props

	if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDatabricksWorkspaceRead(d, meta)
}
