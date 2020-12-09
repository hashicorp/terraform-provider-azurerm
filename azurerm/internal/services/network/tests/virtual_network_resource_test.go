package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMVirtualNetwork_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet.1472110187.id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualNetwork_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetwork_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualNetwork_basicUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet.1472110187.id"),
				),
			},
			{
				Config: testAccAzureRMVirtualNetwork_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.#", "2"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet.1472110187.id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualNetwork_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualNetwork_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_network"),
			},
		},
	})
}

func TestAccAzureRMVirtualNetwork_ddosProtectionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetwork_ddosProtectionPlan(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ddos_protection_plan.0.enable", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ddos_protection_plan.0.id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualNetwork_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
					testCheckAzureRMVirtualNetworkDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMVirtualNetwork_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetwork_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet.1472110187.id"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccAzureRMVirtualNetwork_withTagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet.1472110187.id"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetwork_deleteSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualNetwork_noSubnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualNetwork_bgpCommunity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualNetwork_bgpCommunity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualNetwork_vmProtection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualNetwork_vmProtection(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualNetwork_vmProtection(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMVirtualNetworkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VnetClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		virtualNetworkName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual network: %s", virtualNetworkName)
		}

		// Ensure resource group/virtual network combination exists in API
		resp, err := client.Get(ctx, resourceGroup, virtualNetworkName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on vnetClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual Network %q (resource group: %q) does not exist", virtualNetworkName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMVirtualNetworkDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VnetClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		virtualNetworkName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual network: %s", virtualNetworkName)
		}

		// Ensure resource group/virtual network combination exists in API
		future, err := client.Delete(ctx, resourceGroup, virtualNetworkName)
		if err != nil {
			return fmt.Errorf("Error deleting Virtual Network %q (RG %q): %+v", virtualNetworkName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Virtual Network %q (RG %q): %+v", virtualNetworkName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualNetworkDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VnetClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_network" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Virtual Network still exists:\n%#v", resp.VirtualNetworkPropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMVirtualNetwork_basic(data acceptance.TestData) string {
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

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMVirtualNetwork_complete(data acceptance.TestData) string {
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
  address_space       = ["10.0.0.0/16", "10.10.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_servers         = ["10.7.7.2", "10.7.7.7", "10.7.7.1", ]

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }

  subnet {
    name           = "subnet2"
    address_prefix = "10.10.1.0/24"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMVirtualNetwork_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualNetwork_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "import" {
  name                = azurerm_virtual_network.test.name
  location            = azurerm_virtual_network.test.location
  resource_group_name = azurerm_virtual_network.test.resource_group_name
  address_space       = ["10.0.0.0/16"]

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}
`, template)
}

func testAccAzureRMVirtualNetwork_ddosProtectionPlan(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ddos_protection_plan {
    id     = azurerm_network_ddos_protection_plan.test.id
    enable = true
  }

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualNetwork_withTags(data acceptance.TestData) string {
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

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMVirtualNetwork_withTagsUpdated(data acceptance.TestData) string {
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

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMVirtualNetwork_noSubnet(data acceptance.TestData) string {
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
  subnet              = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMVirtualNetwork_bgpCommunity(data acceptance.TestData) string {
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

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }

  bgp_community = "12076:20000"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMVirtualNetwork_vmProtection(data acceptance.TestData, enabled bool) string {
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

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }

  vm_protection_enabled = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enabled)
}
