package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMStorageAccountNetworkRules_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rules", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccountNetworkRules_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageAccountNetworkRules_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rules", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccountNetworkRules_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.test"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccountNetworkRules_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.test"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccountNetworkRules_empty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.test"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccountNetworkRules_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageAccountNetworkRules_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rules", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccountNetworkRules_empty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMStorageAccountNetworkRules_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_name = azurerm_storage_account.test.name

  default_action             = "Deny"
  ip_rules                   = ["127.0.0.1"]
  virtual_network_subnet_ids = [azurerm_subnet.test.id]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func testAccAzureRMStorageAccountNetworkRules_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestsubnet2%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.3.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_name = azurerm_storage_account.test.name

  default_action             = "Allow"
  ip_rules                   = ["127.0.0.2", "127.0.0.3"]
  virtual_network_subnet_ids = [azurerm_subnet.test.id, azurerm_subnet.test2.id]
  bypass                     = ["Metrics"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func testAccAzureRMStorageAccountNetworkRules_empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_name = azurerm_storage_account.test.name

  default_action             = "Deny"
  bypass                     = ["None"]
  ip_rules                   = []
  virtual_network_subnet_ids = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
