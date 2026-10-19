// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall_policy

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandFirewallPolicyThreatIntelWhitelist(input []any) *firewallpolicies.FirewallPolicyThreatIntelWhitelist {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	return &firewallpolicies.FirewallPolicyThreatIntelWhitelist{
		IPAddresses: pluginsdk.ExpandStringSlice(raw["ip_addresses"].(*pluginsdk.Set).List()),
		Fqdns:       pluginsdk.ExpandStringSlice(raw["fqdns"].(*pluginsdk.Set).List()),
	}
}

func expandFirewallPolicyDNSSetting(input []any) *firewallpolicies.DnsSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	return &firewallpolicies.DnsSettings{
		Servers:     pluginsdk.ExpandStringSlice(raw["servers"].([]any)),
		EnableProxy: pointer.To(raw["proxy_enabled"].(bool)),
	}
}

func expandFirewallPolicyIntrusionDetection(input []any) *firewallpolicies.FirewallPolicyIntrusionDetection {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)

	signatureOverridesRaw := raw["signature_overrides"].([]any)
	signatureOverrides := make([]firewallpolicies.FirewallPolicyIntrusionDetectionSignatureSpecification, 0, len(signatureOverridesRaw))
	for _, v := range signatureOverridesRaw {
		overrides := v.(map[string]any)
		signatureOverrides = append(signatureOverrides, firewallpolicies.FirewallPolicyIntrusionDetectionSignatureSpecification{
			Id:   pointer.To(overrides["id"].(string)),
			Mode: pointer.ToEnum[firewallpolicies.FirewallPolicyIntrusionDetectionStateType](overrides["state"].(string)),
		})
	}

	trafficBypassRaw := raw["traffic_bypass"].([]any)
	trafficBypass := make([]firewallpolicies.FirewallPolicyIntrusionDetectionBypassTrafficSpecifications, 0, len(trafficBypassRaw))
	for _, v := range trafficBypassRaw {
		bypass := v.(map[string]any)
		trafficBypass = append(trafficBypass, firewallpolicies.FirewallPolicyIntrusionDetectionBypassTrafficSpecifications{
			Name:                 pointer.To(bypass["name"].(string)),
			Description:          pointer.To(bypass["description"].(string)),
			Protocol:             pointer.ToEnum[firewallpolicies.FirewallPolicyIntrusionDetectionProtocol](bypass["protocol"].(string)),
			SourceAddresses:      pluginsdk.ExpandStringSlice(bypass["source_addresses"].(*pluginsdk.Set).List()),
			DestinationAddresses: pluginsdk.ExpandStringSlice(bypass["destination_addresses"].(*pluginsdk.Set).List()),
			DestinationPorts:     pluginsdk.ExpandStringSlice(bypass["destination_ports"].(*pluginsdk.Set).List()),
			SourceIPGroups:       pluginsdk.ExpandStringSlice(bypass["source_ip_groups"].(*pluginsdk.Set).List()),
			DestinationIPGroups:  pluginsdk.ExpandStringSlice(bypass["destination_ip_groups"].(*pluginsdk.Set).List()),
		})
	}

	privateRangesRaw := raw["private_ranges"].([]any)
	privateRanges := make([]string, 0, len(privateRangesRaw))
	for _, v := range privateRangesRaw {
		privateRanges = append(privateRanges, v.(string))
	}

	return &firewallpolicies.FirewallPolicyIntrusionDetection{
		Mode: pointer.ToEnum[firewallpolicies.FirewallPolicyIntrusionDetectionStateType](raw["mode"].(string)),
		Configuration: &firewallpolicies.FirewallPolicyIntrusionDetectionConfiguration{
			SignatureOverrides:    &signatureOverrides,
			PrivateRanges:         &privateRanges,
			BypassTrafficSettings: &trafficBypass,
		},
	}
}

func expandFirewallPolicyTransportSecurity(input []any) *firewallpolicies.FirewallPolicyTransportSecurity {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)

	return &firewallpolicies.FirewallPolicyTransportSecurity{
		CertificateAuthority: &firewallpolicies.FirewallPolicyCertificateAuthority{
			KeyVaultSecretId: pointer.To(raw["key_vault_secret_id"].(string)),
			Name:             pointer.To(raw["name"].(string)),
		},
	}
}

func expandFirewallPolicyInsights(input []any) *firewallpolicies.FirewallPolicyInsights {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	return &firewallpolicies.FirewallPolicyInsights{
		IsEnabled:             pointer.To(raw["enabled"].(bool)),
		RetentionDays:         pointer.To(int64(raw["retention_in_days"].(int))),
		LogAnalyticsResources: expandFirewallPolicyLogAnalyticsResources(raw["default_log_analytics_workspace_id"].(string), raw["log_analytics_workspace"].([]any)),
	}
}

func expandFirewallPolicyExplicitProxy(input []any) *firewallpolicies.ExplicitProxy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	if raw == nil {
		return nil
	}

	output := &firewallpolicies.ExplicitProxy{
		EnableExplicitProxy: pointer.To(raw["enabled"].(bool)),
		HTTPPort:            pointer.To(int64(raw["http_port"].(int))),
		HTTPSPort:           pointer.To(int64(raw["https_port"].(int))),
		PacFilePort:         pointer.To(int64(raw["pac_file_port"].(int))),
		PacFile:             pointer.To(raw["pac_file"].(string)),
	}

	if val, ok := raw["enable_pac_file"]; ok {
		output.EnablePacFile = pointer.To(val.(bool))
	}

	return output
}

func expandFirewallPolicyLogAnalyticsResources(defaultWorkspaceId string, workspaces []any) *firewallpolicies.FirewallPolicyLogAnalyticsResources {
	output := &firewallpolicies.FirewallPolicyLogAnalyticsResources{
		DefaultWorkspaceId: &firewallpolicies.SubResource{
			Id: &defaultWorkspaceId,
		},
	}

	workspaceList := make([]firewallpolicies.FirewallPolicyLogAnalyticsWorkspace, 0, len(workspaces))
	for _, workspace := range workspaces {
		workspace := workspace.(map[string]any)
		workspaceList = append(workspaceList, firewallpolicies.FirewallPolicyLogAnalyticsWorkspace{
			Region: pointer.To(location.Normalize(workspace["firewall_location"].(string))),
			WorkspaceId: &firewallpolicies.SubResource{
				Id: pointer.To(workspace["id"].(string)),
			},
		})
	}
	if workspaceList != nil {
		output.Workspaces = &workspaceList
	}

	return output
}
