package kusto_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccKustoCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKustoCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccKustoCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKustoCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_disk_encryption", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_streaming_ingest", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_purge", "false"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccKustoCluster_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_disk_encryption", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_streaming_ingest", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_purge", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccKustoCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_disk_encryption", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_streaming_ingest", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_purge", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccKustoCluster_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKustoCluster_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.label", "test"),
				),
			},
			{
				Config: testAccKustoCluster_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.label", "test1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "prod"),
				),
			},
		},
	})
}

func TestAccKustoCluster_sku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKustoCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Dev(No SLA)_Standard_D11_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
				),
			},
			{
				Config: testAccKustoCluster_skuUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_D11_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "2"),
				),
			},
		},
	})
}

func TestAccKustoCluster_zones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKustoCluster_withZones(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.0", "1"),
				),
			},
		},
	})
}

func TestAccKustoCluster_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKustoCluster_identitySystemAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.identity_ids.#", "0"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccKustoCluster_vnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKustoCluster_vnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_network_configuration.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_configuration.0.subnet_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_configuration.0.engine_public_ip_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_configuration.0.data_management_public_ip_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccKustoCluster_languageExtensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKustoCluster_languageExtensions(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "language_extensions.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "language_extensions.0", "PYTHON"),
					resource.TestCheckResourceAttr(data.ResourceName, "language_extensions.1", "R"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccKustoCluster_languageExtensionsRemove(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "language_extensions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "language_extensions.0", "R"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccKustoCluster_optimizedAutoScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKustoCluster_optimizedAutoScale(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "optimized_auto_scale.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "optimized_auto_scale.0.minimum_instances", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "optimized_auto_scale.0.maximum_instances", "3"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccKustoCluster_optimizedAutoScaleUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "optimized_auto_scale.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "optimized_auto_scale.0.minimum_instances", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "optimized_auto_scale.0.maximum_instances", "4"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccKustoCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccKustoCluster_engineV3(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKustoCluster_engineV3(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccKustoCluster_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccKustoCluster_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  tags = {
    label = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccKustoCluster_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  tags = {
    label = "test1"
    ENV   = "prod"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccKustoCluster_skuUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_D11_v2"
    capacity = 2
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccKustoCluster_withZones(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  zones = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccKustoCluster_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                    = "acctestkc%s"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  enable_disk_encryption  = true
  enable_streaming_ingest = true
  enable_purge            = true

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccKustoCluster_identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccKustoCluster_languageExtensions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  language_extensions = ["PYTHON", "R"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccKustoCluster_languageExtensionsRemove(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  language_extensions = ["R"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccKustoCluster_optimizedAutoScale(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name = "Standard_D11_v2"
  }

  optimized_auto_scale {
    minimum_instances = 2
    maximum_instances = 3
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccKustoCluster_optimizedAutoScaleUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name = "Standard_D11_v2"
  }

  optimized_auto_scale {
    minimum_instances = 3
    maximum_instances = 4
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccKustoCluster_vnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestkc%s-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestkc%s-subnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestkc%s-nsg"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_rule" "test_allow_management_inbound" {
  name                        = "AllowAzureDataExplorerManagement"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "443"
  source_address_prefix       = "AzureDataExplorerManagement"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_public_ip" "engine_pip" {
  name                = "acctestkc%s-engine-pip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_public_ip" "management_pip" {
  name                = "acctestkc%s-management-pip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
  allocation_method   = "Static"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  virtual_network_configuration {
    subnet_id                    = azurerm_subnet.test.id
    engine_public_ip_id          = azurerm_public_ip.engine_pip.id
    data_management_public_ip_id = azurerm_public_ip.management_pip.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString)
}

func testAccKustoCluster_engineV3(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
  engine = "V3"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testCheckKustoClusterDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.ClustersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kusto_cluster" {
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

		return nil
	}

	return nil
}

func testCheckKustoClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.ClustersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		kustoCluster := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Kusto Cluster: %s", kustoCluster)
		}

		resp, err := client.Get(ctx, resourceGroup, kustoCluster)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Kusto Cluster %q (resource group: %q) does not exist", kustoCluster, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on ClustersClient: %+v", err)
		}

		return nil
	}
}
