package cdnfrontdoorsecurityparams

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontdoorSecurityParameters struct {
	TypeName   cdn.Type
	ConfigName string
}

type CdnFrontdoorSecurityMappings struct {
	Firewall CdnFrontdoorSecurityParameters
}

func ExpandCdnFrontdoorFirewallPolicyParameters(input []interface{}, isStandardSku bool) (cdn.SecurityPolicyWebApplicationFirewallParameters, error) {
	results := cdn.SecurityPolicyWebApplicationFirewallParameters{}
	if len(input) == 0 {
		return results, nil
	}

	associations := make([]cdn.SecurityPolicyWebApplicationFirewallAssociation, 0)

	// pull off only the firewall policy from the security_policies list
	policyType := input[0].(map[string]interface{})
	firewallPolicy := policyType["firewall"].([]interface{})
	v := firewallPolicy[0].(map[string]interface{})

	if id := v["cdn_frontdoor_firewall_policy_id"].(string); id != "" {
		results.WafPolicy = &cdn.ResourceReference{
			ID: utils.String(id),
		}
	}

	configAssociations := v["association"].([]interface{})

	for _, item := range configAssociations {
		v := item.(map[string]interface{})
		domains := expandSecurityPoliciesActivatedResourceReference(v["domain"].([]interface{}))

		if isStandardSku {
			if len(*domains) > 100 {
				return results, fmt.Errorf("the %q sku is only allowed to have 100 or less domains associated with the firewall policy, got %d", cdn.SkuNameStandardAzureFrontDoor, len(*domains))
			}
		} else {
			if len(*domains) > 500 {
				return results, fmt.Errorf("the %q sku is only allowed to have 500 or less domains associated with the firewall policy, got %d", cdn.SkuNamePremiumAzureFrontDoor, len(*domains))
			}
		}

		association := cdn.SecurityPolicyWebApplicationFirewallAssociation{
			Domains:         domains,
			PatternsToMatch: utils.ExpandStringSlice(v["patterns_to_match"].([]interface{})),
		}

		associations = append(associations, association)
	}

	results.Associations = &associations

	return results, nil
}

func expandSecurityPoliciesActivatedResourceReference(input []interface{}) *[]cdn.ActivatedResourceReference {
	results := make([]cdn.ActivatedResourceReference, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})

		if id := v["cdn_frontdoor_domain_id"].(string); id != "" {
			results = append(results, cdn.ActivatedResourceReference{
				ID: utils.String(id),
			})
		}
	}

	return &results
}

func FlattenCdnFrontdoorFirewallPolicyParameters(input cdn.BasicSecurityPolicyPropertiesParameters) ([]interface{}, error) {
	waf, ok := input.AsSecurityPolicyWebApplicationFirewallParameters()
	if !ok {
		return nil, fmt.Errorf("expected security policy web application firewall parameters")
	}

	// we know it's a firewall policy at this point,
	// create the objects to hold the policy data
	associations := make([]interface{}, 0)

	wafPolicyId := ""
	if waf.WafPolicy != nil && waf.WafPolicy.ID != nil {
		wafPolicyId = *waf.WafPolicy.ID
	}

	if waf.Associations != nil {
		for _, item := range *waf.Associations {
			associations = append(associations, map[string]interface{}{
				"domain":            flattenSecurityPoliciesActivatedResourceReference(item.Domains),
				"patterns_to_match": utils.FlattenStringSlice(item.PatternsToMatch),
			})
		}
	}

	return []interface{}{
		map[string]interface{}{
			"firewall": []interface{}{
				map[string]interface{}{
					"association":                      associations,
					"cdn_frontdoor_firewall_policy_id": wafPolicyId,
				},
			},
		},
	}, nil
}

func flattenSecurityPoliciesActivatedResourceReference(input *[]cdn.ActivatedResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		frontDoorDomainId := ""
		if item.ID != nil {
			frontDoorDomainId = *item.ID
		}

		active := false
		if item.IsActive != nil {
			active = *item.IsActive
		}

		results = append(results, map[string]interface{}{
			"active":                  active,
			"cdn_frontdoor_domain_id": frontDoorDomainId,
		})
	}

	return results
}
