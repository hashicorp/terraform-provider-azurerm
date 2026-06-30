// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package signalr_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccSignalrServiceCustomDomain_resourceIdentity(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_signalr_service_custom_domain", "test")
	r := SignalrServiceCustomDomainResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"resource_group_name": {},
		"signalr_name":        {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_signalr_service_custom_domain.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_signalr_service_custom_domain.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_signalr_service_custom_domain.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("signalr_service_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_signalr_service_custom_domain.test", tfjsonpath.New("signalr_name"), tfjsonpath.New("signalr_service_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_signalr_service_custom_domain.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("signalr_service_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
