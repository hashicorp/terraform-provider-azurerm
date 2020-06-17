package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFirewallPolicyPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicyPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFirewallPolicyPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicyPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFirewallPolicyPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicyPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMFirewallPolicyPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMFirewallPolicyPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFirewallPolicyPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicyPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyPolicyExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMFirewallPolicyPolicy_requiresImport),
		},
	})
}

func TestAccAzureRMFirewallPolicyPolicy_inherit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicyPolicy_inherit(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMFirewallPolicyPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.FirewallPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Firewall Policy Policy not found: %s", resourceName)
		}

		id, err := parse.FirewallPolicyPolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Firewall Policy Policy %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Getting on Network.FirewallPolicies: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMFirewallPolicyPolicyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.FirewallPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_firewall_policy_policy" {
			continue
		}

		id, err := parse.FirewallPolicyPolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Getting on Network.FirewallPolicies: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMFirewallPolicyPolicy_basic(data acceptance.TestData) string {
	template := testAccAzureRMFirewallPolicyPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy_policy" "test" {
  name                = "acctest-fwpolicy-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func testAccAzureRMFirewallPolicyPolicy_complete(data acceptance.TestData) string {
	template := testAccAzureRMFirewallPolicyPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy_policy" "test" {
  name                     = "acctest-fwpolicy-Policy-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  threat_intelligence_mode = "Off"
  tags = {
    env = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMFirewallPolicyPolicy_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMFirewallPolicyPolicy_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy_policy" "import" {
  name                = azurerm_firewall_policy_policy.test.name
  resource_group_name = azurerm_firewall_policy_policy.test.resource_group_name
  location            = azurerm_firewall_policy_policy.test.location
}
`, template)
}

func testAccAzureRMFirewallPolicyPolicy_inherit(data acceptance.TestData) string {
	template := testAccAzureRMFirewallPolicyPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy_policy" "test-parent" {
  name                = "acctest-fwpolicy-Policy-%d-parent"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_firewall_policy_policy" "test" {
  name                = "acctest-fwpolicy-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  base_policy_id      = azurerm_firewall_policy_policy.test-parent.id
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFirewallPolicyPolicy_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fwpolicy-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
