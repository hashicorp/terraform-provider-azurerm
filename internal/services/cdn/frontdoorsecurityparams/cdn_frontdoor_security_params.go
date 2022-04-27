package cdnfrontdoorsecurityparams

import (
	"fmt"

	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontdoorSecurityParameters struct {
	TypeName   track1.Type
	ConfigName string
}

type CdnFrontdoorSecurityMappings struct {
	Firewall CdnFrontdoorSecurityParameters
}

func InitializeCdnFrontdoorSecurityMappings() *CdnFrontdoorSecurityMappings {
	m := new(CdnFrontdoorSecurityMappings)

	m.Firewall = CdnFrontdoorSecurityParameters{
		TypeName:   track1.TypeWebApplicationFirewall,
		ConfigName: "firewall",
	}

	return m
}

func ExpandCdnFrontdoorFirewallPolicyParameters(input []interface{}, isStandardSku bool) (track1.SecurityPolicyWebApplicationFirewallParameters, error) {
	results := track1.SecurityPolicyWebApplicationFirewallParameters{}
	if len(input) == 0 {
		return results, nil
	}

	m := InitializeCdnFrontdoorSecurityMappings()
	associations := make([]track1.SecurityPolicyWebApplicationFirewallAssociation, 0)

	// pull off only the firewall policy from the security_policies list
	policyType := input[0].(map[string]interface{})
	firewallPolicy := policyType[m.Firewall.ConfigName].([]interface{})
	v := firewallPolicy[0].(map[string]interface{})

	if id := v["cdn_frontdoor_firewall_policy_id"].(string); id != "" {
		results.WafPolicy = &track1.ResourceReference{
			ID: utils.String(id),
		}
	}

	configAssociations := v["association"].([]interface{})

	for _, item := range configAssociations {
		v := item.(map[string]interface{})
		domains := expandSecurityPoliciesActivatedResourceReference(v["domain"].([]interface{}))

		if isStandardSku {
			if len(*domains) > 100 {
				return results, fmt.Errorf("the %q sku is only allowed to have 100 or less domains associated with the firewall policy, got %d", track1.SkuNameStandardAzureFrontDoor, len(*domains))
			}
		} else {
			if len(*domains) > 500 {
				return results, fmt.Errorf("the %q sku is only allowed to have 500 or less domains associated with the firewall policy, got %d", track1.SkuNamePremiumAzureFrontDoor, len(*domains))
			}
		}

		association := track1.SecurityPolicyWebApplicationFirewallAssociation{
			Domains:         domains,
			PatternsToMatch: utils.ExpandStringSlice(v["patterns_to_match"].([]interface{})),
		}

		associations = append(associations, association)
	}

	results.Associations = &associations
	results.Type = m.Firewall.TypeName

	return results, nil
}

func expandSecurityPoliciesActivatedResourceReference(input []interface{}) *[]track1.ActivatedResourceReference {
	results := make([]track1.ActivatedResourceReference, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		activatedResourceReference := track1.ActivatedResourceReference{}

		if id := v["cdn_frontdoor_custom_domain_id"].(string); id != "" {
			activatedResourceReference.ID = utils.String(id)

			// This is a read-only field
			// enabled := v["active"].(bool)

			// if !enabled {
			// 	activatedResourceReference.IsActive = utils.Bool(enabled)
			// }

			results = append(results, activatedResourceReference)
		}
	}

	return &results
}

func FlattenCdnFrontdoorFirewallPolicyParameters(input track1.BasicSecurityPolicyPropertiesParameters) ([]interface{}, error) {
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

func flattenSecurityPoliciesActivatedResourceReference(input *[]track1.ActivatedResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		domain := make(map[string]interface{})
		domain["cdn_frontdoor_custom_domain_id"] = *item.ID
		domain["active"] = *item.IsActive

		results = append(results, domain)
	}

	return results
}
