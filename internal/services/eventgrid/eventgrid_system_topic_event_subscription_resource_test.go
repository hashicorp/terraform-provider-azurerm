// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/eventsubscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type EventGridSystemTopicEventSubscriptionResource struct{}

func TestAccEventGridSystemTopicEventSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("event_delivery_schema").HasValue("EventGridSchema"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopicEventSubscription_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_eventgrid_system_topic_event_subscription"),
		},
	})
}

func TestAccEventGridSystemTopicEventSubscription_eventHubID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.eventHubID(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("event_delivery_schema").HasValue("CloudEventSchemaV1_0"),
				check.That(data.ResourceName).Key("eventhub_endpoint_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopicEventSubscription_serviceBusQueueID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceBusQueueID(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("event_delivery_schema").HasValue("CloudEventSchemaV1_0"),
				check.That(data.ResourceName).Key("service_bus_queue_endpoint_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopicEventSubscription_serviceBusTopicID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceBusTopicID(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("event_delivery_schema").HasValue("CloudEventSchemaV1_0"),
				check.That(data.ResourceName).Key("service_bus_topic_endpoint_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopicEventSubscription_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("event_delivery_schema").HasValue("EventGridSchema"),
				check.That(data.ResourceName).Key("storage_queue_endpoint.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_blob_dead_letter_destination.#").HasValue("1"),
				check.That(data.ResourceName).Key("included_event_types.0").HasValue("Microsoft.Resources.ResourceWriteSuccess"),
				check.That(data.ResourceName).Key("retry_policy.0.max_delivery_attempts").HasValue("11"),
				check.That(data.ResourceName).Key("retry_policy.0.event_time_to_live").HasValue("11"),
				check.That(data.ResourceName).Key("labels.0").HasValue("test"),
				check.That(data.ResourceName).Key("labels.2").HasValue("test2"),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("included_event_types.0").HasValue("Microsoft.Storage.BlobCreated"),
				check.That(data.ResourceName).Key("included_event_types.1").HasValue("Microsoft.Storage.BlobDeleted"),
				check.That(data.ResourceName).Key("storage_queue_endpoint.0.queue_message_time_to_live_in_seconds").HasValue("3600"),
				check.That(data.ResourceName).Key("subject_filter.0.subject_ends_with").HasValue(".jpg"),
				check.That(data.ResourceName).Key("subject_filter.0.subject_begins_with").HasValue("test/test"),
				check.That(data.ResourceName).Key("retry_policy.0.max_delivery_attempts").HasValue("10"),
				check.That(data.ResourceName).Key("retry_policy.0.event_time_to_live").HasValue("12"),
				check.That(data.ResourceName).Key("labels.0").HasValue("test4"),
				check.That(data.ResourceName).Key("labels.2").HasValue("test6"),
			),
		},
	})
}

