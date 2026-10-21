// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceFirewallPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallPolicyName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(firewallpolicies.FirewallPolicySkuTierStandard),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(firewallpolicies.PossibleValuesForFirewallPolicySkuTier(), false),
		},

		"location": commonschema.Location(),

		"base_policy_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: firewallpolicies.ValidateFirewallPolicyID,
		},

		"dns": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"servers": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsIPv4Address,
						},
					},
					"proxy_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"threat_intelligence_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(firewallpolicies.AzureFirewallThreatIntelModeAlert),
			ValidateFunc: validation.StringInSlice(firewallpolicies.PossibleValuesForAzureFirewallThreatIntelMode(), false),
		},

		"threat_intelligence_allowlist": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_addresses": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.Any(validation.IsCIDR, validation.IsIPv4Address),
						},
						AtLeastOneOf: []string{"threat_intelligence_allowlist.0.ip_addresses", "threat_intelligence_allowlist.0.fqdns"},
					},
					"fqdns": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						AtLeastOneOf: []string{"threat_intelligence_allowlist.0.ip_addresses", "threat_intelligence_allowlist.0.fqdns"},
					},
				},
			},
		},

		"intrusion_detection": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(firewallpolicies.PossibleValuesForFirewallPolicyIntrusionDetectionStateType(), false),
						Optional:     true,
					},
					"signature_overrides": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"state": {
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringInSlice(firewallpolicies.PossibleValuesForFirewallPolicyIntrusionDetectionStateType(), false),
									Optional:     true,
								},
								"id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
					"private_ranges": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"traffic_bypass": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"description": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"protocol": {
									Type:     pluginsdk.TypeString,
									Required: true,
									// protocol to be one of [ICMP ANY TCP UDP] but response may be "Any"
									DiffSuppressFunc: suppress.CaseDifference,
									ValidateFunc:     validation.StringInSlice(firewallpolicies.PossibleValuesForFirewallPolicyIntrusionDetectionProtocol(), true),
								},
								"source_addresses": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"destination_addresses": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"destination_ports": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"source_ip_groups": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"destination_ip_groups": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tls_certificate": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_secret_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"insights": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"default_log_analytics_workspace_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: workspaces.ValidateWorkspaceID,
					},
					"retention_in_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"log_analytics_workspace": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: workspaces.ValidateWorkspaceID,
								},
								"firewall_location": commonschema.LocationWithoutForceNew(),
							},
						},
					},
				},
			},
		},

		"explicit_proxy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"http_port": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 35536),
					},
					"https_port": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 35536),
					},
					"enable_pac_file": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"pac_file_port": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 35536),
					},
					"pac_file": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"sql_redirect_allowed": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"child_policies": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"firewalls": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"rule_collection_groups": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"private_ip_ranges": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.Any(
					validation.IsCIDR,
					validation.IsIPv4Address,
				),
			},
		},

		"auto_learn_private_ranges_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tags": commonschema.Tags(),
	}
}
