package policy_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
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
				check.That(data.ResourceName).Key("parameters").HasValue("{\"IncludeArcMachines\":{\"type\":\"String\",\"allowedValues\":[\"true\",\"false\"],\"defaultValue\":\"false\",\"metadata\":{\"description\":\"By selecting this option, you agree to be charged monthly per Arc connected machine.\",\"displayName\":\"Include Arc connected servers\"}}}"),
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
