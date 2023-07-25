// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PolicyDefinitionBuiltInDataSource struct{}

func TestAccDataSourceAzureRMPolicyDefinitionBuiltIn_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_definition_built_in", "test")
	d := PolicyDefinitionBuiltInDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic("Allowed resource types"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("/providers/Microsoft.Authorization/policyDefinitions/a08ec900-254a-4555-9bf5-e42af04b5c5c"),
				check.That(data.ResourceName).Key("name").HasValue("a08ec900-254a-4555-9bf5-e42af04b5c5c"),
				check.That(data.ResourceName).Key("display_name").HasValue("Allowed resource types"),
				check.That(data.ResourceName).Key("type").HasValue("Microsoft.Authorization/policyDefinitions"),
				check.That(data.ResourceName).Key("description").HasValue("This policy enables you to specify the resource types that your organization can deploy. Only resource types that support 'tags' and 'location' will be affected by this policy. To restrict all resources please duplicate this policy and change the 'mode' to 'All'."),
				check.That(data.ResourceName).Key("mode").HasValue("Indexed"),
			),
		},
	})
}

func (d PolicyDefinitionBuiltInDataSource) basic(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_policy_definition_built_in" "test" {
  display_name = "%s"
}
`, name)
}
