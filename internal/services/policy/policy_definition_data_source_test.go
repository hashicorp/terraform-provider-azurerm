// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PolicyDefinitionDataSource struct{}

func TestAccDataSourceAzureRMPolicyDefinition_builtIn(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_definition", "test")
	d := PolicyDefinitionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.builtIn("Allowed resource types"),
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

func TestAccDataSourceAzureRMPolicyDefinition_builtInLogAnalytics(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_definition", "test")
	d := PolicyDefinitionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.builtInByName("04d53d87-841c-4f23-8a5b-21564380b55e"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("/providers/Microsoft.Authorization/policyDefinitions/04d53d87-841c-4f23-8a5b-21564380b55e"),
				check.That(data.ResourceName).Key("role_definition_ids.0").Exists(),
			),
		},
	})
}

func TestAccDataSourceAzureRMPolicyDefinition_builtInByName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_definition", "test")
	d := PolicyDefinitionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.builtInByName("a08ec900-254a-4555-9bf5-e42af04b5c5c"),
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

func TestAccDataSourceAzureRMPolicyDefinition_builtIn_AtManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_definition", "test")
	d := PolicyDefinitionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.builtInAtManagementGroup("Allowed resource types"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("/providers/Microsoft.Authorization/policyDefinitions/a08ec900-254a-4555-9bf5-e42af04b5c5c"),
			),
		},
	})
}

func TestAccDataSourceAzureRMPolicyDefinition_customByDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_definition", "test")
	d := PolicyDefinitionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.customByDisplayName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestPol-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestPol-display-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("type").HasValue("Microsoft.Authorization/policyDefinitions"),
				check.That(data.ResourceName).Key("policy_type").HasValue("Custom"),
				check.That(data.ResourceName).Key("policy_rule").HasValue("{\"if\":{\"not\":{\"field\":\"location\",\"in\":\"[parameters('allowedLocations')]\"}},\"then\":{\"effect\":\"audit\"}}"),
				check.That(data.ResourceName).Key("parameters").HasValue("{\"allowedLocations\":{\"type\":\"Array\",\"metadata\":{\"description\":\"The list of allowed locations for resources.\",\"displayName\":\"Allowed locations\",\"strongType\":\"location\"}}}"),
				check.That(data.ResourceName).Key("metadata").Exists(),
				check.That(data.ResourceName).Key("mode").HasValue("All"),
			),
		},
	})
}

func TestAccDataSourceAzureRMPolicyDefinition_customByName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_definition", "test")
	d := PolicyDefinitionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.customByName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestPol-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestPol-display-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("type").HasValue("Microsoft.Authorization/policyDefinitions"),
				check.That(data.ResourceName).Key("policy_type").HasValue("Custom"),
				check.That(data.ResourceName).Key("policy_rule").HasValue("{\"if\":{\"not\":{\"field\":\"location\",\"in\":\"[parameters('allowedLocations')]\"}},\"then\":{\"effect\":\"audit\"}}"),
				check.That(data.ResourceName).Key("parameters").HasValue("{\"allowedLocations\":{\"type\":\"Array\",\"metadata\":{\"description\":\"The list of allowed locations for resources.\",\"displayName\":\"Allowed locations\",\"strongType\":\"location\"}}}"),
				check.That(data.ResourceName).Key("metadata").Exists(),
				check.That(data.ResourceName).Key("mode").HasValue("All"),
			),
		},
	})
}

func (d PolicyDefinitionDataSource) builtIn(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_policy_definition" "test" {
  display_name = "%s"
}
`, name)
}

func (d PolicyDefinitionDataSource) builtInByName(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_policy_definition" "test" {
  name = "%s"
}
`, name)
}

func (d PolicyDefinitionDataSource) builtInAtManagementGroup(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

data "azurerm_policy_definition" "test" {
  display_name          = "%s"
  management_group_name = data.azurerm_client_config.current.tenant_id
}
`, name)
}

func (d PolicyDefinitionDataSource) customByDisplayName(data acceptance.TestData) string {
	template := d.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_definition" "test" {
  display_name = azurerm_policy_definition.test_policy.display_name
}
`, template)
}

func (d PolicyDefinitionDataSource) customByName(data acceptance.TestData) string {
	template := d.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_definition" "test" {
  name = azurerm_policy_definition.test_policy.name
}
`, template)
}

func (d PolicyDefinitionDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_policy_definition" "test_policy" {
  name         = "acctestPol-%d"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestPol-display-%d"

  policy_rule = <<POLICY_RULE
  {
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE

  parameters = <<PARAMETERS
  {
    "allowedLocations": {
      "type": "Array",
      "metadata": {
    	"description": "The list of allowed locations for resources.",
    	"displayName": "Allowed locations",
    	"strongType": "location"
      }
    }
  }
PARAMETERS
}
`, data.RandomInteger, data.RandomInteger)
}
