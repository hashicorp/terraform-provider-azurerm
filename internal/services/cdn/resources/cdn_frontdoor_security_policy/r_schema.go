// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_security_policy

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2024-02-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cdnFrontdoorSecurityPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorSecurityPolicyName,
		},

		"cdn_frontdoor_profile_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: profiles.ValidateProfileID,
		},

		"security_policies": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,

			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"firewall": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,

						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"cdn_frontdoor_firewall_policy_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: webapplicationfirewallpolicies.ValidateFrontDoorWebApplicationFirewallPolicyID,
								},

								"association": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,

									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											// NOTE: The max number of domains vary depending on sku: 100 Standard, 500 Premium
											"domain": {
												Type:     pluginsdk.TypeList,
												Required: true,
												MaxItems: 500,

												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"cdn_frontdoor_domain_id": {
															Type:     pluginsdk.TypeString,
															Required: true,
															ValidateFunc: validation.Any(
																afddomains.ValidateCustomDomainID,
																afdendpoints.ValidateAfdEndpointID,
															),
														},

														"active": {
															Type:     pluginsdk.TypeBool,
															Computed: true,
														},
													},
												},
											},

											// NOTE: Per the service team the only acceptable value as of GA is "/*"
											"patterns_to_match": {
												Type:     pluginsdk.TypeList,
												Required: true,
												ForceNew: true,
												MaxItems: 1,

												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
													ValidateFunc: validation.StringInSlice([]string{
														"/*",
													}, false),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
