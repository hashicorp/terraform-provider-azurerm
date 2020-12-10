package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusNamespaceNetworkRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespaceNetworkRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceNetworkRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusNamespaceNetworkRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespaceNetworkRule_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceNetworkRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusNamespaceNetworkRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespaceNetworkRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceNetworkRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMServiceBusNamespaceNetworkRule_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceNetworkRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMServiceBusNamespaceNetworkRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceNetworkRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusNamespaceNetworkRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespaceNetworkRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceNetworkRuleExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMServiceBusNamespaceNetworkRule_requiresImport),
		},
	})
}

func testCheckAzureRMServiceBusNamespaceNetworkRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.NamespacesClientPreview
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Service Bus Namespace Network Rule Set not found: %s", resourceName)
		}

		id, err := parse.NamespaceNetworkRuleSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.GetNetworkRuleSet(ctx, id.ResourceGroup, id.NamespaceName); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Service Bus Namespace Network Rule Set (Namespace %q / Resource Group %q) does not exist", id.NamespaceName, id.ResourceGroup)
			}
			return fmt.Errorf("failed to GetNetworkRuleSet on ServiceBus.NamespacesClientPreview: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMServiceBusNamespaceNetworkRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.NamespacesClientPreview
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_namespace_network_rule_set" {
			continue
		}

		id, err := parse.NamespaceNetworkRuleSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		// this resource cannot be deleted, instead, we check if this setting was set back to empty
		resp, err := client.GetNetworkRuleSet(ctx, id.ResourceGroup, id.NamespaceName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("failed to GetNetworkRuleSet on ServiceBus.NamespacesClientPreview: %+v", err)
			}
			return nil
		}

		if !servicebus.CheckNetworkRuleNullified(resp) {
			return fmt.Errorf("the Service Bus Namespace Network Rule Set (Namespace %q / Resource Group %q) still exists", id.NamespaceName, id.ResourceGroup)
		}
	}

	return nil
}

func testAccAzureRMServiceBusNamespaceNetworkRule_basic(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusNamespaceNetworkRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_network_rule_set" "test" {
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  default_action = "Deny"

  network_rules {
    subnet_id                            = azurerm_subnet.test.id
    ignore_missing_vnet_service_endpoint = false
  }
}
`, template)
}

func testAccAzureRMServiceBusNamespaceNetworkRule_complete(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusNamespaceNetworkRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_network_rule_set" "test" {
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  default_action = "Deny"

  network_rules {
    subnet_id                            = azurerm_subnet.test.id
    ignore_missing_vnet_service_endpoint = false
  }

  ip_rules = ["1.1.1.1"]
}
`, template)
}

func testAccAzureRMServiceBusNamespaceNetworkRule_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sb-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-sb-namespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Premium"

  capacity = 1
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["172.17.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "test" {
  name                 = "${azurerm_virtual_network.test.name}-default"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "172.17.0.0/24"

  service_endpoints = ["Microsoft.ServiceBus"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMServiceBusNamespaceNetworkRule_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusNamespaceNetworkRule_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_network_rule_set" "import" {
  namespace_name      = azurerm_servicebus_namespace_network_rule_set.test.namespace_name
  resource_group_name = azurerm_servicebus_namespace_network_rule_set.test.resource_group_name
}
`, template)
}
