package eventgrid_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMEventGridEventSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_event_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridEventSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "event_delivery_schema", "EventGridSchema"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMEventGridEventSubscription_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_event_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridEventSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMEventGridEventSubscription_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_eventgrid_event_subscription"),
			},
		},
	})
}

func TestAccAzureRMEventGridEventSubscription_eventHubID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_event_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridEventSubscription_eventHubID(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "event_delivery_schema", "CloudEventSchemaV1_0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "eventhub_endpoint_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMEventGridEventSubscription_serviceBusQueueID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_event_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridEventSubscription_serviceBusQueueID(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "event_delivery_schema", "CloudEventSchemaV1_0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "service_bus_queue_endpoint_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMEventGridEventSubscription_serviceBusTopicID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_event_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridEventSubscription_serviceBusTopicID(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "event_delivery_schema", "CloudEventSchemaV1_0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "service_bus_topic_endpoint_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMEventGridEventSubscription_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_event_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridEventSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "event_delivery_schema", "EventGridSchema"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_queue_endpoint.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_blob_dead_letter_destination.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "included_event_types.0", "Microsoft.Resources.ResourceWriteSuccess"),
					resource.TestCheckResourceAttr(data.ResourceName, "retry_policy.0.max_delivery_attempts", "11"),
					resource.TestCheckResourceAttr(data.ResourceName, "retry_policy.0.event_time_to_live", "11"),
					resource.TestCheckResourceAttr(data.ResourceName, "labels.0", "test"),
					resource.TestCheckResourceAttr(data.ResourceName, "labels.2", "test2"),
				),
			},
			{
				Config: testAccAzureRMEventGridEventSubscription_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "included_event_types.0", "Microsoft.Storage.BlobCreated"),
					resource.TestCheckResourceAttr(data.ResourceName, "included_event_types.1", "Microsoft.Storage.BlobDeleted"),
					resource.TestCheckResourceAttr(data.ResourceName, "subject_filter.0.subject_ends_with", ".jpg"),
					resource.TestCheckResourceAttr(data.ResourceName, "subject_filter.0.subject_begins_with", "test/test"),
					resource.TestCheckResourceAttr(data.ResourceName, "retry_policy.0.max_delivery_attempts", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "retry_policy.0.event_time_to_live", "12"),
					resource.TestCheckResourceAttr(data.ResourceName, "labels.0", "test4"),
					resource.TestCheckResourceAttr(data.ResourceName, "labels.2", "test6"),
				),
			},
		},
	})
}