func TestAccEventGridSystemTopicEventSubscription_filter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.filter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("included_event_types.0").HasValue("Microsoft.Storage.BlobCreated"),
				check.That(data.ResourceName).Key("included_event_types.1").HasValue("Microsoft.Storage.BlobDeleted"),
				check.That(data.ResourceName).Key("subject_filter.0.subject_ends_with").HasValue(".jpg"),
				check.That(data.ResourceName).Key("subject_filter.0.subject_begins_with").HasValue("test/test"),
				check.That(data.ResourceName).Key("advanced_filtering_on_arrays_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopicEventSubscription_advancedFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test1")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.advancedFilter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_eventgrid_system_topic_event_subscription.test2").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopicEventSubscription_advancedFilterMaxItems(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.advancedFilterMaxItems(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopicEventSubscription_systemIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("delivery_identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("delivery_identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("dead_letter_identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("dead_letter_identity.0.type").HasValue("SystemAssigned"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithTopicSystemIdentityEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("delivery_identity.#").HasValue("0"),
				check.That(data.ResourceName).Key("dead_letter_identity.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.systemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("delivery_identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("delivery_identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("dead_letter_identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("dead_letter_identity.0.type").HasValue("SystemAssigned"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopicEventSubscription_userIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("delivery_identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("delivery_identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("dead_letter_identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("dead_letter_identity.0.type").HasValue("UserAssigned"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithTopicUserIdentityEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("delivery_identity.#").HasValue("0"),
				check.That(data.ResourceName).Key("dead_letter_identity.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.userIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("delivery_identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("delivery_identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("dead_letter_identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("dead_letter_identity.0.type").HasValue("UserAssigned"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopicEventSubscription_deliveryPropertiesStatic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deliveryProperties(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),

				check.That(data.ResourceName).Key("delivery_property.0.header_name").HasValue("test-1"),
				check.That(data.ResourceName).Key("delivery_property.0.type").HasValue("Static"),
				check.That(data.ResourceName).Key("delivery_property.0.value").HasValue("1"),
				check.That(data.ResourceName).Key("delivery_property.0.secret").HasValue("false"),

				check.That(data.ResourceName).Key("delivery_property.1.header_name").HasValue("test-2"),
				check.That(data.ResourceName).Key("delivery_property.1.type").HasValue("Static"),
				check.That(data.ResourceName).Key("delivery_property.1.value").HasValue("string"),
				check.That(data.ResourceName).Key("delivery_property.1.secret").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopicEventSubscription_deliveryPropertiesSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deliveryPropertiesSecret(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),

				check.That(data.ResourceName).Key("delivery_property.0.header_name").HasValue("test-secret-1"),
				check.That(data.ResourceName).Key("delivery_property.0.type").HasValue("Static"),
				check.That(data.ResourceName).Key("delivery_property.0.secret").HasValue("true"),
			),
		},
		data.ImportStep("delivery_property.0.value"),
	})
}

func TestAccEventGridSystemTopicEventSubscription_deliveryPropertiesMixed(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deliveryPropertiesWithMultipleTypes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),

				check.That(data.ResourceName).Key("delivery_property.0.header_name").HasValue("test-static-1"),
				check.That(data.ResourceName).Key("delivery_property.0.type").HasValue("Static"),
				check.That(data.ResourceName).Key("delivery_property.0.value").HasValue("1"),
				check.That(data.ResourceName).Key("delivery_property.0.secret").HasValue("false"),

				check.That(data.ResourceName).Key("delivery_property.1.header_name").HasValue("test-dynamic-1"),
				check.That(data.ResourceName).Key("delivery_property.1.type").HasValue("Dynamic"),
				check.That(data.ResourceName).Key("delivery_property.1.source_field").HasValue("data.system"),

				check.That(data.ResourceName).Key("delivery_property.2.header_name").HasValue("test-secret-1"),
				check.That(data.ResourceName).Key("delivery_property.2.type").HasValue("Static"),
				check.That(data.ResourceName).Key("delivery_property.2.secret").HasValue("true"),
				check.That(data.ResourceName).Key("delivery_property.2.value").HasValue("this-value-is-secret!"),
			),
		},
		data.ImportStep("delivery_property.2.value"),
	})
}

func TestAccEventGridSystemTopicEventSubscription_deliveryPropertiesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deliveryPropertiesWithMultipleTypes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("delivery_property.0.header_name").HasValue("test-static-1"),
				check.That(data.ResourceName).Key("delivery_property.0.type").HasValue("Static"),
				check.That(data.ResourceName).Key("delivery_property.0.value").HasValue("1"),
				check.That(data.ResourceName).Key("delivery_property.0.secret").HasValue("false"),
				check.That(data.ResourceName).Key("delivery_property.1.header_name").HasValue("test-dynamic-1"),
				check.That(data.ResourceName).Key("delivery_property.1.type").HasValue("Dynamic"),
				check.That(data.ResourceName).Key("delivery_property.1.source_field").HasValue("data.system"),
				check.That(data.ResourceName).Key("delivery_property.2.header_name").HasValue("test-secret-1"),
				check.That(data.ResourceName).Key("delivery_property.2.type").HasValue("Static"),
				check.That(data.ResourceName).Key("delivery_property.2.secret").HasValue("true"),
				check.That(data.ResourceName).Key("delivery_property.2.value").HasValue("this-value-is-secret!"),
			),
		},
		{
			Config: r.deliveryPropertiesUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("delivery_property.0.header_name").HasValue("test-static-1"),
				check.That(data.ResourceName).Key("delivery_property.0.type").HasValue("Static"),
				check.That(data.ResourceName).Key("delivery_property.0.value").HasValue("2"),
				check.That(data.ResourceName).Key("delivery_property.0.secret").HasValue("false"),
				check.That(data.ResourceName).Key("delivery_property.1.header_name").HasValue("test-dynamic-1"),
				check.That(data.ResourceName).Key("delivery_property.1.type").HasValue("Dynamic"),
				check.That(data.ResourceName).Key("delivery_property.1.source_field").HasValue("data.topic"),
				check.That(data.ResourceName).Key("delivery_property.2.header_name").HasValue("test-secret-1"),
				check.That(data.ResourceName).Key("delivery_property.2.type").HasValue("Static"),
				check.That(data.ResourceName).Key("delivery_property.2.secret").HasValue("true"),
				check.That(data.ResourceName).Key("delivery_property.2.value").HasValue("this-value-is-still-secret!"),
			),
		},
	})
}

