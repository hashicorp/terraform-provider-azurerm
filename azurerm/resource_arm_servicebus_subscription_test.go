package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusSubscription_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusSubscription_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists("azurerm_servicebus_subscription.test"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscription_updateEnableBatched(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMServiceBusSubscription_basic(ri, location)
	postConfig := testAccAzureRMServiceBusSubscription_updateEnableBatched(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(resourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_batched_operations", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscription_updateRequiresSession(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMServiceBusSubscription_basic(ri, location)
	postConfig := testAccAzureRMServiceBusSubscription_updateRequiresSession(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(resourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "requires_session", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscription_updateForwardTo(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMServiceBusSubscription_basic(ri, location)
	postConfig := testAccAzureRMServiceBusSubscription_updateForwardTo(ri, location)

	expectedValue := fmt.Sprintf("acctestservicebustopic-forward_to-%d", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(resourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "forward_to", expectedValue),
				),
			},
		},
	})
}

func testCheckAzureRMServiceBusSubscriptionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).serviceBusSubscriptionsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_subscription" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		topicName := rs.Primary.Attributes["topic_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, namespaceName, topicName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("ServiceBus Subscription still exists:\n%+v", resp.SBSubscriptionProperties)
		}
	}

	return nil
}

func testCheckAzureRMServiceBusSubscriptionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		subscriptionName := rs.Primary.Attributes["name"]
		topicName := rs.Primary.Attributes["topic_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for subscription: %q", topicName)
		}

		client := testAccProvider.Meta().(*ArmClient).serviceBusSubscriptionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, namespaceName, topicName, subscriptionName)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceBusSubscriptionsClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Subscription %q (resource group: %q) does not exist", subscriptionName, resourceGroup)
		}

		return nil
	}
}

const testAccAzureRMServiceBusSubscription_tfTemplate = `
resource "azurerm_resource_group" "test" {
    name     = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
    name                = "acctestservicebusnamespace-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    sku                 = "standard"
}

resource "azurerm_servicebus_topic" "test" {
    name                = "acctestservicebustopic-%d"
    namespace_name      = "${azurerm_servicebus_namespace.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_servicebus_subscription" "test" {
    name                = "acctestservicebussubscription-%d"
    namespace_name      = "${azurerm_servicebus_namespace.test.name}"
    topic_name          = "${azurerm_servicebus_topic.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    max_delivery_count  = 10
	%s
}
`

func testAccAzureRMServiceBusSubscription_basic(rInt int, location string) string {
	return fmt.Sprintf(testAccAzureRMServiceBusSubscription_tfTemplate, rInt, location, rInt, rInt, rInt, "")
}

func testAccAzureRMServiceBusSubscription_updateEnableBatched(rInt int, location string) string {
	return fmt.Sprintf(testAccAzureRMServiceBusSubscription_tfTemplate, rInt, location, rInt, rInt, rInt,
		"enable_batched_operations = true\n")
}

func testAccAzureRMServiceBusSubscription_updateRequiresSession(rInt int, location string) string {
	return fmt.Sprintf(testAccAzureRMServiceBusSubscription_tfTemplate, rInt, location, rInt, rInt, rInt,
		"requires_session = true\n")
}

func testAccAzureRMServiceBusSubscription_updateForwardTo(rInt int, location string) string {
	forwardToTf := testAccAzureRMServiceBusSubscription_tfTemplate + `

resource "azurerm_servicebus_topic" "forward_to" {
    name = "acctestservicebustopic-forward_to-%d"
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

`
	return fmt.Sprintf(forwardToTf, rInt, location, rInt, rInt, rInt,
		"forward_to = \"${azurerm_servicebus_topic.forward_to.name}\"\n", rInt)
}
