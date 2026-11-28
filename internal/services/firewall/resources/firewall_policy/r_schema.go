// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name which should be used for this Firewall Policy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallPolicyName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Description:  "The SKU Tier of the Firewall Policy.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(firewallpolicies.FirewallPolicySkuTierStandard),
			ValidateFunc: validation.StringInSlice(firewallpolicies.PossibleValuesForFirewallPolicySkuTier(), false),
		},

		"location": commonschema.Location(),

		"base_policy_id": {
			Description:  "The ID of the base Firewall Policy.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: firewallpolicies.ValidateFirewallPolicyID,
		},

		"dns": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"servers": {
						Description: "A list of custom DNS servers' IP addresses.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsIPv4Address,
						},
					},
					"proxy_enabled": {
						Description: "Whether to enable DNS proxy on Firewalls attached to this Firewall Policy.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"threat_intelligence_mode": {
			Description:  "The operation mode for Threat Intelligence.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(firewallpolicies.AzureFirewallThreatIntelModeAlert),
			ValidateFunc: validation.StringInSlice(firewallpolicies.PossibleValuesForAzureFirewallThreatIntelMode(), false),
		},

		"threat_intelligence_allowlist": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_addresses": {
						Description:  "A list of IP addresses or CIDR ranges that will be skipped for threat detection.",
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"threat_intelligence_allowlist.0.ip_addresses", "threat_intelligence_allowlist.0.fqdns"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.Any(validation.IsCIDR, validation.IsIPv4Address),
						},
					},
					"fqdns": {
						Description:  "A list of FQDNs that will be skipped for threat detection.",
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"threat_intelligence_allowlist.0.ip_addresses", "threat_intelligence_allowlist.0.fqdns"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
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
						Description:  "In which mode you want to run intrusion detection: `Off`, `Alert` or `Deny`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(firewallpolicies.PossibleValuesForFirewallPolicyIntrusionDetectionStateType(), false),
					},
					"signature_overrides": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"state": {
									Description:  "State can be any of `Off`, `Alert` or `Deny`.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(firewallpolicies.PossibleValuesForFirewallPolicyIntrusionDetectionStateType(), false),
								},
								"id": {
									Description: "12-digit number (ID) which identifies your signature.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
								},
							},
						},
					},
					"private_ranges": {
						Description: "A list of Private IP address ranges to identify traffic direction. By default, only ranges defined by IANA RFC 1918 are considered private IP addresses.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
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
									Description: "The name which should be used for this bypass traffic setting.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"description": {
									Description: "The description for this bypass traffic setting.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
								},
								"protocol": {
									Description:  "The protocols any of `ANY`, `TCP`, `ICMP`, `UDP` that shall be bypassed by intrusion detection.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(firewallpolicies.PossibleValuesForFirewallPolicyIntrusionDetectionProtocol(), true),
									// protocol to be one of [ICMP ANY TCP UDP] but response may be "Any"
									DiffSuppressFunc: suppress.CaseDifference,
								},
								"source_addresses": {
									Description: "A list of source addresses that shall be bypassed by intrusion detection.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"destination_addresses": {
									Description: "A list of destination IP addresses that shall be bypassed by intrusion detection.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"destination_ports": {
									Description: "A list of destination IP ports that shall be bypassed by intrusion detection.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"source_ip_groups": {
									Description: "A list of source IP groups that shall be bypassed by intrusion detection.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"destination_ip_groups": {
									Description: "A list of destination IP groups that shall be bypassed by intrusion detection.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
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
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_secret_id": {
						Description: "The Secret Identifier (URI) of the certificate stored in Azure Key Vault, either as a secret or certificate.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"name": {
						Description: "The name of the certificate.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
						Description: "Whether the insights functionality is enabled for this Firewall Policy.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
					"default_log_analytics_workspace_id": {
						Description:  "The ID of the default Log Analytics Workspace that the Firewalls associated with this Firewall Policy will send their logs to, when there is no location matches in the `log_analytics_workspace`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: workspaces.ValidateWorkspaceID,
					},
					"retention_in_days": {
						Description:  "The log retention period in days.",
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
									Description:  "The ID of the Log Analytics Workspace that the Firewalls associated with this Firewall Policy will send their logs to when their locations match the `firewall_location`.",
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
						Description: "Whether the explicit proxy is enabled for this Firewall Policy.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"http_port": {
						Description:  "The port number for explicit HTTP protocol.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 35536),
					},
					"https_port": {
						Description:  "The port number for explicit proxy HTTPS protocol.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 35536),
					},
					"enable_pac_file": {
						Description: "Whether the pac file port and URL need to be provided.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"pac_file_port": {
						Description:  "A port number for firewall to serve PAC file.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 35536),
					},
					"pac_file": {
						Description:  "A SAS URL for PAC file.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"sql_redirect_allowed": {
			Description: "Whether SQL Redirect traffic filtering is allowed. Enabling this flag requires no rule using ports between `11000`-`11999`.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"child_policies": {
			Description: "A list of reference to child Firewall Policies of this Firewall Policy.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"firewalls": {
			Description: "A list of references to Azure Firewalls that this Firewall Policy is associated with.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"rule_collection_groups": {
			Description: "A list of references to Firewall Policy Rule Collection Groups that belongs to this Firewall Policy.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"private_ip_ranges": {
			Description: "A list of private IP ranges to which traffic will not be SNAT.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.Any(
					validation.IsCIDR,
					validation.IsIPv4Address,
				),
			},
		},

		"auto_learn_private_ranges_enabled": {
			Description: "Whether enable auto learn private IP range.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"tags": commonschema.Tags(),
	}
}