func TestAccEventGridSystemTopicEventSubscription_deliveryPropertiesForEventHubs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deliveryPropertiesForEventHubs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),

				check.That(data.ResourceName).Key("delivery_property.0.header_name").HasValue("test-static-1"),
				check.That(data.ResourceName).Key("delivery_property.0.type").HasValue("Static"),
				check.That(data.ResourceName).Key("delivery_property.0.value").HasValue("1"),
				check.That(data.ResourceName).Key("delivery_property.0.secret").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopicEventSubscription_deliveryPropertiesHybridRelay(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic_event_subscription", "test")
	r := EventGridSystemTopicEventSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deliveryPropertiesForHybridRelay(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),

				check.That(data.ResourceName).Key("delivery_property.0.header_name").HasValue("test-static-1"),
				check.That(data.ResourceName).Key("delivery_property.0.type").HasValue("Static"),
				check.That(data.ResourceName).Key("delivery_property.0.value").HasValue("1"),
				check.That(data.ResourceName).Key("delivery_property.0.secret").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func (EventGridSystemTopicEventSubscriptionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := eventsubscriptions.ParseSystemTopicEventSubscriptionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.EventGrid.EventSubscriptions.SystemTopicEventSubscriptionsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (EventGridSystemTopicEventSubscriptionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%[1]d"
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

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.test.id
    queue_name         = azurerm_storage_queue.test.name
  }

  storage_blob_dead_letter_destination {
    storage_account_id          = azurerm_storage_account.test.id
    storage_blob_container_name = azurerm_storage_container.test.name
  }

  retry_policy {
    event_time_to_live    = 11
    max_delivery_attempts = 11
  }

  labels = ["test", "test1", "test2"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) basicWithTopicSystemIdentityEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%[1]d"
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

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.test.id
    queue_name         = azurerm_storage_queue.test.name
  }

  storage_blob_dead_letter_destination {
    storage_account_id          = azurerm_storage_account.test.id
    storage_blob_container_name = azurerm_storage_container.test.name
  }

  retry_policy {
    event_time_to_live    = 11
    max_delivery_attempts = 11
  }

  labels = ["test", "test1", "test2"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) basicWithTopicUserIdentityEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%[1]d"
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

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.test.id
    queue_name         = azurerm_storage_queue.test.name
  }

  storage_blob_dead_letter_destination {
    storage_account_id          = azurerm_storage_account.test.id
    storage_blob_container_name = azurerm_storage_container.test.name
  }

  retry_policy {
    event_time_to_live    = 11
    max_delivery_attempts = 11
  }

  labels = ["test", "test1", "test2"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) requiresImport(data acceptance.TestData) string {
	template := EventGridSystemTopicEventSubscriptionResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_system_topic_event_subscription" "import" {
  name                = azurerm_eventgrid_system_topic_event_subscription.test.name
  system_topic        = azurerm_eventgrid_system_topic_event_subscription.test.system_topic
  resource_group_name = azurerm_eventgrid_system_topic_event_subscription.test.resource_group_name
}
`, template)
}

func (EventGridSystemTopicEventSubscriptionResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%[1]d"
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

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  storage_queue_endpoint {
    storage_account_id                    = azurerm_storage_account.test.id
    queue_name                            = azurerm_storage_queue.test.name
    queue_message_time_to_live_in_seconds = 3600
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) eventHubID(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  event_delivery_schema = "CloudEventSchemaV1_0"

  eventhub_endpoint_id = azurerm_eventhub.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) serviceBusQueueID(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "example" {
  name                         = "acctestservicebusnamespace-%[1]d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku                          = "Premium"
  premium_messaging_partitions = 2
  capacity                     = 2
}

resource "azurerm_servicebus_queue" "test" {
  name                 = "acctestservicebusqueue-%[1]d"
  namespace_id         = azurerm_servicebus_namespace.example.id
  partitioning_enabled = true
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  event_delivery_schema         = "CloudEventSchemaV1_0"
  service_bus_queue_endpoint_id = azurerm_servicebus_queue.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) serviceBusTopicID(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "acctestservicebusnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                 = "acctestservicebustopic-%[1]d"
  namespace_id         = azurerm_servicebus_namespace.example.id
  partitioning_enabled = true
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  event_delivery_schema         = "CloudEventSchemaV1_0"
  service_bus_topic_endpoint_id = azurerm_servicebus_topic.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) filter(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%[1]d"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.test.id
    queue_name         = azurerm_storage_queue.test.name
  }

  advanced_filtering_on_arrays_enabled = true

  included_event_types = ["Microsoft.Storage.BlobCreated", "Microsoft.Storage.BlobDeleted"]

  subject_filter {
    subject_begins_with = "test/test"
    subject_ends_with   = ".jpg"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) advancedFilter(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%[1]d"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test1" {
  name                = "acctesteg-%[1]d-1"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.test.id
    queue_name         = azurerm_storage_queue.test.name
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
    number_in_range {
      key    = "data.contentLength"
      values = [[0, 1], [2, 3]]
    }
    number_not_in_range {
      key    = "data.contentLength"
      values = [[5, 13], [21, 34]]
    }
    string_begins_with {
      key    = "subject"
      values = ["foo"]
    }
  }
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test2" {
  name                = "acctesteg-%[1]d-2"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.test.id
    queue_name         = azurerm_storage_queue.test.name
  }

  advanced_filter {
    string_ends_with {
      key    = "subject"
      values = ["bar"]
    }
    string_not_begins_with {
      key    = "subject"
      values = ["lorem"]
    }
    string_not_ends_with {
      key    = "subject"
      values = ["ipsum"]
    }
    string_not_contains {
      key    = "data.contentType"
      values = ["text"]
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
    is_not_null {
      key = "subject"
    }
    is_null_or_undefined {
      key = "subject"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) advancedFilterMaxItems(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%[1]d"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.test.id
    queue_name         = azurerm_storage_queue.test.name
  }

  advanced_filter {
    bool_equals {
      key   = "subject"
      value = true
    }
    number_greater_than {
      key   = "data.metadataVersion"
      value = 2
    }
    number_greater_than_or_equals {
      key   = "data.contentLength"
      value = 3
    }
    number_less_than {
      key   = "data.contentLength"
      value = 4
    }
    number_less_than_or_equals {
      key   = "data.metadataVersion"
      value = 5
    }
    number_in {
      key    = "data.contentLength"
      values = [6, 7, 8]
    }
    number_not_in {
      key    = "data.contentLength"
      values = [9, 10, 11]
    }
    string_begins_with {
      key    = "subject"
      values = ["12", "13", "14"]
    }
    string_ends_with {
      key    = "subject"
      values = ["15", "16", "17"]
    }
    string_contains {
      key    = "data.contentType"
      values = ["18", "19", "20"]
    }
    string_in {
      key    = "data.blobType"
      values = ["21", "22", "23"]
    }
    string_not_in {
      key    = "data.blobType"
      values = ["24", "25"]
    }
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) systemIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%[1]d"
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

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_role_assignment" "contributor" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_eventgrid_system_topic.test.identity.0.principal_id
}

resource "azurerm_role_assignment" "sender" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Queue Data Message Sender"
  principal_id         = azurerm_eventgrid_system_topic.test.identity.0.principal_id
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  delivery_identity {
    type = "SystemAssigned"
  }

  dead_letter_identity {
    type = "SystemAssigned"
  }

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.test.id
    queue_name         = azurerm_storage_queue.test.name
  }

  storage_blob_dead_letter_destination {
    storage_account_id          = azurerm_storage_account.test.id
    storage_blob_container_name = azurerm_storage_container.test.name
  }

  depends_on = [
    azurerm_role_assignment.contributor,
    azurerm_role_assignment.sender
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) userIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%[1]d"
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

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}

resource "azurerm_role_assignment" "contributor" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "sender" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Queue Data Message Sender"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  delivery_identity {
    type                   = "UserAssigned"
    user_assigned_identity = azurerm_user_assigned_identity.test.id
  }

  dead_letter_identity {
    type                   = "UserAssigned"
    user_assigned_identity = azurerm_user_assigned_identity.test.id
  }

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.test.id
    queue_name         = azurerm_storage_queue.test.name
  }

  storage_blob_dead_letter_destination {
    storage_account_id          = azurerm_storage_account.test.id
    storage_blob_container_name = azurerm_storage_container.test.name
  }

  depends_on = [
    azurerm_role_assignment.contributor,
    azurerm_role_assignment.sender
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) deliveryProperties(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "acctestservicebusnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                 = "acctestservicebustopic-%[1]d"
  namespace_id         = azurerm_servicebus_namespace.example.id
  partitioning_enabled = true
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  service_bus_topic_endpoint_id = azurerm_servicebus_topic.test.id

  advanced_filtering_on_arrays_enabled = true

  subject_filter {
    subject_begins_with = "test/test"
  }

  delivery_property {
    header_name = "test-1"
    type        = "Static"
    value       = "1"
    secret      = false
  }

  delivery_property {
    header_name = "test-2"
    type        = "Static"
    value       = "string"
    secret      = false
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) deliveryPropertiesSecret(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "acctestservicebusnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                 = "acctestservicebustopic-%[1]d"
  namespace_id         = azurerm_servicebus_namespace.example.id
  partitioning_enabled = true
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  service_bus_topic_endpoint_id = azurerm_servicebus_topic.test.id

  advanced_filtering_on_arrays_enabled = true

  subject_filter {
    subject_begins_with = "test/test"
  }

  delivery_property {
    header_name = "test-secret-1"
    type        = "Static"
    value       = "1"
    secret      = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) deliveryPropertiesWithMultipleTypes(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "acctestservicebusnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                 = "acctestservicebustopic-%[1]d"
  namespace_id         = azurerm_servicebus_namespace.example.id
  partitioning_enabled = true
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  service_bus_topic_endpoint_id = azurerm_servicebus_topic.test.id

  advanced_filtering_on_arrays_enabled = true

  subject_filter {
    subject_begins_with = "test/test"
  }

  delivery_property {
    header_name = "test-static-1"
    type        = "Static"
    value       = "1"
    secret      = false
  }

  delivery_property {
    header_name  = "test-dynamic-1"
    type         = "Dynamic"
    source_field = "data.system"
  }

  delivery_property {
    header_name = "test-secret-1"
    type        = "Static"
    value       = "this-value-is-secret!"
    secret      = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) deliveryPropertiesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "acctestservicebusnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                 = "acctestservicebustopic-%[1]d"
  namespace_id         = azurerm_servicebus_namespace.example.id
  partitioning_enabled = true
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  service_bus_topic_endpoint_id = azurerm_servicebus_topic.test.id

  advanced_filtering_on_arrays_enabled = true

  subject_filter {
    subject_begins_with = "test/test"
  }

  delivery_property {
    header_name = "test-static-1"
    type        = "Static"
    value       = "2"
    secret      = false
  }

  delivery_property {
    header_name  = "test-dynamic-1"
    type         = "Dynamic"
    source_field = "data.topic"
  }

  delivery_property {
    header_name = "test-secret-1"
    type        = "Static"
    value       = "this-value-is-still-secret!"
    secret      = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (EventGridSystemTopicEventSubscriptionResource) deliveryPropertiesForEventHubs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  eventhub_endpoint_id = azurerm_eventhub.test.id

  delivery_property {
    header_name = "test-static-1"
    type        = "Static"
    value       = "1"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (EventGridSystemTopicEventSubscriptionResource) deliveryPropertiesForHybridRelay(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_relay_namespace" "test" {
  name                = "acctest-%[1]d-rly-eventsub-repo"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Standard"
}

resource "azurerm_relay_hybrid_connection" "test" {
  name                          = "acctest-%[1]d-rhc-eventsub-repo"
  resource_group_name           = azurerm_resource_group.test.name
  relay_namespace_name          = azurerm_relay_namespace.test.name
  requires_client_authorization = false
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_resource_group.test.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "test" {
  name                = "acctesteg-%[1]d"
  system_topic        = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  hybrid_connection_endpoint_id = azurerm_relay_hybrid_connection.test.id

  delivery_property {
    header_name = "test-static-1"
    type        = "Static"
    value       = "1"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
