// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenWorkspaceManagedIdentity(input *workspaces.ManagedIdentityConfiguration) []any {
	if input == nil {
		return []any{}
	}

	e := make(map[string]any)

	if t := input.PrincipalId; t != nil {
		e["principal_id"] = *t
	}

	if t := input.TenantId; t != nil {
		e["tenant_id"] = *t
	}

	if t := input.Type; t != nil {
		e["type"] = *t
	}

	if len(e) != 0 {
		return []any{e}
	}

	return []any{e}
}

func flattenWorkspaceCustomParameters(input *workspaces.WorkspaceCustomParameters, publicSubnetAssociation, privateSubnetAssociation *string) ([]any, string) {
	if input == nil {
		return []any{}, ""
	}

	var backendAddressPoolId, backendName, loadBalancerId string
	parameters := make(map[string]any)

	if publicSubnetAssociation != nil && *publicSubnetAssociation != "" {
		parameters["public_subnet_network_security_group_association_id"] = *publicSubnetAssociation
	}

	if privateSubnetAssociation != nil && *privateSubnetAssociation != "" {
		parameters["private_subnet_network_security_group_association_id"] = *privateSubnetAssociation
	}

	if v := input.LoadBalancerBackendPoolName; v != nil {
		backendName = v.Value
	}

	if v := input.LoadBalancerId; v != nil {
		loadBalancerId = v.Value
	}

	if v := input.NatGatewayName; v != nil {
		parameters["nat_gateway_name"] = v.Value
	}

	if v := input.PublicIPName; v != nil {
		parameters["public_ip_name"] = v.Value
	}

	if v := input.StorageAccountName; v != nil {
		parameters["storage_account_name"] = v.Value
	}

	if v := input.StorageAccountSkuName; v != nil {
		parameters["storage_account_sku_name"] = v.Value
	}

	if v := input.VnetAddressPrefix; v != nil {
		parameters["vnet_address_prefix"] = v.Value
	}

	if v := input.AmlWorkspaceId; v != nil {
		parameters["machine_learning_workspace_id"] = v.Value
	}

	if v := input.EnableNoPublicIP; v != nil {
		parameters["no_public_ip"] = v.Value
	}

	if v := input.CustomPrivateSubnetName; v != nil {
		parameters["private_subnet_name"] = v.Value
	}

	if v := input.CustomPublicSubnetName; v != nil {
		parameters["public_subnet_name"] = v.Value
	}

	if v := input.CustomVirtualNetworkId; v != nil {
		parameters["virtual_network_id"] = v.Value
	}

	lbId, err := loadbalancers.ParseLoadBalancerIDInsensitively(loadBalancerId)

	if err == nil {
		backendId := loadbalancers.NewLoadBalancerBackendAddressPoolID(lbId.SubscriptionId, lbId.ResourceGroupName, lbId.LoadBalancerName, backendName)
		backendAddressPoolId = backendId.ID()
	}

	return []any{parameters}, backendAddressPoolId
}

func flattenWorkspaceEnhancedSecurity(input *workspaces.EnhancedSecurityComplianceDefinition) []any {
	if input == nil {
		return []any{}
	}

	enhancedSecurityCompliance := make(map[string]any)

	if v := input.AutomaticClusterUpdate; v != nil {
		enhancedSecurityCompliance["automatic_cluster_update_enabled"] = pointer.From(v.Value) != workspaces.AutomaticClusterUpdateValueDisabled
	}

	if v := input.EnhancedSecurityMonitoring; v != nil {
		enhancedSecurityCompliance["enhanced_security_monitoring_enabled"] = pointer.From(v.Value) != workspaces.EnhancedSecurityMonitoringValueDisabled
	}

	if v := input.ComplianceSecurityProfile; v != nil {
		enhancedSecurityCompliance["compliance_security_profile_enabled"] = pointer.From(v.Value) != workspaces.ComplianceSecurityProfileValueDisabled

		standards := pluginsdk.NewSet(pluginsdk.HashString, nil)
		for _, s := range pointer.From(v.ComplianceStandards) {
			if s == string(validate.ComplianceStandardNONE) {
				continue
			}
			standards.Add(s)
		}

		enhancedSecurityCompliance["compliance_security_profile_standards"] = standards
	}

	return []any{enhancedSecurityCompliance}
}
