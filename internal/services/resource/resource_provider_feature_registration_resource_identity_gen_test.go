// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccResourceProviderFeatureRegistration_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_feature_registration", "test")
	r := ResourceProviderFeatureRegistrationResource{}

	checkedFields := map[string]struct{}{
		"subscription_id": {},
		"name":            {},
		"provider_name":   {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basicForResourceIdentity(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_resource_provider_feature_registration.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_resource_provider_feature_registration.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Secondary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_resource_provider_feature_registration.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_resource_provider_feature_registration.test", tfjsonpath.New("provider_name"), tfjsonpath.New("provider_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
