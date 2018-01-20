package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusTopicAuthorizationRule_listen(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusTopicAuthorizationRule_listen(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicAuthorizationRuleExists("azurerm_servicebus_topic_authorization_rule.test"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusTopicAuthorizationRule_send(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusTopicAuthorizationRule_send(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicAuthorizationRuleExists("azurerm_servicebus_topic_authorization_rule.test"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusTopicAuthorizationRule_readwrite(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusTopicAuthorizationRule_readWrite(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicAuthorizationRuleExists("azurerm_servicebus_topic_authorization_rule.test"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusTopicAuthorizationRule_manage(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusTopicAuthorizationRule_manage(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicAuthorizationRuleExists("azurerm_servicebus_topic_authorization_rule.test"),
				),
			},
		},
	})
}

func testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).serviceBusTopicsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testCheckAzureRMServiceBusTopicAuthorizationRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		topicName := rs.Primary.Attributes["topic_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for ServiceBus Topic: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).serviceBusTopicsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testAccAzureRMServiceBusTopicAuthorizationRule_listen(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  name                = "acctestservicebustopicrule-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name       = "${azurerm_servicebus_topic.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  listen              = true
  send                = false
  manage              = false
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMServiceBusTopicAuthorizationRule_send(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  name                = "acctestservicebustopicrule-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name       = "${azurerm_servicebus_topic.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  listen              = false
  send                = true
  manage              = false
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMServiceBusTopicAuthorizationRule_readWrite(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  name                = "acctestservicebustopicrule-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name       = "${azurerm_servicebus_topic.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  listen              = true
  send                = true
  manage              = false
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMServiceBusTopicAuthorizationRule_manage(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  name                = "acctestservicebustopicrule-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name       = "${azurerm_servicebus_topic.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  listen              = true
  send                = true
  manage              = true
}
`, rInt, location, rInt, rInt, rInt)
}
