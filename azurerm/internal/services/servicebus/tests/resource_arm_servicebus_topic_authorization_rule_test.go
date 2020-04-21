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

func TestAccAzureRMServiceBusTopicAuthorizationRule_listen(t *testing.T) {
	testAccAzureRMServiceBusTopicAuthorizationRule(t, true, false, false)
}

func TestAccAzureRMServiceBusTopicAuthorizationRule_send(t *testing.T) {
	testAccAzureRMServiceBusTopicAuthorizationRule(t, false, true, false)
}

func TestAccAzureRMServiceBusTopicAuthorizationRule_listensend(t *testing.T) {
	testAccAzureRMServiceBusTopicAuthorizationRule(t, true, true, false)
}

func TestAccAzureRMServiceBusTopicAuthorizationRule_manage(t *testing.T) {
	testAccAzureRMServiceBusTopicAuthorizationRule(t, true, true, true)
}

func testAccAzureRMServiceBusTopicAuthorizationRule(t *testing.T, listen, send, manage bool) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusTopicAuthorizationRule_base(data, listen, send, manage),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "namespace_name"),
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

func TestAccAzureRMServiceBusTopicAuthorizationRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusTopicAuthorizationRule_base(data, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "false"),
				),
			},
			{
				Config:      testAccAzureRMServiceBusTopicAuthorizationRule_requiresImport(data, true, false, false),
				ExpectError: acceptance.RequiresImportError("azurerm_servicebus_topic_authorization_rule"),
			},
		},
	})
}

func TestAccAzureRMServiceBusTopicAuthorizationRule_rightsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusTopicAuthorizationRule_base(data, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "false"),
				),
			},
			{
				Config: testAccAzureRMServiceBusTopicAuthorizationRule_base(data, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicAuthorizationRuleExists(data.ResourceName),
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

func testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.TopicsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_topic_authorization_rule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		topicName := rs.Primary.Attributes["topic_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetAuthorizationRule(ctx, resourceGroup, namespaceName, topicName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testCheckAzureRMServiceBusTopicAuthorizationRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.TopicsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		topicName := rs.Primary.Attributes["topic_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for ServiceBus Topic: %s", name)
		}

		resp, err := conn.GetAuthorizationRule(ctx, resourceGroup, namespaceName, topicName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: ServiceBus Topic Authorization Rule %q (topic %s, namespace %s / resource group: %s) does not exist", name, topicName, namespaceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on serviceBusTopicsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMServiceBusTopicAuthorizationRule_base(data acceptance.TestData, listen, send, manage bool) string {
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

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%[1]d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  topic_name          = azurerm_servicebus_topic.test.name

  listen = %[3]t
  send   = %[4]t
  manage = %[5]t
}
`, data.RandomInteger, data.Locations.Primary, listen, send, manage)
}

func testAccAzureRMServiceBusTopicAuthorizationRule_requiresImport(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_topic_authorization_rule" "import" {
  name                = azurerm_servicebus_topic_authorization_rule.test.name
  namespace_name      = azurerm_servicebus_topic_authorization_rule.test.namespace_name
  resource_group_name = azurerm_servicebus_topic_authorization_rule.test.resource_group_name
  topic_name          = azurerm_servicebus_topic_authorization_rule.test.topic_name

  listen = azurerm_servicebus_topic_authorization_rule.test.listen
  send   = azurerm_servicebus_topic_authorization_rule.test.send
  manage = azurerm_servicebus_topic_authorization_rule.test.manage
}
`, testAccAzureRMServiceBusTopicAuthorizationRule_base(data, listen, send, manage))
}
