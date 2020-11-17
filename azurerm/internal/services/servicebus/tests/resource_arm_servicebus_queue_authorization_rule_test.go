package tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusQueueAuthorizationRule_listen(t *testing.T) {
	testAccAzureRMServiceBusQueueAuthorizationRule(t, true, false, false)
}

func TestAccAzureRMServiceBusQueueAuthorizationRule_send(t *testing.T) {
	testAccAzureRMServiceBusQueueAuthorizationRule(t, false, true, false)
}

func TestAccAzureRMServiceBusQueueAuthorizationRule_listensend(t *testing.T) {
	testAccAzureRMServiceBusQueueAuthorizationRule(t, true, true, false)
}

func TestAccAzureRMServiceBusQueueAuthorizationRule_manage(t *testing.T) {
	testAccAzureRMServiceBusQueueAuthorizationRule(t, true, true, true)
}

func testAccAzureRMServiceBusQueueAuthorizationRule(t *testing.T, listen, send, manage bool) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueueAuthorizationRule_base(data, listen, send, manage),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "queue_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", strconv.FormatBool(listen)),
					resource.TestCheckResourceAttr(data.ResourceName, "send", strconv.FormatBool(send)),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", strconv.FormatBool(manage)),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusQueueAuthorizationRule_rightsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueueAuthorizationRule_base(data, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "false"),
				),
			},
			{
				Config: testAccAzureRMServiceBusQueueAuthorizationRule_base(data, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusQueueAuthorizationRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueueAuthorizationRule_base(data, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "false"),
				),
			},
			{
				Config:      testAccAzureRMServiceBusQueueAuthorizationRule_requiresImport(data, true, false, false),
				ExpectError: acceptance.RequiresImportError("azurerm_servicebus_queue_authorization_rule"),
			},
		},
	})
}

func testCheckAzureRMServiceBusQueueAuthorizationRuleDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.QueuesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_queue_authorization_rule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		queueName := rs.Primary.Attributes["queue_name"]

		resp, err := conn.GetAuthorizationRule(ctx, resourceGroup, namespaceName, queueName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testCheckAzureRMServiceBusQueueAuthorizationRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.QueuesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		queueName := rs.Primary.Attributes["queue_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]

		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for ServiceBus Queue Authorization Rule: %s", name)
		}

		resp, err := conn.GetAuthorizationRule(ctx, resourceGroup, namespaceName, queueName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: ServiceBus Queue Authorization Rule %q (Queue %q / Namespace %q / Resource Group: %q) does not exist", name, queueName, namespaceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on ServiceBus Queue: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMServiceBusQueueAuthorizationRule_base(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name

  enable_partitioning = true
}

resource "azurerm_servicebus_queue_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  queue_name          = azurerm_servicebus_queue.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = %[3]t
  send   = %[4]t
  manage = %[5]t
}
`, data.RandomInteger, data.Locations.Primary, listen, send, manage)
}

func testAccAzureRMServiceBusQueueAuthorizationRule_requiresImport(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_queue_authorization_rule" "import" {
  name                = azurerm_servicebus_queue_authorization_rule.test.name
  namespace_name      = azurerm_servicebus_queue_authorization_rule.test.namespace_name
  queue_name          = azurerm_servicebus_queue_authorization_rule.test.queue_name
  resource_group_name = azurerm_servicebus_queue_authorization_rule.test.resource_group_name

  listen = azurerm_servicebus_queue_authorization_rule.test.listen
  send   = azurerm_servicebus_queue_authorization_rule.test.send
  manage = azurerm_servicebus_queue_authorization_rule.test.manage
}
`, testAccAzureRMServiceBusQueueAuthorizationRule_base(data, listen, send, manage))
}
