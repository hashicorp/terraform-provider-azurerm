package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusSubscription_basic(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscription_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscription_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_servicebus_subscription.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscription_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMServiceBusSubscription_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_servicebus_subscription"),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscription_defaultTtl(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMServiceBusSubscription_withDefaultTtl(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_servicebus_subscription.test", "default_message_ttl", "PT1H"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscription_updateEnableBatched(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMServiceBusSubscription_basic(ri, location)
	postConfig := testAccAzureRMServiceBusSubscription_updateEnableBatched(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscription_updateRequiresSession(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMServiceBusSubscription_basic(ri, location)
	postConfig := testAccAzureRMServiceBusSubscription_updateRequiresSession(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscription_updateForwardTo(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMServiceBusSubscription_basic(ri, location)
	postConfig := testAccAzureRMServiceBusSubscription_updateForwardTo(ri, location)

	expectedValue := fmt.Sprintf("acctestservicebustopic-forward_to-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscription_updateForwardDeadLetteredMessagesTo(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMServiceBusSubscription_basic(ri, location)
	postConfig := testAccAzureRMServiceBusSubscription_updateForwardDeadLetteredMessagesTo(ri, location)

	expectedValue := fmt.Sprintf("acctestservicebustopic-forward_dl_messages_to-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
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
					resource.TestCheckResourceAttr(resourceName, "forward_dead_lettered_messages_to", expectedValue),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMServiceBusSubscriptionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.SubscriptionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckAzureRMServiceBusSubscriptionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		subscriptionName := rs.Primary.Attributes["name"]
		topicName := rs.Primary.Attributes["topic_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Subscription: %q", topicName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.SubscriptionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
    sku                 = "Standard"
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

func testAccAzureRMServiceBusSubscription_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`

%s

resource "azurerm_servicebus_subscription" "import" {
    name                = "${azurerm_servicebus_subscription.test.name}"
    namespace_name      = "${azurerm_servicebus_subscription.test.namespace_name}"
    topic_name          = "${azurerm_servicebus_subscription.test.topic_name}"
    resource_group_name = "${azurerm_servicebus_subscription.test.resource_group_name}"
    max_delivery_count  = "${azurerm_servicebus_subscription.test.max_delivery_count}"
}
`, testAccAzureRMServiceBusSubscription_basic(rInt, location))
}

func testAccAzureRMServiceBusSubscription_withDefaultTtl(rInt int, location string) string {
	return fmt.Sprintf(testAccAzureRMServiceBusSubscription_tfTemplate, rInt, location, rInt, rInt, rInt,
		"default_message_ttl = \"PT1H\"\n")
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

func testAccAzureRMServiceBusSubscription_updateForwardDeadLetteredMessagesTo(rInt int, location string) string {
	forwardToTf := testAccAzureRMServiceBusSubscription_tfTemplate + `

resource "azurerm_servicebus_topic" "forward_dl_messages_to" {
    name = "acctestservicebustopic-forward_dl_messages_to-%d"
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

`
	return fmt.Sprintf(forwardToTf, rInt, location, rInt, rInt, rInt,
		"forward_dead_lettered_messages_to = \"${azurerm_servicebus_topic.forward_dl_messages_to.name}\"\n", rInt)
}
