// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PolicySetDefinitionDataSource struct{}

func TestAccDataSourceAzureRMPolicySetDefinition_builtIn(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_set_definition", "test")
	d := PolicySetDefinitionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.builtIn("Audit machines with insecure password security settings"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("095e4ed9-c835-4ab6-9439-b5644362a06c"),
				check.That(data.ResourceName).Key("display_name").HasValue("Audit machines with insecure password security settings"),
				check.That(data.ResourceName).Key("policy_type").HasValue("BuiltIn"),
				check.That(data.ResourceName).Key("parameters").HasValue("{\"EnforcePasswordHistory\":{\"type\":\"String\",\"defaultValue\":\"24\",\"metadata\":{\"description\":\"The Enforce password history setting determines the number of unique new passwords that must be associated with a user account before an old password can be reused.\",\"displayName\":\"Enforce password history\"}},\"IncludeArcMachines\":{\"type\":\"String\",\"allowedValues\":[\"true\",\"false\"],\"defaultValue\":\"false\",\"metadata\":{\"description\":\"By selecting this option, you agree to be charged monthly per Arc connected machine.\",\"displayName\":\"Include Arc connected servers\"}},\"MaximumPasswordAge\":{\"type\":\"String\",\"defaultValue\":\"70\",\"metadata\":{\"description\":\"The Maximum password age setting determines the period of time (in days) that a password can be used before the system requires the user to change it.\",\"displayName\":\"Maximum password age\"}},\"MinimumPasswordAge\":{\"type\":\"String\",\"defaultValue\":\"1\",\"metadata\":{\"description\":\"The Minimum password age setting determines the period of time (in days) that a password must be used before the user can change it.\",\"displayName\":\"Minimum password age\"}},\"MinimumPasswordLength\":{\"type\":\"String\",\"defaultValue\":\"14\",\"metadata\":{\"description\":\"The Minimum password length setting determines the least number of characters that can make up a password for a user account.\",\"displayName\":\"Minimum password length\"}}}"),
				check.That(data.ResourceName).Key("policy_definitions").Exists(),
				check.That(data.ResourceName).Key("policy_definition_reference.#").HasValue("9"),
			),
		},
	})
}

func TestAccDataSourceAzureRMPolicySetDefinition_customByName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_set_definition", "test")
	d := PolicySetDefinitionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.customByName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestPolSet-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestPolSet-display-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("policy_type").HasValue("Custom"),
				check.That(data.ResourceName).Key("parameters").Exists(),
				check.That(data.ResourceName).Key("policy_definitions").Exists(),
				check.That(data.ResourceName).Key("policy_definition_reference.#").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceAzureRMPolicySetDefinition_customByDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_set_definition", "test")
	d := PolicySetDefinitionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.customByDisplayName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestPolSet-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestPolSet-display-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("policy_type").HasValue("Custom"),
				check.That(data.ResourceName).Key("parameters").Exists(),
				check.That(data.ResourceName).Key("policy_definitions").Exists(),
				check.That(data.ResourceName).Key("policy_definition_reference.#").HasValue("1"),
			),
		},
	})
}

func (d PolicySetDefinitionDataSource) builtIn(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_policy_set_definition" "test" {
  display_name = "%s"
}
`, name)
}

func (d PolicySetDefinitionDataSource) customByName(data acceptance.TestData) string {
	template := PolicySetDefinitionResource{}.custom(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_set_definition" "test" {
  name = azurerm_policy_set_definition.test.name
}
`, template)
}

func (d PolicySetDefinitionDataSource) customByDisplayName(data acceptance.TestData) string {
	template := PolicySetDefinitionResource{}.custom(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_set_definition" "test" {
  display_name = azurerm_policy_set_definition.test.display_name
}
`, template)
}
