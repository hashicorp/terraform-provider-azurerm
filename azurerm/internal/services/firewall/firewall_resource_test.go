package firewall_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccFirewall_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewall_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.name", "configuration"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_configuration.0.private_ip_address"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewall_enableDNS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewall_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewall_enableDNS(data, "1.1.1.1", "8.8.8.8"),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewall_enableDNS(data, "1.1.1.1"),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewall_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewall_withManagementIp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewall_withManagementIp(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.name", "configuration"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_configuration.0.private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "management_ip_configuration.0.name", "management_configuration"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_ip_configuration.0.public_ip_address_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewall_withMultiplePublicIPs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewall_multiplePublicIps(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.name", "configuration"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_configuration.0.private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.1.name", "configuration_2"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_configuration.1.public_ip_address_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewall_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewall_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
				),
			},
			{
				Config:      testAccFirewall_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_firewall"),
			},
		},
	})
}

func TestAccFirewall_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewall_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccFirewall_withUpdatedTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewall_withZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	zones := []string{"1"}
	zonesUpdate := []string{"1", "3"}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewall_withZones(data, zones),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.0", "1"),
				),
			},
			{
				Config: testAccFirewall_withZones(data, zonesUpdate),
				Check: resource.ComposeTestCheckFunc(

					testCheckFirewallExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.0", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.1", "3"),
				),
			},
		},
	})
}

func TestAccFirewall_withoutZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewall_withoutZone(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewall_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewall_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
					testCheckFirewallDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccFirewall_withFirewallPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewall_withFirewallPolicy(data, "pol-01"),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewall_withFirewallPolicy(data, "pol-02"),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewall_inVirtualHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewall_inVirtualHub(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_hub.0.public_ip_addresses.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_hub.0.private_ip_address"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewall_inVirtualHub(data, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_hub.0.public_ip_addresses.#", "2"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_hub.0.private_ip_address"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewall_inVirtualHub(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_hub.0.public_ip_addresses.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_hub.0.private_ip_address"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckFirewallExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Firewall.AzureFirewallsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Azure Firewall: %q", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure Firewall %q (Resource Group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on azureFirewallsClient: %+v", err)
		}

		return nil
	}
}

func testCheckFirewallDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Firewall.AzureFirewallsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Azure Firewall: %q", name)
		}

		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Delete on azureFirewallsClient: %+v", err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Bad: waiting for Deletion on azureFirewallsClient: %+v", err)
		}

		return nil
	}
}

func testCheckFirewallDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Firewall.AzureFirewallsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_firewall" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Firewall still exists:\n%#v", resp.AzureFirewallPropertiesFormat)
	}

	return nil
}

func testAccFirewall_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
  threat_intel_mode = "Deny"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccFirewall_enableDNS(data acceptance.TestData, dnsServers ...string) string {
	servers := make([]string, len(dnsServers))
	for idx, server := range dnsServers {
		servers[idx] = fmt.Sprintf(`"%s"`, server)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
  threat_intel_mode = "Deny"
  dns_servers       = [%s]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, strings.Join(servers, ","))
}

func testAccFirewall_withManagementIp(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_subnet" "test_mgmt" {
  name                 = "AzureFirewallManagementSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test_mgmt" {
  name                = "acctestmgmtpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  management_ip_configuration {
    name                 = "management_configuration"
    subnet_id            = azurerm_subnet.test_mgmt.id
    public_ip_address_id = azurerm_public_ip.test_mgmt.id
  }

  threat_intel_mode = "Alert"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccFirewall_multiplePublicIps(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test_2" {
  name                = "acctestpip2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  ip_configuration {
    name                 = "configuration_2"
    public_ip_address_id = azurerm_public_ip.test_2.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccFirewall_requiresImport(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall" "import" {
  name                = azurerm_firewall.test.name
  location            = azurerm_firewall.test.location
  resource_group_name = azurerm_firewall.test.resource_group_name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
  threat_intel_mode = azurerm_firewall.test.threat_intel_mode
}
`, template)
}

func testAccFirewall_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccFirewall_withUpdatedTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccFirewall_withZones(data acceptance.TestData, zones []string) string {
	zoneString := strings.Join(zones, ",")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  zones = [%s]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, zoneString)
}

func testAccFirewall_withoutZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  zones = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccFirewall_withFirewallPolicy(data acceptance.TestData, policyName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall_policy" "test" {
  name                = "acctestfirewall-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  firewall_policy_id = azurerm_firewall_policy.test.id

  lifecycle {
    create_before_destroy = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, policyName, data.RandomInteger)
}

func testAccFirewall_inVirtualHub(data acceptance.TestData, pipCount int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%[1]d"
  location = "%s"
}

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-firewallpolicy-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-virtualwan-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-virtualhub-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_firewall" "test" {
  name                = "acctest-firewall-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "AZFW_Hub"

  virtual_hub {
    virtual_hub_id  = azurerm_virtual_hub.test.id
    public_ip_count = %[3]d
  }

  firewall_policy_id = azurerm_firewall_policy.test.id
  threat_intel_mode  = ""
}
`, data.RandomInteger, data.Locations.Primary, pipCount)
}
