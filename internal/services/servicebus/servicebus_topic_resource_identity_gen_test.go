// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccServicebusTopic_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")
	r := ServicebusTopicResource{}

	checkedFields := map[string]struct{}{
		"namespace_name":      {},
		"resource_group_name": {},
		"subscription_id":     {},
		"topic_name":          {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_servicebus_topic.test", checkedFields),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_servicebus_topic.test", tfjsonpath.New("namespace_name"), tfjsonpath.New("namespace_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_servicebus_topic.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("namespace_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_servicebus_topic.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("namespace_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_servicebus_topic.test", tfjsonpath.New("topic_name"), tfjsonpath.New("name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
