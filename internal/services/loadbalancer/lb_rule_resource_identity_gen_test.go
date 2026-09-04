// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loadbalancer_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccLbRule_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	checkedFields := map[string]struct{}{
		"load_balancer_name":  {},
		"name":                {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_lb_rule.test", checkedFields),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_lb_rule.test", tfjsonpath.New("load_balancer_name"), tfjsonpath.New("loadbalancer_id")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_lb_rule.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_lb_rule.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("loadbalancer_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_lb_rule.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("loadbalancer_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
