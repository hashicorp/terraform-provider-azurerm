// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccAppServiceCustomHostnameBinding_resourceIdentity(t *testing.T) {
	for _, v := range []string{"ARM_TEST_DNS_ZONE", "ARM_TEST_DATA_RESOURCE_GROUP"} {
		if os.Getenv(v) == "" {
			t.Skipf("Skipping as `%s` was not specified", v)
		}
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_custom_hostname_binding", "test")
	r := AppServiceCustomHostnameBindingResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"name":                {},
		"resource_group_name": {},
		"site_name":           {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_app_service_custom_hostname_binding.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_app_service_custom_hostname_binding.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_app_service_custom_hostname_binding.test", tfjsonpath.New("name"), tfjsonpath.New("hostname")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_app_service_custom_hostname_binding.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_app_service_custom_hostname_binding.test", tfjsonpath.New("site_name"), tfjsonpath.New("app_service_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
