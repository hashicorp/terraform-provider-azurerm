package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFrontDoorFirewallPolicy_basic(t *testing.T) {
	resourceName := "azurerm_frontdoor_firewall_policy.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMFrontDoorFirewallPolicy_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorFirewallPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAccFrontDoorWAF%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "mode", "Prevention"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMFrontDoorFirewallPolicy_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_frontdoor_firewall_policy.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoorFirewallPolicy_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorFirewallPolicyExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMFrontDoorFirewallPolicy_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_frontdoor_firewall_policy"),
			},
		},
	})
}

func TestAccAzureRMFrontDoorFirewallPolicy_update(t *testing.T) {
	resourceName := "azurerm_frontdoor_firewall_policy.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMFrontDoorFirewallPolicy_update(ri, "", acceptance.Location())
	configUpdate := testAccAzureRMFrontDoorFirewallPolicy_update(ri, testAccAzureRMFrontDoorFirewallPolicy_updateTemplate(), acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorFirewallPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAccFrontDoorWAF%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "mode", "Prevention"),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorFirewallPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAccFrontDoorWAF%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "mode", "Prevention"),
					resource.TestCheckResourceAttr(resourceName, "custom_rule.1.name", "Rule2"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorFirewallPolicyExists(resourceName),
					testCheckAzureRMFrontDoorFirewallPolicyAttrNotExists(resourceName, "custom_rule.1.name"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAccFrontDoorWAF%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "mode", "Prevention"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMFrontDoorFirewallPolicy_complete(t *testing.T) {
	resourceName := "azurerm_frontdoor_firewall_policy.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMFrontDoorFirewallPolicy_update(ri, testAccAzureRMFrontDoorFirewallPolicy_updateTemplate(), acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorFirewallPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAccFrontDoorWAF%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "mode", "Prevention"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "https://www.contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "custom_block_response_status_code", "403"),
					resource.TestCheckResourceAttr(resourceName, "custom_rule.0.name", "Rule1"),
					resource.TestCheckResourceAttr(resourceName, "custom_rule.1.name", "Rule2"),
					resource.TestCheckResourceAttr(resourceName, "managed_rule.0.type", "DefaultRuleSet"),
					resource.TestCheckResourceAttr(resourceName, "managed_rule.1.type", "BotProtection"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMFrontDoorFirewallPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Front Door Firewall Policy not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).Frontdoor.FrontDoorsPolicyClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMFrontDoorFirewallPolicy_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testAccRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_frontdoor_firewall_policy" "test" {
  name                = "testAccFrontDoorWAF%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}
`, rInt, location)
}

func testAccAzureRMFrontDoorFirewallPolicy_requiresImport(rInt int, location string) string {
	template := testAccAzureRMFrontDoorFirewallPolicy_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_frontdoor_firewall_policy" "import" {
  name                = azurerm_frontdoor_firewall_policy.test.name
  resource_group_name = azurerm_frontdoor_firewall_policy.test.resource_group_name
}
`, template)
}

func testAccAzureRMFrontDoorFirewallPolicy_updateTemplate() string {
	return fmt.Sprintf(`
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
`)
}

func testAccAzureRMFrontDoorFirewallPolicy_update(rInt int, sTemplate string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testAccRG-%[1]d"
  location = "%[3]s"
}

resource "azurerm_frontdoor_firewall_policy" "test" {
  name                              = "testAccFrontDoorWAF%[1]d"
  resource_group_name               = "${azurerm_resource_group.test.name}"
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

  %[2]s

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
`, rInt, sTemplate, location)
}
