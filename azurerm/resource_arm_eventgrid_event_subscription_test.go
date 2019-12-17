package azurerm

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMEventGridEventSubscription_basic(t *testing.T) {
	resourceName := "azurerm_eventgrid_event_subscription.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))

	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridEventSubscription_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "event_delivery_schema", "EventGridSchema"),
					resource.TestCheckResourceAttr(resourceName, "storage_queue_endpoint.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_blob_dead_letter_destination.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "included_event_types.0", "All"),
					resource.TestCheckResourceAttr(resourceName, "retry_policy.0.max_delivery_attempts", "11"),
					resource.TestCheckResourceAttr(resourceName, "retry_policy.0.event_time_to_live", "11"),
					resource.TestCheckResourceAttr(resourceName, "labels.0", "test"),
					resource.TestCheckResourceAttr(resourceName, "labels.2", "test2"),
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

func TestAccAzureRMEventGridEventSubscription_eventhub(t *testing.T) {
	resourceName := "azurerm_eventgrid_event_subscription.test"
	ri := tf.AccRandTimeInt()

	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridEventSubscription_eventhub(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "event_delivery_schema", "CloudEventV01Schema"),
					resource.TestCheckResourceAttr(resourceName, "eventhub_endpoint.#", "1"),
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

func TestAccAzureRMEventGridEventSubscription_update(t *testing.T) {
	resourceName := "azurerm_eventgrid_event_subscription.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridEventSubscription_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "event_delivery_schema", "EventGridSchema"),
					resource.TestCheckResourceAttr(resourceName, "storage_queue_endpoint.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_blob_dead_letter_destination.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "included_event_types.0", "All"),
					resource.TestCheckResourceAttr(resourceName, "retry_policy.0.max_delivery_attempts", "11"),
					resource.TestCheckResourceAttr(resourceName, "retry_policy.0.event_time_to_live", "11"),
					resource.TestCheckResourceAttr(resourceName, "labels.0", "test"),
					resource.TestCheckResourceAttr(resourceName, "labels.2", "test2"),
				),
			},
			{
				Config: testAccAzureRMEventGridEventSubscription_update(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "included_event_types.0", "Microsoft.Storage.BlobCreated"),
					resource.TestCheckResourceAttr(resourceName, "included_event_types.1", "Microsoft.Storage.BlobDeleted"),
					resource.TestCheckResourceAttr(resourceName, "subject_filter.0.subject_ends_with", ".jpg"),
					resource.TestCheckResourceAttr(resourceName, "subject_filter.0.subject_begins_with", "test/test"),
					resource.TestCheckResourceAttr(resourceName, "retry_policy.0.max_delivery_attempts", "10"),
					resource.TestCheckResourceAttr(resourceName, "retry_policy.0.event_time_to_live", "12"),
					resource.TestCheckResourceAttr(resourceName, "labels.0", "test4"),
					resource.TestCheckResourceAttr(resourceName, "labels.2", "test6"),
				),
			},
		},
	})
}

func TestAccAzureRMEventGridEventSubscription_filter(t *testing.T) {
	resourceName := "azurerm_eventgrid_event_subscription.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))

	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridEventSubscription_filter(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "included_event_types.0", "Microsoft.Storage.BlobCreated"),
					resource.TestCheckResourceAttr(resourceName, "included_event_types.1", "Microsoft.Storage.BlobDeleted"),
					resource.TestCheckResourceAttr(resourceName, "subject_filter.0.subject_ends_with", ".jpg"),
					resource.TestCheckResourceAttr(resourceName, "subject_filter.0.subject_begins_with", "test/test"),
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

func testCheckAzureRMEventGridEventSubscriptionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.EventSubscriptionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventgrid_event_subscription" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		scope := rs.Primary.Attributes["scope"]

		resp, err := client.Get(ctx, scope, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("EventGrid Event Subscription still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMEventGridEventSubscriptionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		scope, hasScope := rs.Primary.Attributes["scope"]
		if !hasScope {
			return fmt.Errorf("Bad: no scope found in state for EventGrid Event Subscription: %s", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.EventSubscriptionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, scope, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: EventGrid Event Subscription %q (scope: %s) does not exist", name, scope)
			}

			return fmt.Errorf("Bad: Get on eventGridEventSubscriptionsClient: %s", err)
		}

		return nil
	}
}

func testAccAzureRMEventGridEventSubscription_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name = "herpderp1.vhd"

  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"

  type = "page"
  size = 5120
}

resource "azurerm_eventgrid_event_subscription" "test" {
  name  = "acctesteg-%d"
  scope = "${azurerm_resource_group.test.id}"

  storage_queue_endpoint {
    storage_account_id = "${azurerm_storage_account.test.id}"
    queue_name         = "${azurerm_storage_queue.test.name}"
  }

  storage_blob_dead_letter_destination {
    storage_account_id          = "${azurerm_storage_account.test.id}"
    storage_blob_container_name = "${azurerm_storage_container.test.name}"
  }

  retry_policy {
    event_time_to_live    = 11
    max_delivery_attempts = 11
  }

  labels = ["test", "test1", "test2"]
}
`, rInt, location, rString, rInt, rInt)
}

func testAccAzureRMEventGridEventSubscription_update(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name = "herpderp1.vhd"

  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"

  type = "page"
  size = 5120
}

resource "azurerm_eventgrid_event_subscription" "test" {
  name  = "acctesteg-%d"
  scope = "${azurerm_resource_group.test.id}"

  storage_queue_endpoint {
    storage_account_id = "${azurerm_storage_account.test.id}"
    queue_name         = "${azurerm_storage_queue.test.name}"
  }

  storage_blob_dead_letter_destination {
    storage_account_id          = "${azurerm_storage_account.test.id}"
    storage_blob_container_name = "${azurerm_storage_container.test.name}"
  }

  retry_policy {
    event_time_to_live    = 12
    max_delivery_attempts = 10
  }

  subject_filter {
    subject_begins_with = "test/test"
    subject_ends_with   = ".jpg"
  }

  included_event_types = ["Microsoft.Storage.BlobCreated", "Microsoft.Storage.BlobDeleted"]
  labels               = ["test4", "test5", "test6"]
}
`, rInt, location, rString, rInt, rInt)
}

func testAccAzureRMEventGridEventSubscription_eventhub(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventgrid_event_subscription" "test" {
  name                  = "acctesteg-%d"
  scope                 = "${azurerm_resource_group.test.id}"
  event_delivery_schema = "CloudEventV01Schema"

  eventhub_endpoint {
    eventhub_id = "${azurerm_eventhub.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMEventGridEventSubscription_filter(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}

resource "azurerm_eventgrid_event_subscription" "test" {
  name  = "acctesteg-%d"
  scope = "${azurerm_resource_group.test.id}"

  storage_queue_endpoint {
    storage_account_id = "${azurerm_storage_account.test.id}"
    queue_name         = "${azurerm_storage_queue.test.name}"
  }

  included_event_types = ["Microsoft.Storage.BlobCreated", "Microsoft.Storage.BlobDeleted"]

  subject_filter {
    subject_begins_with = "test/test"
    subject_ends_with   = ".jpg"
  }
}
`, rInt, location, rString, rInt, rInt)
}
