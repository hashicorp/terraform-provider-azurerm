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

func TestAccAzureRMFirewallPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFirewallPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFirewallPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMFirewallPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMFirewallPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFirewallPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMFirewallPolicy_requiresImport),
		},
	})
}

func TestAccAzureRMFirewallPolicy_inherit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicy_inherit(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMFirewallPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Firewall.FirewallPolicyClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Firewall Policy not found: %s", resourceName)
		}

		id, err := parse.FirewallPolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Firewall Policy %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Getting on Network.FirewallPolicies: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMFirewallPolicyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Firewall.FirewallPolicyClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_firewall_policy" {
			continue
		}

		id, err := parse.FirewallPolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err == nil {
			return fmt.Errorf("Network.FirewallPolicies still exists")
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Getting on Network.FirewallPolicies: %+v", err)
		}

		return nil
	}

	return nil
}

func testAccAzureRMFirewallPolicy_basic(data acceptance.TestData) string {
	template := testAccAzureRMFirewallPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-networkfw-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func testAccAzureRMFirewallPolicy_complete(data acceptance.TestData) string {
	template := testAccAzureRMFirewallPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy" "test" {
  name                     = "acctest-networkfw-Policy-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  threat_intelligence_mode = "Off"
  threat_intelligence_allowlist {
    ip_addresses = ["1.1.1.1", "2.2.2.2"]
    fqdns        = ["foo.com", "bar.com"]
  }
  dns {
    servers       = ["1.1.1.1", "2.2.2.2"]
    proxy_enabled = true
  }
  tags = {
    env = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMFirewallPolicy_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMFirewallPolicy_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy" "import" {
  name                = azurerm_firewall_policy.test.name
  resource_group_name = azurerm_firewall_policy.test.resource_group_name
  location            = azurerm_firewall_policy.test.location
}
`, template)
}

func testAccAzureRMFirewallPolicy_inherit(data acceptance.TestData) string {
	template := testAccAzureRMFirewallPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy" "test-parent" {
  name                = "acctest-networkfw-Policy-%d-parent"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-networkfw-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  base_policy_id      = azurerm_firewall_policy.test-parent.id
  threat_intelligence_allowlist {
    ip_addresses = ["1.1.1.1", "2.2.2.2"]
    fqdns        = ["foo.com", "bar.com"]
  }
  dns {
    servers       = ["1.1.1.1", "2.2.2.2"]
    proxy_enabled = true
  }
  tags = {
    env = "Test"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFirewallPolicy_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-networkfw-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
