// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdnfrontdoorsecurityparams

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontdoorSecurityParameters struct {
	TypeName   cdn.Type
	ConfigName string
}

type CdnFrontdoorSecurityMappings struct {
	Firewall CdnFrontdoorSecurityParameters
}

func ExpandCdnFrontdoorFirewallPolicyParameters(input []interface{}, isStandardSku bool) (*cdn.SecurityPolicyWebApplicationFirewallParameters, error) {
	results := cdn.SecurityPolicyWebApplicationFirewallParameters{}
	if len(input) == 0 {
		return &results, nil
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
				return &results, fmt.Errorf("the 'Standard_AzureFrontDoor' sku is only allowed to have 100 or less domains associated with the firewall policy, got %d", len(*domains))
			}
		} else {
			if len(*domains) > 500 {
				return &results, fmt.Errorf("the 'Premium_AzureFrontDoor' sku is only allowed to have 500 or less domains associated with the firewall policy, got %d", len(*domains))
			}
		}

		association := cdn.SecurityPolicyWebApplicationFirewallAssociation{
			Domains:         domains,
			PatternsToMatch: utils.ExpandStringSlice(v["patterns_to_match"].([]interface{})),
		}

		associations = append(associations, association)
	}

	results.Associations = &associations

	return &results, nil
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

func FlattenSecurityPoliciesActivatedResourceReference(input *[]cdn.ActivatedResourceReference) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, item := range *input {
		frontDoorDomainId := ""
		if item.ID != nil {
			if parsedFrontDoorCustomDomainId, frontDoorCustomDomainIdErr := parse.FrontDoorCustomDomainIDInsensitively(*item.ID); frontDoorCustomDomainIdErr == nil {
				frontDoorDomainId = parsedFrontDoorCustomDomainId.ID()
			} else if parsedFrontDoorEndpointId, frontDoorEndpointIdErr := parse.FrontDoorEndpointIDInsensitively(*item.ID); frontDoorEndpointIdErr == nil {
				frontDoorDomainId = parsedFrontDoorEndpointId.ID()
			} else {
				return nil, fmt.Errorf("flattening `cdn_frontdoor_domain_id`: %+v; %+v", frontDoorCustomDomainIdErr, frontDoorEndpointIdErr)
			}
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

	return results, nil
}
