// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package signalr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccWebPubsubCustomCertificate_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_custom_certificate", "test")
	r := WebPubsubCustomCertificateResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"resource_group_name": {},
		"subscription_id":     {},
		"web_pub_sub_name":    {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_web_pubsub_custom_certificate.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_web_pubsub_custom_certificate.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_web_pubsub_custom_certificate.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("web_pubsub_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_web_pubsub_custom_certificate.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("web_pubsub_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_web_pubsub_custom_certificate.test", tfjsonpath.New("web_pub_sub_name"), tfjsonpath.New("web_pubsub_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
