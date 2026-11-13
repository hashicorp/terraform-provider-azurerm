// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandDatabricksWorkspaceEncryption(d *pluginsdk.ResourceData) (*workspaces.WorkspacePropertiesEncryption, error) {
	diskCMK := d.Get("managed_disk_cmk_key_vault_key_id").(string)
	servicesCMK := d.Get("managed_services_cmk_key_vault_key_id").(string)
	if diskCMK == "" && servicesCMK == "" {
		return nil, nil
	}

	result := &workspaces.WorkspacePropertiesEncryption{
		Entities: workspaces.EncryptionEntitiesDefinition{},
	}

	if diskCMK != "" {
		key, err := keyvault.ParseNestedItemID(diskCMK, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
		if err != nil {
			return nil, err
		}

		rotateToLatest := (*bool)(nil)
		if d.Get("managed_disk_cmk_rotation_to_latest_version_enabled").(bool) {
			rotateToLatest = pointer.To(true)
		}

		result.Entities.ManagedDisk = &workspaces.ManagedDiskEncryption{
			KeySource: workspaces.EncryptionKeySourceMicrosoftPointKeyvault,
			KeyVaultProperties: workspaces.ManagedDiskEncryptionKeyVaultProperties{
				KeyName:     key.Name,
				KeyVaultUri: key.KeyVaultBaseURL,
				KeyVersion:  key.Version,
			},
			RotationToLatestKeyVersionEnabled: rotateToLatest,
		}
	}

	if servicesCMK != "" {
		key, err := keyvault.ParseNestedItemID(servicesCMK, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
		if err != nil {
			return nil, err
		}

		result.Entities.ManagedServices = &workspaces.EncryptionV2{
			KeySource: workspaces.EncryptionKeySourceMicrosoftPointKeyvault,
			KeyVaultProperties: &workspaces.EncryptionV2KeyVaultProperties{
				KeyName:     key.Name,
				KeyVaultUri: key.KeyVaultBaseURL,
				KeyVersion:  key.Version,
			},
		}
	}

	return result, nil
}

func expandWorkspaceCustomParameters(input []any, customerManagedKeyEnabled, infrastructureEncryptionEnabled bool, backendAddressPoolName, loadBalancerId string) (workspaceCustomParameters *workspaces.WorkspaceCustomParameters, publicSubnetAssociation, privateSubnetAssociation *string) {
	if len(input) == 0 || input[0] == nil {
		// This will be hit when there are no custom params set but we still
		// need to pass the customerManagedKeyEnabled and infrastructureEncryptionEnabled
		// flags anyway...
		parameters := workspaces.WorkspaceCustomParameters{}

		parameters.PrepareEncryption = &workspaces.WorkspaceCustomBooleanParameter{
			Value: customerManagedKeyEnabled,
		}

		parameters.RequireInfrastructureEncryption = &workspaces.WorkspaceCustomBooleanParameter{
			Value: infrastructureEncryptionEnabled,
		}

		return &parameters, nil, nil
	}

	config := input[0].(map[string]any)
	var pubSubnetAssoc, priSubnetAssoc *string
	parameters := workspaces.WorkspaceCustomParameters{}

	if v, ok := config["public_subnet_network_security_group_association_id"].(string); ok && v != "" {
		pubSubnetAssoc = &v
	}

	if v, ok := config["private_subnet_network_security_group_association_id"].(string); ok && v != "" {
		priSubnetAssoc = &v
	}

	if backendAddressPoolName != "" {
		parameters.LoadBalancerBackendPoolName = &workspaces.WorkspaceCustomStringParameter{
			Value: backendAddressPoolName,
		}
	}

	if loadBalancerId != "" {
		parameters.LoadBalancerId = &workspaces.WorkspaceCustomStringParameter{
			Value: loadBalancerId,
		}
	}

	if v, ok := config["nat_gateway_name"].(string); ok && v != "" {
		parameters.NatGatewayName = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["public_ip_name"].(string); ok && v != "" {
		parameters.PublicIPName = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["storage_account_name"].(string); ok && v != "" {
		parameters.StorageAccountName = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["storage_account_sku_name"].(string); ok && v != "" {
		parameters.StorageAccountSkuName = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["vnet_address_prefix"].(string); ok && v != "" {
		parameters.VnetAddressPrefix = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["machine_learning_workspace_id"].(string); ok && v != "" {
		parameters.AmlWorkspaceId = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["no_public_ip"].(bool); ok {
		parameters.EnableNoPublicIP = &workspaces.WorkspaceNoPublicIPBooleanParameter{
			Value: v,
		}
	}

	if v, ok := config["public_subnet_name"].(string); ok && v != "" {
		parameters.CustomPublicSubnetName = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	parameters.PrepareEncryption = &workspaces.WorkspaceCustomBooleanParameter{
		Value: customerManagedKeyEnabled,
	}

	parameters.RequireInfrastructureEncryption = &workspaces.WorkspaceCustomBooleanParameter{
		Value: infrastructureEncryptionEnabled,
	}

	if v, ok := config["private_subnet_name"].(string); ok && v != "" {
		parameters.CustomPrivateSubnetName = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["virtual_network_id"].(string); ok && v != "" {
		parameters.CustomVirtualNetworkId = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	return &parameters, pubSubnetAssoc, priSubnetAssoc
}

func expandWorkspaceEnhancedSecurity(input []any) *workspaces.EnhancedSecurityComplianceDefinition {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]any)

	automaticClusterUpdateEnabled := workspaces.AutomaticClusterUpdateValueDisabled
	if enabled, ok := config["automatic_cluster_update_enabled"].(bool); ok && enabled {
		automaticClusterUpdateEnabled = workspaces.AutomaticClusterUpdateValueEnabled
	}

	enhancedSecurityMonitoringEnabled := workspaces.EnhancedSecurityMonitoringValueDisabled
	if enabled, ok := config["enhanced_security_monitoring_enabled"].(bool); ok && enabled {
		enhancedSecurityMonitoringEnabled = workspaces.EnhancedSecurityMonitoringValueEnabled
	}

	complianceSecurityProfileEnabled := workspaces.ComplianceSecurityProfileValueDisabled
	if enabled, ok := config["compliance_security_profile_enabled"].(bool); ok && enabled {
		complianceSecurityProfileEnabled = workspaces.ComplianceSecurityProfileValueEnabled
	}

	complianceStandards := make([]string, 0)
	if standardSet, ok := config["compliance_security_profile_standards"].(*pluginsdk.Set); ok {
		for _, s := range standardSet.List() {
			complianceStandards = append(complianceStandards, s.(string))
		}
	}

	if complianceSecurityProfileEnabled == workspaces.ComplianceSecurityProfileValueEnabled && len(complianceStandards) == 0 {
		complianceStandards = append(complianceStandards, string(validate.ComplianceStandardNONE))
	}

	return &workspaces.EnhancedSecurityComplianceDefinition{
		AutomaticClusterUpdate: &workspaces.AutomaticClusterUpdateDefinition{
			Value: &automaticClusterUpdateEnabled,
		},
		EnhancedSecurityMonitoring: &workspaces.EnhancedSecurityMonitoringDefinition{
			Value: &enhancedSecurityMonitoringEnabled,
		},
		ComplianceSecurityProfile: &workspaces.ComplianceSecurityProfileDefinition{
			Value:               &complianceSecurityProfileEnabled,
			ComplianceStandards: &complianceStandards,
		},
	}
}
