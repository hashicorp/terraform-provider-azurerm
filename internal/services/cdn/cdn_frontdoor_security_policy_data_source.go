// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/securitypolicies"
	waf "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2024-02-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceCdnFrontDoorSecurityPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCdnFrontDoorSecurityPolicyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorSecurityPolicyName,
			},

			"profile_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"cdn_frontdoor_profile_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"security_policies": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"firewall": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"association": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"domain": {
													Type:     pluginsdk.TypeList,
													Computed: true,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"active": {
																Type:     pluginsdk.TypeBool,
																Computed: true,
															},
															"cdn_frontdoor_domain_id": {
																Type:     pluginsdk.TypeString,
																Computed: true,
															},
														},
													},
												},
												"patterns_to_match": {
													Type:     pluginsdk.TypeList,
													Computed: true,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
													},
												},
											},
										},
									},
									"cdn_frontdoor_firewall_policy_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
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

func dataSourceCdnFrontDoorSecurityPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecurityPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := securitypolicies.NewSecurityPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.SecurityPolicyName)
	d.Set("profile_name", id.ProfileName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cdn_frontdoor_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			securityPolicies, err := flattenCdnFrontDoorSecurityPolicyDataSource(props.Parameters)
			if err != nil {
				return fmt.Errorf("flattening `security_policies`: %+v", err)
			}

			if err := d.Set("security_policies", securityPolicies); err != nil {
				return fmt.Errorf("setting `security_policies`: %+v", err)
			}
		}
	}

	return nil
}

func flattenCdnFrontDoorSecurityPolicyDataSource(input securitypolicies.SecurityPolicyPropertiesParameters) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input.SecurityPolicyPropertiesParameters().Type != securitypolicies.SecurityPolicyTypeWebApplicationFirewall {
		return results, fmt.Errorf("unexpected security policy `Type` %q, expected `WebApplicationFirewall`", input.SecurityPolicyPropertiesParameters().Type)
	}

	wafParams := input.(securitypolicies.SecurityPolicyWebApplicationFirewallParameters)
	associations := make([]interface{}, 0)
	wafPolicyId := ""

	if wafParams.WafPolicy != nil {
		parsedId, err := waf.ParseFrontDoorWebApplicationFirewallPolicyIDInsensitively(pointer.From(wafParams.WafPolicy.Id))
		if err != nil {
			return results, fmt.Errorf("flattening `cdn_frontdoor_firewall_policy_id`: %+v", err)
		}

		wafPolicyId = parsedId.ID()
	}

	if wafParams.Associations != nil {
		for _, item := range *wafParams.Associations {
			domain, err := flattenSecurityPoliciesActivatedResourceReference(item.Domains)
			if err != nil {
				return results, fmt.Errorf("flattening `domain`: %+v", err)
			}

			associations = append(associations, map[string]interface{}{
				"domain":            domain,
				"patterns_to_match": utils.FlattenStringSlice(item.PatternsToMatch),
			})
		}
	}

	results = []interface{}{
		map[string]interface{}{
			"firewall": []interface{}{
				map[string]interface{}{
					"association":                      associations,
					"cdn_frontdoor_firewall_policy_id": wafPolicyId,
				},
			},
		},
	}

	return results, nil
}
