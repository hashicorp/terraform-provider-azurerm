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

func TestAccAzureRMFrontDoorFirewallPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_firewall_policy", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoorFirewallPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorFirewallPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("testAccFrontDoorWAF%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "mode", "Prevention"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFrontDoorFirewallPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_firewall_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoorFirewallPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorFirewallPolicyExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMFrontDoorFirewallPolicy_requiresImport),
		},
	})
}

func TestAccAzureRMFrontDoorFirewallPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_firewall_policy", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoorFirewallPolicy_update(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorFirewallPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("testAccFrontDoorWAF%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "mode", "Prevention"),
				),
			},
			{
				Config: testAccAzureRMFrontDoorFirewallPolicy_update(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorFirewallPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("testAccFrontDoorWAF%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "mode", "Prevention"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rule.1.name", "Rule2"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rule.2.name", "Rule3"),
				),
			},
			{
				Config: testAccAzureRMFrontDoorFirewallPolicy_update(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorFirewallPolicyExists(data.ResourceName),
					testCheckAzureRMFrontDoorFirewallPolicyAttrNotExists(data.ResourceName, "custom_rule.1.name"),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("testAccFrontDoorWAF%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "mode", "Prevention"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFrontDoorFirewallPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_firewall_policy", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoorFirewallPolicy_update(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorFirewallPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("testAccFrontDoorWAF%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "mode", "Prevention"),
					resource.TestCheckResourceAttr(data.ResourceName, "redirect_url", "https://www.contoso.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_block_response_status_code", "403"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rule.0.name", "Rule1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_rule.1.name", "Rule2"),
					resource.TestCheckResourceAttr(data.ResourceName, "managed_rule.0.type", "DefaultRuleSet"),
					resource.TestCheckResourceAttr(data.ResourceName, "managed_rule.0.exclusion.0.match_variable", "QueryStringArgNames"),
					resource.TestCheckResourceAttr(data.ResourceName, "managed_rule.0.override.1.exclusion.0.selector", "really_not_suspicious"),
					resource.TestCheckResourceAttr(data.ResourceName, "managed_rule.0.override.1.rule.0.exclusion.0.selector", "innocent"),
					resource.TestCheckResourceAttr(data.ResourceName, "managed_rule.1.type", "Microsoft_BotManagerRuleSet"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMFrontDoorFirewallPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Frontdoor.FrontDoorsPolicyClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Front Door Firewall Policy not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Front Door Firewall Policy %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on FrontDoorsPolicyClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMFrontDoorFirewallPolicyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Frontdoor.FrontDoorsPolicyClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_frontdoor_firewall_policy" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on FrontDoorsPolicyClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testCheckAzureRMFrontDoorFirewallPolicyAttrNotExists(name string, attribute string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if testAttr := rs.Primary.Attributes[attribute]; testAttr != "" {
			return fmt.Errorf("Attribute still exists: %s", attribute)
		}

		return nil
	}
}

func testAccAzureRMFrontDoorFirewallPolicy_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d"
  location = "%s"
}

resource "azurerm_frontdoor_firewall_policy" "test" {
  name                = "testAccFrontDoorWAF%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMFrontDoorFirewallPolicy_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMFrontDoorFirewallPolicy_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_frontdoor_firewall_policy" "import" {
  name                = azurerm_frontdoor_firewall_policy.test.name
  resource_group_name = azurerm_frontdoor_firewall_policy.test.resource_group_name
}
`, template)
}

func testAccAzureRMFrontDoorFirewallPolicy_update(data acceptance.TestData, update bool) string {
	if update {
		return testAccAzureRMFrontDoorFirewallPolicy_updated(data)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d"
  location = "%s"
}

resource "azurerm_frontdoor_firewall_policy" "test" {
  name                              = "testAccFrontDoorWAF%[1]d"
  resource_group_name               = azurerm_resource_group.test.name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  custom_rule {
    name                           = "Rule1"
    enabled                        = true
    priority                       = 1
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "RemoteAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24", "10.0.0.0/24"]
    }
  }

  managed_rule {
    type    = "DefaultRuleSet"
    version = "preview-0.1"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933111"
        enabled = false
        action  = "Block"
      }
    }
  }

  managed_rule {
    type    = "BotProtection"
    version = "preview-0.1"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMFrontDoorFirewallPolicy_updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d"
  location = "%[2]s"
}

resource "azurerm_frontdoor_firewall_policy" "test" {
  name                              = "testAccFrontDoorWAF%[1]d"
  resource_group_name               = azurerm_resource_group.test.name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  custom_rule {
    name                           = "Rule1"
    enabled                        = true
    priority                       = 1
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "RemoteAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24", "10.0.0.0/24"]
    }
  }

  custom_rule {
    name                           = "Rule2"
    enabled                        = true
    priority                       = 2
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "RemoteAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24"]
    }

    match_condition {
      match_variable     = "RequestHeader"
      selector           = "UserAgent"
      operator           = "Contains"
      negation_condition = false
      match_values       = ["windows"]
      transforms         = ["Lowercase", "Trim"]
    }
  }

  custom_rule {
    name                           = "Rule3"
    enabled                        = true
    priority                       = 3
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "SocketAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24"]
    }

    match_condition {
      match_variable     = "RequestHeader"
      selector           = "UserAgent"
      operator           = "Contains"
      negation_condition = false
      match_values       = ["windows"]
      transforms         = ["Lowercase", "Trim"]
    }
  }

  managed_rule {
    type    = "DefaultRuleSet"
    version = "1.0"

    exclusion {
      match_variable = "QueryStringArgNames"
      operator       = "Equals"
      selector       = "not_suspicious"
    }

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "Block"
      }
    }

    override {
      rule_group_name = "SQLI"

      exclusion {
        match_variable = "QueryStringArgNames"
        operator       = "Equals"
        selector       = "really_not_suspicious"
      }

      rule {
        rule_id = "942200"
        action  = "Block"

        exclusion {
          match_variable = "QueryStringArgNames"
          operator       = "Equals"
          selector       = "innocent"
        }
      }
    }
  }

  managed_rule {
    type    = "Microsoft_BotManagerRuleSet"
    version = "1.0"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
