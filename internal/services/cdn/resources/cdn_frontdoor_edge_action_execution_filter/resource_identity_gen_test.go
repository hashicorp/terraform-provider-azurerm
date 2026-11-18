// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_execution_filter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccCdnFrontdoorEdgeActionExecutionFilter_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_execution_filter", "test")
	r := CdnFrontdoorEdgeActionExecutionFilterResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"edge_action_name":    {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_cdn_frontdoor_edge_action_execution_filter.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_cdn_frontdoor_edge_action_execution_filter.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cdn_frontdoor_edge_action_execution_filter.test", tfjsonpath.New("edge_action_name"), tfjsonpath.New("cdn_frontdoor_edge_action_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cdn_frontdoor_edge_action_execution_filter.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("cdn_frontdoor_edge_action_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cdn_frontdoor_edge_action_execution_filter.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("cdn_frontdoor_edge_action_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
