// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceFirewallPolicySetFlatten(d *pluginsdk.ResourceData, id *firewallpolicies.FirewallPolicyId, model *firewallpolicies.FirewallPolicy) error {
	d.Set("name", id.FirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			basePolicyID := ""
			if props.BasePolicy != nil && props.BasePolicy.Id != nil {
				basePolicyID = *props.BasePolicy.Id
			}
			d.Set("base_policy_id", basePolicyID)

			d.Set("threat_intelligence_mode", pointer.FromEnum(props.ThreatIntelMode))

			if sku := props.Sku; sku != nil {
				d.Set("sku", pointer.FromEnum(sku.Tier))
			}

			if err := d.Set("threat_intelligence_allowlist", flattenFirewallPolicyThreatIntelWhitelist(props.ThreatIntelWhitelist)); err != nil {
				return fmt.Errorf(`setting "threat_intelligence_allowlist": %+v`, err)
			}

			if err := d.Set("dns", flattenFirewallPolicyDNSSetting(props.DnsSettings)); err != nil {
				return fmt.Errorf(`setting "dns": %+v`, err)
			}

			if err := d.Set("intrusion_detection", flattenFirewallPolicyIntrusionDetection(props.IntrusionDetection)); err != nil {
				return fmt.Errorf(`setting "intrusion_detection": %+v`, err)
			}

			if err := d.Set("tls_certificate", flattenFirewallPolicyTransportSecurity(props.TransportSecurity)); err != nil {
				return fmt.Errorf(`setting "tls_certificate": %+v`, err)
			}

			if err := d.Set("child_policies", flattenNetworkSubResourceID(props.ChildPolicies)); err != nil {
				return fmt.Errorf(`setting "child_policies": %+v`, err)
			}

			if err := d.Set("firewalls", flattenNetworkSubResourceID(props.Firewalls)); err != nil {
				return fmt.Errorf(`setting "firewalls": %+v`, err)
			}

			if err := d.Set("rule_collection_groups", flattenNetworkSubResourceID(props.RuleCollectionGroups)); err != nil {
				return fmt.Errorf(`setting "rule_collection_groups": %+v`, err)
			}

			var privateIPRanges []any
			var isAutoLearnPrivateRangeEnabled bool
			if props.Snat != nil {
				privateIPRanges = pluginsdk.FlattenSlice(props.Snat.PrivateRanges)
				isAutoLearnPrivateRangeEnabled = pointer.From(props.Snat.AutoLearnPrivateRanges) == firewallpolicies.AutoLearnPrivateRangesModeEnabled
			}
			if err := d.Set("private_ip_ranges", privateIPRanges); err != nil {
				return fmt.Errorf("setting `private_ip_ranges`: %+v", err)
			}

			if err := d.Set("auto_learn_private_ranges_enabled", isAutoLearnPrivateRangeEnabled); err != nil {
				return fmt.Errorf("setting `auto_learn_private_ranges_enabled`: %+v", err)
			}

			if err := d.Set("insights", flattenFirewallPolicyInsights(props.Insights)); err != nil {
				return fmt.Errorf(`setting "insights": %+v`, err)
			}

			if err := d.Set("explicit_proxy", flattenFirewallPolicyExplicitProxy(props.ExplicitProxy)); err != nil {
				return fmt.Errorf("setting `explicit_proxy`: %+v", err)
			}

			if props.Sql != nil && props.Sql.AllowSqlRedirect != nil {
				if err := d.Set("sql_redirect_allowed", props.Sql.AllowSqlRedirect); err != nil {
					return fmt.Errorf("setting `sql_redirect_allowed`: %+v", err)
				}
			}
		}

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("flattening `tags`: %+v", err)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenFirewallPolicyThreatIntelWhitelist(input *firewallpolicies.FirewallPolicyThreatIntelWhitelist) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"ip_addresses": pluginsdk.FlattenSlice(input.IPAddresses),
			"fqdns":        pluginsdk.FlattenSlice(input.Fqdns),
		},
	}
}

func flattenFirewallPolicyDNSSetting(input *firewallpolicies.DnsSettings) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"servers":       pluginsdk.FlattenSlice(input.Servers),
			"proxy_enabled": pointer.From(input.EnableProxy),
		},
	}
}

