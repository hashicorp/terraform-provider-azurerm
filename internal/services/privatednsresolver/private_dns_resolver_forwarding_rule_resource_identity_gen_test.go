// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccPrivateDnsResolverForwardingRule_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_forwarding_rule", "test")
	r := PrivateDnsResolverForwardingRuleResource{}

	checkedFields := map[string]struct{}{
		"dns_forwarding_ruleset_name": {},
		"name":                        {},
		"resource_group_name":         {},
		"subscription_id":             {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_private_dns_resolver_forwarding_rule.test", checkedFields),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_private_dns_resolver_forwarding_rule.test", tfjsonpath.New("dns_forwarding_ruleset_name"), tfjsonpath.New("dns_forwarding_ruleset_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_private_dns_resolver_forwarding_rule.test", tfjsonpath.New("name"), tfjsonpath.New("dns_forwarding_ruleset_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_private_dns_resolver_forwarding_rule.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("dns_forwarding_ruleset_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_private_dns_resolver_forwarding_rule.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("dns_forwarding_ruleset_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
