package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_virtual_network_swift_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotVirtualNetworkSwiftConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotVirtualNetworkSwiftConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_virtual_network_swift_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotVirtualNetworkSwiftConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_requiresImport),
		},
	})
}

func TestAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_virtual_network_swift_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotVirtualNetworkSwiftConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotVirtualNetworkSwiftConnectionExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotVirtualNetworkSwiftConnectionExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_virtual_network_swift_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotVirtualNetworkSwiftConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlotVirtualNetworkSwiftConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotVirtualNetworkSwiftConnectionExists(data.ResourceName),
					testCheckAzureRMAppServiceSlotVirtualNetworkSwiftConnectionDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_base(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-functionapp-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lifecycle {
    ignore_changes = ["ddos_protection_plan"]
  }
}

resource "azurerm_subnet" "test1" {
  name               = "acctestSubnet1"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.1.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_subnet" "test2" {
  name               = "acctestSubnet2"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.2.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctest-ASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-FA-%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test-staging" {
  name                       = "acctest-FA-%d-staging"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_basic(data acceptance.TestData) string {
	template := testAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot_virtual_network_swift_connection" "test" {
  slot_name      = azurerm_function_app_slot.test-staging.name
  app_service_id = azurerm_function_app.test.id
  subnet_id      = azurerm_subnet.test1.id
}
`, template)
}

func testAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_update(data acceptance.TestData) string {
	template := testAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot_virtual_network_swift_connection" "test" {
  slot_name      = azurerm_function_app_slot.test-staging.name
  app_service_id = azurerm_function_app.test.id
  subnet_id      = azurerm_subnet.test2.id
}
`, template)
}

func testAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMFunctionAppSlotVirtualNetworkSwiftConnection_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot_virtual_network_swift_connection" "import" {
  slot_name      = azurerm_app_service_slot_virtual_network_swift_connection.test.slot_name
  app_service_id = azurerm_app_service_slot_virtual_network_swift_connection.test.app_service_id
  subnet_id      = azurerm_app_service_slot_virtual_network_swift_connection.test.subnet_id
}
`, template)
}
