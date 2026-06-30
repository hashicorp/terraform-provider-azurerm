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

func TestAccWebPubsubCustomDomain_resourceIdentity(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_custom_domain", "test")
	r := WebPubsubCustomDomainResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"resource_group_name": {},
		"subscription_id":     {},
		"web_pubsub_name":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_web_pubsub_custom_domain.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_web_pubsub_custom_domain.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_web_pubsub_custom_domain.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("web_pubsub_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_web_pubsub_custom_domain.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("web_pubsub_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_web_pubsub_custom_domain.test", tfjsonpath.New("web_pubsub_name"), tfjsonpath.New("web_pubsub_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
