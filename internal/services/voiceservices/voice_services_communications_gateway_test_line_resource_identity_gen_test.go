// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package voiceservices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccVoiceServicesCommunicationsGatewayTestLine_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_voice_services_communications_gateway_test_line", "test")
	r := VoiceServicesCommunicationsGatewayTestLineResource{}

	checkedFields := map[string]struct{}{
		"name":                        {},
		"communications_gateway_name": {},
		"resource_group_name":         {},
		"subscription_id":             {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_voice_services_communications_gateway_test_line.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_voice_services_communications_gateway_test_line.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_voice_services_communications_gateway_test_line.test", tfjsonpath.New("communications_gateway_name"), tfjsonpath.New("voice_services_communications_gateway_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_voice_services_communications_gateway_test_line.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("voice_services_communications_gateway_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_voice_services_communications_gateway_test_line.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("voice_services_communications_gateway_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
