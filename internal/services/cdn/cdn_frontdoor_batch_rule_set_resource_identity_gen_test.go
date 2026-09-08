// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccCdnFrontdoorBatchRuleSet_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"profile_name":        {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basicAttachedRoute(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_cdn_frontdoor_batch_rule_set.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_cdn_frontdoor_batch_rule_set.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cdn_frontdoor_batch_rule_set.test", tfjsonpath.New("profile_name"), tfjsonpath.New("cdn_frontdoor_profile_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cdn_frontdoor_batch_rule_set.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("cdn_frontdoor_profile_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cdn_frontdoor_batch_rule_set.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("cdn_frontdoor_profile_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