func TestAccAzureRMEventGridEventSubscription_filter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_event_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridEventSubscription_filter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "included_event_types.0", "Microsoft.Storage.BlobCreated"),
					resource.TestCheckResourceAttr(data.ResourceName, "included_event_types.1", "Microsoft.Storage.BlobDeleted"),
					resource.TestCheckResourceAttr(data.ResourceName, "subject_filter.0.subject_ends_with", ".jpg"),
					resource.TestCheckResourceAttr(data.ResourceName, "subject_filter.0.subject_begins_with", "test/test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMEventGridEventSubscription_advancedFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_event_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridEventSubscription_advancedFilter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridEventSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.bool_equals.0.key", "subject"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.bool_equals.0.value", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.number_greater_than.0.key", "data.metadataVersion"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.number_greater_than.0.value", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.number_greater_than_or_equals.0.key", "data.contentLength"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.number_greater_than_or_equals.0.value", "42"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.number_less_than.0.key", "data.contentLength"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.number_less_than.0.value", "42.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.number_less_than_or_equals.0.key", "data.metadataVersion"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.number_less_than_or_equals.0.value", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.number_in.0.key", "data.contentLength"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.number_in.0.values.0", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.number_not_in.0.key", "data.contentLength"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.number_not_in.0.values.0", "5"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.string_begins_with.0.key", "subject"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.string_begins_with.0.values.0", "foo"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.string_ends_with.0.key", "subject"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.string_ends_with.0.values.0", "bar"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.string_contains.0.key", "data.contentType"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.string_contains.0.values.0", "application"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.string_in.0.key", "data.blobType"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.string_in.0.values.0", "Block"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.string_not_in.0.key", "data.blobType"),
					resource.TestCheckResourceAttr(data.ResourceName, "advanced_filter.0.string_not_in.0.values.0", "Page"),
				),
			},
			data.ImportStep(),
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
		client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.EventSubscriptionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMEventGridEventSubscription_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%d"
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
  storage_account_name = "${azurerm_storage_account.test.name}"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name = "herpderp1.vhd"

  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"

  type = "Page"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMEventGridEventSubscription_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMEventGridEventSubscription_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_event_subscription" "import" {
  name  = azurerm_eventgrid_event_subscription.test.name
  scope = azurerm_eventgrid_event_subscription.test.scope
}
`, template)
}

func testAccAzureRMEventGridEventSubscription_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%d"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name = "herpderp1.vhd"

  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name

  type = "Page"
  size = 5120
}

resource "azurerm_eventgrid_event_subscription" "test" {
  name  = "acctest-eg-%d"
  scope = azurerm_resource_group.test.id

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.test.id
    queue_name         = azurerm_storage_queue.test.name
  }

  storage_blob_dead_letter_destination {
    storage_account_id          = azurerm_storage_account.test.id
    storage_blob_container_name = azurerm_storage_container.test.name
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMEventGridEventSubscription_eventHubID(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventgrid_event_subscription" "test" {
  name                  = "acctest-eg-%d"
  scope                 = azurerm_resource_group.test.id
  event_delivery_schema = "CloudEventSchemaV1_0"

  eventhub_endpoint_id = azurerm_eventhub.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMEventGridEventSubscription_serviceBusQueueID(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}
resource "azurerm_servicebus_queue" "test" {
  name                = "acctestservicebusqueue-%d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.example.name
  enable_partitioning = true
}
resource "azurerm_eventgrid_event_subscription" "test" {
  name                          = "acctest-eg-%d"
  scope                         = azurerm_resource_group.test.id
  event_delivery_schema         = "CloudEventSchemaV1_0"
  service_bus_queue_endpoint_id = azurerm_servicebus_queue.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMEventGridEventSubscription_serviceBusTopicID(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%d"
  location = "%s"
}
resource "azurerm_servicebus_namespace" "example" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}
resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.example.name
  enable_partitioning = true
}
resource "azurerm_eventgrid_event_subscription" "test" {
  name                          = "acctest-eg-%d"
  scope                         = azurerm_resource_group.test.id
  event_delivery_schema         = "CloudEventSchemaV1_0"
  service_bus_topic_endpoint_id = azurerm_servicebus_topic.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMEventGridEventSubscription_filter(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%d"
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
  storage_account_name = "${azurerm_storage_account.test.name}"
}

resource "azurerm_eventgrid_event_subscription" "test" {
  name  = "acctest-eg-%d"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMEventGridEventSubscription_advancedFilter(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%d"
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
  storage_account_name = "${azurerm_storage_account.test.name}"
}

resource "azurerm_eventgrid_event_subscription" "test" {
  name  = "acctesteg-%d"
  scope = "${azurerm_storage_account.test.id}"

  storage_queue_endpoint {
    storage_account_id = "${azurerm_storage_account.test.id}"
    queue_name         = "${azurerm_storage_queue.test.name}"
  }

  advanced_filter {
    bool_equals {
      key   = "subject"
      value = true
    }
    number_greater_than {
      key   = "data.metadataVersion"
      value = 1
    }
    number_greater_than_or_equals {
      key   = "data.contentLength"
      value = 42.0
    }
    number_less_than {
      key   = "data.contentLength"
      value = 42.1
    }
    number_less_than_or_equals {
      key   = "data.metadataVersion"
      value = 2
    }
    number_in {
      key    = "data.contentLength"
      values = [0, 1, 1, 2, 3]
    }
    number_not_in {
      key    = "data.contentLength"
      values = [5, 8, 13, 21, 34]
    }
    string_begins_with {
      key    = "subject"
      values = ["foo"]
    }
    string_ends_with {
      key    = "subject"
      values = ["bar"]
    }
    string_contains {
      key    = "data.contentType"
      values = ["application", "octet-stream"]
    }
    string_in {
      key    = "data.blobType"
      values = ["Block"]
    }
    string_not_in {
      key    = "data.blobType"
      values = ["Page"]
    }
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}
