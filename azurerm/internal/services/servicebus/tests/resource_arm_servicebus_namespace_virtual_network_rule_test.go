package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusNamespaceVirtualNetworkRule_create(t *testing.T) {
	testAccAzureRMServiceBusNamespaceVirtualNetworkRule(t)
}

func testAccAzureRMServiceBusNamespaceVirtualNetworkRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_virtual_network_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespaceVirtualNetworkRule_base(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceVirtualNetworkRuleExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusNamespaceVirtualNetworkRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_virtual_network_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespaceVirtualNetworkRule_base(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceVirtualNetworkRuleExists(data.ResourceName)),
			},
			{
				Config:      testAccAzureRMServiceBusNamespaceVirtualNetworkRule_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_servicebus_namespace_virtual_network_rule"),
			},
		},
	})
}

func testCheckAzureRMServiceBusNamespaceVirtualNetworkRuleDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.NamespacesClientPreview
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_namespace_virtual_network_rule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetVirtualNetworkRule(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testCheckAzureRMServiceBusNamespaceVirtualNetworkRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.NamespacesClientPreview
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for ServiceBus Namespace: %s", name)
		}

		resp, err := conn.GetVirtualNetworkRule(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: ServiceBus Namespace Virtual Network Rule %q (namespace %s / resource group: %s) does not exist", name, namespaceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on ServiceBus Namespace: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMServiceBusNamespaceVirtualNetworkRule_base(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%[1]d"
  address_space       = ["10.7.29.0/29"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.ServiceBus"]
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Premium"
  capacity            = "1"
}

resource "azurerm_servicebus_namespace_virtual_network_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMServiceBusNamespaceVirtualNetworkRule_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_virtual_network_rule" "import" {
  name                = azurerm_servicebus_namespace_virtual_network_rule.test.name
  namespace_name      = azurerm_servicebus_namespace_virtual_network_rule.test.namespace_name
  resource_group_name = azurerm_servicebus_namespace_virtual_network_rule.test.resource_group_name
  subnet_id           = azurerm_servicebus_namespace_virtual_network_rule.test.subnet_id
}
`, testAccAzureRMServiceBusNamespaceVirtualNetworkRule_base(data))
}
