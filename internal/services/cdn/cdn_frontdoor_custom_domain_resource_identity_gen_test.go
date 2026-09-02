// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccCdnFrontdoorCustomDomain_resourceIdentity(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" {
		t.Skipf("`ARM_TEST_DNS_ZONE` must be set for acceptance tests!")
	}

	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontdoorCustomDomainResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"profile_name":        {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_cdn_frontdoor_custom_domain.test", checkedFields),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cdn_frontdoor_custom_domain.test", tfjsonpath.New("name"), tfjsonpath.New("cdn_frontdoor_profile_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cdn_frontdoor_custom_domain.test", tfjsonpath.New("profile_name"), tfjsonpath.New("cdn_frontdoor_profile_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cdn_frontdoor_custom_domain.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("cdn_frontdoor_profile_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cdn_frontdoor_custom_domain.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("cdn_frontdoor_profile_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
