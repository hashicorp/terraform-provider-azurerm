// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_security_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/securitypolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2024-02-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontdoorSecurityPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitypolicies.ParseSecurityPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, pointer.From(id))
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SecurityPolicyName)
	d.Set("cdn_frontdoor_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			securityPolicies, err := flattenCdnFrontDoorSecurityPolicyResource(props.Parameters)
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

func flattenSecurityPoliciesActivatedResourceReference(input *[]securitypolicies.ActivatedResourceReference) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	for _, item := range *input {
		frontDoorDomainId := ""
		if item.Id != nil {
			if parsedFrontDoorCustomDomainId, frontDoorCustomDomainIdErr := afddomains.ParseCustomDomainIDInsensitively(*item.Id); frontDoorCustomDomainIdErr == nil {
				frontDoorDomainId = parsedFrontDoorCustomDomainId.ID()
			} else if parsedFrontDoorEndpointId, frontDoorEndpointIdErr := afdendpoints.ParseAfdEndpointIDInsensitively(*item.Id); frontDoorEndpointIdErr == nil {
				frontDoorDomainId = parsedFrontDoorEndpointId.ID()
			} else {
				return nil, fmt.Errorf("flattening `cdn_frontdoor_domain_id`: %+v; %+v", frontDoorCustomDomainIdErr, frontDoorEndpointIdErr)
			}
		}

		results = append(results, map[string]any{
			"active":                  pointer.From(item.IsActive),
			"cdn_frontdoor_domain_id": frontDoorDomainId,
		})
	}

	return results, nil
}

func flattenCdnFrontDoorSecurityPolicyResource(input securitypolicies.SecurityPolicyPropertiesParameters) ([]any, error) {
	results := make([]any, 0)
	if input.SecurityPolicyPropertiesParameters().Type != securitypolicies.SecurityPolicyTypeWebApplicationFirewall {
		return results, fmt.Errorf("unexpected security policy `Type` %q, expected `WebApplicationFirewall`", input.SecurityPolicyPropertiesParameters().Type)
	}

	wafParams := input.(securitypolicies.SecurityPolicyWebApplicationFirewallParameters)
	associations := make([]any, 0)
	wafPolicyId := ""

	if wafParams.WafPolicy != nil {
		parsedId, err := webapplicationfirewallpolicies.ParseFrontDoorWebApplicationFirewallPolicyIDInsensitively(pointer.From(wafParams.WafPolicy.Id))
		if err != nil {
			return results, err
		}

		wafPolicyId = parsedId.ID()
	}

	if wafParams.Associations != nil {
		for _, item := range *wafParams.Associations {
			domain, err := flattenSecurityPoliciesActivatedResourceReference(item.Domains)
			if err != nil {
				return results, fmt.Errorf("flattening `domain`: %+v", err)
			}

			associations = append(associations, map[string]any{
				"domain":            domain,
				"patterns_to_match": pluginsdk.FlattenSlice(item.PatternsToMatch),
			})
		}
	}

	results = []any{
		map[string]any{
			"firewall": []any{
				map[string]any{
					"association":                      associations,
					"cdn_frontdoor_firewall_policy_id": wafPolicyId,
				},
			},
		},
	}

	return results, nil
}
