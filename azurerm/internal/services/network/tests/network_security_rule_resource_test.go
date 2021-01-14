package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetworkSecurityRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkSecurityRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkSecurityRule_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_network_security_rule"),
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityRule_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists(data.ResourceName),
					testCheckAzureRMNetworkSecurityRuleDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityRule_addingRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_updateBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists("azurerm_network_security_rule.test1"),
				),
			},

			{
				Config: testAccAzureRMNetworkSecurityRule_updateExtraRule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists("azurerm_network_security_rule.test2"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityRule_augmented(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "test1")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_augmented(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkSecurityRule_applicationSecurityGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "test1")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_applicationSecurityGroups(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMNetworkSecurityRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SecurityRuleClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		sgName := rs.Primary.Attributes["network_security_group_name"]
		sgrName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for network security rule: %q", sgName)
		}

		resp, err := client.Get(ctx, resourceGroup, sgName, sgrName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Network Security Rule %q (resource group: %q) (network security group: %q) does not exist", sgrName, sgName, resourceGroup)
			}
			return fmt.Errorf("Error retrieving Network Security Rule %q (NSG %q / Resource Group %q): %+v", sgrName, sgName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkSecurityRuleDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SecurityRuleClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		sgName := rs.Primary.Attributes["network_security_group_name"]
		sgrName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for network security rule: %s", sgName)
		}

		future, err := client.Delete(ctx, resourceGroup, sgName, sgrName)
		if err != nil {
			if !response.WasNotFound(future.Response()) {
				return fmt.Errorf("Error deleting Network Security Rule %q (NSG %q / Resource Group %q): %+v", sgrName, sgName, resourceGroup, err)
			}
		}

		return nil
	}
}

func testCheckAzureRMNetworkSecurityRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SecurityRuleClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_network_security_rule" {
			continue
		}

		sgName := rs.Primary.Attributes["network_security_group_name"]
		sgrName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, sgName, sgrName)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Network Security Rule still exists:\n%#v", resp.SecurityRulePropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMNetworkSecurityRule_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_rule" "test" {
  name                        = "test123"
  network_security_group_name = azurerm_network_security_group.test.name
  resource_group_name         = azurerm_resource_group.test.name
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMNetworkSecurityRule_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMNetworkSecurityRule_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_security_rule" "import" {
  name                        = azurerm_network_security_rule.test.name
  network_security_group_name = azurerm_network_security_rule.test.network_security_group_name
  resource_group_name         = azurerm_network_security_rule.test.resource_group_name
  priority                    = azurerm_network_security_rule.test.priority
  direction                   = azurerm_network_security_rule.test.direction
  access                      = azurerm_network_security_rule.test.access
  protocol                    = azurerm_network_security_rule.test.protocol
  source_port_range           = azurerm_network_security_rule.test.source_port_range
  destination_port_range      = azurerm_network_security_rule.test.destination_port_range
  source_address_prefix       = azurerm_network_security_rule.test.source_address_prefix
  destination_address_prefix  = azurerm_network_security_rule.test.destination_address_prefix
}
`, template)
}

func testAccAzureRMNetworkSecurityRule_updateBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test1" {
  name                = "acceptanceTestSecurityGroup2"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
}

resource "azurerm_network_security_rule" "test1" {
  name                        = "test123"
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test1.name
  network_security_group_name = azurerm_network_security_group.test1.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMNetworkSecurityRule_updateExtraRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test1" {
  name                = "acceptanceTestSecurityGroup2"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
}

resource "azurerm_network_security_rule" "test1" {
  name                        = "test123"
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test1.name
  network_security_group_name = azurerm_network_security_group.test1.name
}

resource "azurerm_network_security_rule" "test2" {
  name                        = "testing456"
  priority                    = 101
  direction                   = "Inbound"
  access                      = "Deny"
  protocol                    = "Icmp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test1.name
  network_security_group_name = azurerm_network_security_group.test1.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMNetworkSecurityRule_augmented(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test1" {
  name                = "acceptanceTestSecurityGroup2"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
}

resource "azurerm_network_security_rule" "test1" {
  name                         = "test123"
  priority                     = 100
  direction                    = "Outbound"
  access                       = "Allow"
  protocol                     = "Tcp"
  source_port_ranges           = ["10000-40000"]
  destination_port_ranges      = ["80", "443", "8080", "8190"]
  source_address_prefixes      = ["10.0.0.0/8", "192.168.0.0/16"]
  destination_address_prefixes = ["172.16.0.0/20", "8.8.8.8"]
  resource_group_name          = azurerm_resource_group.test1.name
  network_security_group_name  = azurerm_network_security_group.test1.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMNetworkSecurityRule_applicationSecurityGroups(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_security_group" "source1" {
  name                = "acctest-source1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_application_security_group" "source2" {
  name                = "acctest-source2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_application_security_group" "destination1" {
  name                = "acctest-destination1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_application_security_group" "destination2" {
  name                = "acctest-destination2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_rule" "test1" {
  name                                       = "test123"
  resource_group_name                        = azurerm_resource_group.test.name
  network_security_group_name                = azurerm_network_security_group.test.name
  priority                                   = 100
  direction                                  = "Outbound"
  access                                     = "Allow"
  protocol                                   = "Tcp"
  source_application_security_group_ids      = [azurerm_application_security_group.source1.id, azurerm_application_security_group.source2.id]
  destination_application_security_group_ids = [azurerm_application_security_group.destination1.id, azurerm_application_security_group.destination2.id]
  source_port_ranges                         = ["10000-40000"]
  destination_port_ranges                    = ["80", "443", "8080", "8190"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