func flattenFirewallPolicyIntrusionDetection(input *firewallpolicies.FirewallPolicyIntrusionDetection) []any {
	if input == nil {
		return []any{}
	}

	signatureOverrides := make([]any, 0)
	trafficBypass := make([]any, 0)

	if input.Configuration == nil {
		return []any{
			map[string]any{
				"mode":                pointer.FromEnum(input.Mode),
				"signature_overrides": signatureOverrides,
				"traffic_bypass":      trafficBypass,
			},
		}
	}

	if overrides := input.Configuration.SignatureOverrides; overrides != nil {
		for _, override := range *overrides {
			signatureOverrides = append(signatureOverrides, map[string]any{
				"id":    pointer.From(override.Id),
				"state": pointer.FromEnum(override.Mode),
			})
		}
	}

	if bypasses := input.Configuration.BypassTrafficSettings; bypasses != nil {
		for _, bypass := range *bypasses {
			destinationPorts := make([]string, 0)
			if bypass.DestinationPorts != nil {
				destinationPorts = *bypass.DestinationPorts
			}

			sourceIPGroups := make([]string, 0)
			if bypass.SourceIPGroups != nil {
				sourceIPGroups = *bypass.SourceIPGroups
			}

			destinationIPGroups := make([]string, 0)
			if bypass.DestinationIPGroups != nil {
				destinationIPGroups = *bypass.DestinationIPGroups
			}

			trafficBypass = append(trafficBypass, map[string]any{
				"name":                  pointer.From(bypass.Name),
				"description":           pointer.From(bypass.Description),
				"protocol":              pointer.FromEnum(bypass.Protocol),
				"source_addresses":      pointer.From(bypass.SourceAddresses),
				"destination_addresses": pointer.From(bypass.DestinationAddresses),
				"destination_ports":     destinationPorts,
				"source_ip_groups":      sourceIPGroups,
				"destination_ip_groups": destinationIPGroups,
			})
		}
	}

	return []any{
		map[string]any{
			"mode":                pointer.FromEnum(input.Mode),
			"signature_overrides": signatureOverrides,
			"traffic_bypass":      trafficBypass,
			"private_ranges":      pointer.From(input.Configuration.PrivateRanges),
		},
	}
}

func flattenFirewallPolicyTransportSecurity(input *firewallpolicies.FirewallPolicyTransportSecurity) []any {
	if input == nil || input.CertificateAuthority == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"key_vault_secret_id": input.CertificateAuthority.KeyVaultSecretId,
			"name":                input.CertificateAuthority.Name,
		},
	}
}

func flattenFirewallPolicyInsights(input *firewallpolicies.FirewallPolicyInsights) []any {
	if input == nil {
		return []any{}
	}

	var retentionInDays int
	if input.RetentionDays != nil {
		retentionInDays = int(*input.RetentionDays)
	}

	defaultLogAnalyticsWorkspaceId, logAnalyticsWorkspaces := flattenFirewallPolicyLogAnalyticsResources(input.LogAnalyticsResources)

	return []any{
		map[string]any{
			"enabled":                            pointer.From(input.IsEnabled),
			"retention_in_days":                  retentionInDays,
			"default_log_analytics_workspace_id": defaultLogAnalyticsWorkspaceId,
			"log_analytics_workspace":            logAnalyticsWorkspaces,
		},
	}
}

func flattenFirewallPolicyExplicitProxy(input *firewallpolicies.ExplicitProxy) (result []any) {
	if input == nil {
		return []any{}
	}
	output := map[string]any{
		"enabled":         input.EnableExplicitProxy,
		"http_port":       input.HTTPPort,
		"https_port":      input.HTTPSPort,
		"enable_pac_file": input.EnablePacFile,
		"pac_file_port":   input.PacFilePort,
		"pac_file":        input.PacFile,
	}
	return []any{output}
}

func flattenFirewallPolicyLogAnalyticsResources(input *firewallpolicies.FirewallPolicyLogAnalyticsResources) (string, []any) {
	if input == nil {
		return "", []any{}
	}

	var defaultLogAnalyticsWorkspaceId string
	if input.DefaultWorkspaceId != nil && input.DefaultWorkspaceId.Id != nil {
		defaultLogAnalyticsWorkspaceId = *input.DefaultWorkspaceId.Id
	}

	var workspaceList []any
	if input.Workspaces != nil {
		for _, workspace := range *input.Workspaces {
			loc := location.NormalizeNilable(workspace.Region)

			var id string
			if workspace.WorkspaceId != nil && workspace.WorkspaceId.Id != nil {
				id = *workspace.WorkspaceId.Id
			}

			workspaceList = append(workspaceList, map[string]any{
				"id":                id,
				"firewall_location": loc,
			})
		}
	}

	return defaultLogAnalyticsWorkspaceId, workspaceList
}
