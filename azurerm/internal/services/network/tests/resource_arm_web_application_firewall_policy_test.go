package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMWebApplicationFirewallPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWebApplicationFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWebApplicationFirewallPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWebApplicationFirewallPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMWebApplicationFirewallPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWebApplicationFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWebApplicationFirewallPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWebApplicationFirewallPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.name", "Rule1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.priority", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.rule_type", "MatchRule"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.match_variables.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.match_variables.0.variable_name", "RemoteAddr"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.operator", "IPMatch"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.negation_condition", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.match_values.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.match_values.0", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.match_values.1", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.action", "Block"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.name", "Rule2"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.priority", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.rule_type", "MatchRule"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.0.match_variables.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.0.match_variables.0.variable_name", "RemoteAddr"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.0.operator", "IPMatch"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.0.negation_condition", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.0.match_values.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.0.match_values.0", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.match_variables.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.match_variables.0.variable_name", "RequestHeaders"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.match_variables.0.selector", "UserAgent"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.operator", "Contains"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.negation_condition", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.match_values.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.match_values.0", "Windows"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.action", "Block"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMWebApplicationFirewallPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWebApplicationFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWebApplicationFirewallPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWebApplicationFirewallPolicyExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMWebApplicationFirewallPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWebApplicationFirewallPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.name", "Rule1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.priority", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.rule_type", "MatchRule"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.match_variables.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.match_variables.0.variable_name", "RemoteAddr"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.operator", "IPMatch"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.negation_condition", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.match_values.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.match_values.0", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.match_conditions.0.match_values.1", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.0.action", "Block"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.name", "Rule2"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.priority", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.rule_type", "MatchRule"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.0.match_variables.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.0.match_variables.0.variable_name", "RemoteAddr"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.0.operator", "IPMatch"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.0.negation_condition", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.0.match_values.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.0.match_values.0", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.match_variables.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.match_variables.0.variable_name", "RequestHeaders"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.match_variables.0.selector", "UserAgent"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.operator", "Contains"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.negation_condition", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.match_values.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.match_conditions.1.match_values.0", "Windows"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rules.1.action", "Block"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMWebApplicationFirewallPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.WebApplicationFirewallPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Web Application Firewall Policy not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Web Application Firewall Policy %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.WebApplicationFirewallPoliciesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMWebApplicationFirewallPolicyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.WebApplicationFirewallPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_web_application_firewall_policy" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.WebApplicationFirewallPoliciesClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMWebApplicationFirewallPolicy_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_web_application_firewall_policy" "test" {
  name                = "acctestwafpolicy-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMWebApplicationFirewallPolicy_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_web_application_firewall_policy" "test" {
  name                = "acctestwafpolicy-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  custom_rules {
    name      = "Rule1"
    priority  = 1
    rule_type = "MatchRule"

    match_conditions {
      match_variables {
        variable_name = "RemoteAddr"
      }

      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24", "10.0.0.0/24"]
    }

    action = "Block"
  }

  custom_rules {
    name      = "Rule2"
    priority  = 2
    rule_type = "MatchRule"

    match_conditions {
      match_variables {
        variable_name = "RemoteAddr"
      }

      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24"]
    }

    match_conditions {
      match_variables {
        variable_name = "RequestHeaders"
        selector      = "UserAgent"
      }

      operator           = "Contains"
      negation_condition = false
      match_values       = ["Windows"]
    }

    action = "Block"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
