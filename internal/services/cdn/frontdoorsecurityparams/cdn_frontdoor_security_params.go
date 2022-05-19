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

func InitializeCdnFrontdoorSecurityMappings() *CdnFrontdoorSecurityMappings {
	m := new(CdnFrontdoorSecurityMappings)

	m.Firewall = CdnFrontdoorSecurityParameters{
		TypeName:   cdn.TypeWebApplicationFirewall,
		ConfigName: "firewall",
	}

	return m
}

func ExpandCdnFrontdoorFirewallPolicyParameters(input []interface{}, isStandardSku bool) (cdn.SecurityPolicyWebApplicationFirewallParameters, error) {
	results := cdn.SecurityPolicyWebApplicationFirewallParameters{}
	if len(input) == 0 {
		return results, nil
	}

	m := InitializeCdnFrontdoorSecurityMappings()
	associations := make([]cdn.SecurityPolicyWebApplicationFirewallAssociation, 0)

	// pull off only the firewall policy from the security_policies list
	policyType := input[0].(map[string]interface{})
	firewallPolicy := policyType[m.Firewall.ConfigName].([]interface{})
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
	results.Type = m.Firewall.TypeName

	return results, nil
}

func expandSecurityPoliciesActivatedResourceReference(input []interface{}) *[]cdn.ActivatedResourceReference {
	results := make([]cdn.ActivatedResourceReference, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		activatedResourceReference := cdn.ActivatedResourceReference{}

		if id := v["cdn_frontdoor_resource_id"].(string); id != "" {
			activatedResourceReference.ID = utils.String(id)

			results = append(results, activatedResourceReference)
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
	securityPolicy := make([]interface{}, 0)
	firewall := make([]interface{}, 0)
	associations := make([]interface{}, 0)
	wafPolicy := make(map[string]interface{})
	firewallPolicy := make(map[string]interface{})

	wafPolicy["cdn_frontdoor_firewall_policy_id"] = *waf.WafPolicy.ID

	for _, item := range *waf.Associations {
		association := make(map[string]interface{})
		association["domain"] = flattenSecurityPoliciesActivatedResourceReference(item.Domains)
		association["patterns_to_match"] = utils.FlattenStringSlice(item.PatternsToMatch)
		associations = append(associations, association)
	}

	wafPolicy["association"] = associations
	firewall = append(firewall, wafPolicy)
	firewallPolicy["firewall"] = firewall
	securityPolicy = append(securityPolicy, firewallPolicy)

	return securityPolicy, nil
}

func flattenSecurityPoliciesActivatedResourceReference(input *[]cdn.ActivatedResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		domain := make(map[string]interface{})
		domain["cdn_frontdoor_resource_id"] = *item.ID
		domain["active"] = *item.IsActive

		results = append(results, domain)
	}

	return results
}
