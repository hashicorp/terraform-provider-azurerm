// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServiceBusQueueResource struct{}

func TestAccServiceBusQueue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("express_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("partitioning_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("express_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("partitioning_enabled").HasValue("false"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccServiceBusQueue_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("express_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("batched_operations_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("express_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("max_size_in_megabytes").HasValue("2048"),
				check.That(data.ResourceName).Key("batched_operations_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_enablePartitioningStandard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enablePartitioningStandard(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("partitioning_enabled").HasValue("true"),
				// Ensure size is read back in its original value and not the x16 value returned by Azure
				check.That(data.ResourceName).Key("max_size_in_megabytes").HasValue("5120"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_maxMessageSizePremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.PremiumNamespaceNonPartitioned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_partitionedPremiumNamespace(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.PremiumNamespacePartitioned(data, false),
			ExpectError: regexp.MustCompile("non-partitioned entities are not allowed in partitioned namespace"),
		},
		{
			Config: r.PremiumNamespacePartitioned(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("partitioning_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("express_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_nonPartitionedPremiumNamespaceError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.nonPartitionedPremiumNamespaceError(data),
			ExpectError: regexp.MustCompile("the parent premium namespace is not partitioned and the partitioning for premium namespace is only available at the namepsace creation"),
		},
	})
}

func TestAccServiceBusQueue_enableDuplicateDetection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enableDuplicateDetection(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("requires_duplicate_detection").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_enableRequiresSession(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enableRequiresSession(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("requires_session").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_enableDeadLetteringOnMessageExpiration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("dead_lettering_on_message_expiration").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableDeadLetteringOnMessageExpiration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("dead_lettering_on_message_expiration").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_lockDuration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.lockDuration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.lockDurationUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_isoTimeSpanAttributes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.isoTimeSpanAttributes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_delete_on_idle").HasValue("PT10M"),
				check.That(data.ResourceName).Key("default_message_ttl").HasValue("PT30M"),
				check.That(data.ResourceName).Key("requires_duplicate_detection").HasValue("true"),
				check.That(data.ResourceName).Key("duplicate_detection_history_time_window").HasValue("PT15M"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_maxDeliveryCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.maxDeliveryCount(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_forwardTo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("forward_to").HasValue(""),
			),
		},
		data.ImportStep(),
		{
			Config: r.forwardTo(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("forward_to").HasValue(fmt.Sprintf("acctestservicebusqueue-forward_to-%d", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_forwardDeadLetteredMessagesTo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("forward_dead_lettered_messages_to").HasValue(""),
			),
		},
		data.ImportStep(),
		{
			Config: r.forwardDeadLetteredMessagesTo(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("forward_dead_lettered_messages_to").HasValue(fmt.Sprintf("acctestservicebusqueue-forward_dl_messages_to-%d", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_status(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("status").HasValue("Active"),
			),
		},
		{
			Config: r.status(data, "Creating"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("status").HasValue("Creating"),
			),
		},
		{
			Config: r.status(data, "Deleting"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("status").HasValue("Deleting"),
			),
		},
		{
			Config: r.status(data, "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("status").HasValue("Disabled"),
			),
		},
		{
			Config: r.status(data, "ReceiveDisabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("status").HasValue("ReceiveDisabled"),
			),
		},
		{
			Config: r.status(data, "Renaming"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("status").HasValue("Renaming"),
			),
		},
		{
			Config: r.status(data, "SendDisabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("status").HasValue("SendDisabled"),
			),
		},
		{
			Config: r.status(data, "Unknown"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("status").HasValue("Unknown"),
			),
		},
		{
			Config: r.status(data, "Active"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("status").HasValue("Active"),
			),
		},
		data.ImportStep(),
	})
}

func (t ServiceBusQueueResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := queues.ParseQueueID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.QueuesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ServiceBusQueueResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"
}

resource "azurerm_servicebus_queue" "test" {
  name         = "acctestservicebusqueue-%d"
  namespace_id = azurerm_servicebus_namespace.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ServiceBusQueueResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_queue" "import" {
  name         = azurerm_servicebus_queue.test.name
  namespace_id = azurerm_servicebus_queue.test.namespace_id
}
`, r.basic(data))
}

func (ServiceBusQueueResource) PremiumNamespaceNonPartitioned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                         = "acctestservicebusnamespace-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  sku                          = "Premium"
  premium_messaging_partitions = 1
  capacity                     = 1
}

resource "azurerm_servicebus_queue" "test" {
  name                 = "acctestservicebusqueue-%d"
  namespace_id         = azurerm_servicebus_namespace.test.id
  partitioning_enabled = false
  express_enabled      = false

  max_message_size_in_kilobytes = 102400
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) PremiumNamespacePartitioned(data acceptance.TestData, enabled bool) string {
	// Limited regional availability for premium namespace partitions
	data.Locations.Primary = "westus"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                         = "acctestservicebusnamespace-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  sku                          = "Premium"
  premium_messaging_partitions = 2
  capacity                     = 2
}

resource "azurerm_servicebus_queue" "test" {
  name                 = "acctestservicebusqueue-%d"
  namespace_id         = azurerm_servicebus_namespace.test.id
  partitioning_enabled = %t
  express_enabled      = false

  max_message_size_in_kilobytes = 102400
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, enabled)
}

func (ServiceBusQueueResource) nonPartitionedPremiumNamespaceError(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                         = "acctestservicebusnamespace-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  sku                          = "Premium"
  premium_messaging_partitions = 1
  capacity                     = 1
}

resource "azurerm_servicebus_queue" "test" {
  name         = "acctestservicebusqueue-%d"
  namespace_id = azurerm_servicebus_namespace.test.id

  partitioning_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                       = "acctestservicebusqueue-%d"
  namespace_id               = azurerm_servicebus_namespace.test.id
  express_enabled            = true
  max_size_in_megabytes      = 2048
  batched_operations_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) enablePartitioningStandard(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                  = "acctestservicebusqueue-%d"
  namespace_id          = azurerm_servicebus_namespace.test.id
  partitioning_enabled  = true
  max_size_in_megabytes = 5120
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) enableDuplicateDetection(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                         = "acctestservicebusqueue-%d"
  namespace_id                 = azurerm_servicebus_namespace.test.id
  requires_duplicate_detection = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) enableRequiresSession(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name             = "acctestservicebusqueue-%d"
  namespace_id     = azurerm_servicebus_namespace.test.id
  requires_session = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) enableDeadLetteringOnMessageExpiration(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                                 = "acctestservicebusqueue-%d"
  namespace_id                         = azurerm_servicebus_namespace.test.id
  dead_lettering_on_message_expiration = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) lockDuration(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name          = "acctestservicebusqueue-%d"
  namespace_id  = azurerm_servicebus_namespace.test.id
  lock_duration = "PT40S"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) lockDurationUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name          = "acctestservicebusqueue-%d"
  namespace_id  = azurerm_servicebus_namespace.test.id
  lock_duration = "PT2M"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) isoTimeSpanAttributes(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                                    = "acctestservicebusqueue-%d"
  namespace_id                            = azurerm_servicebus_namespace.test.id
  auto_delete_on_idle                     = "PT10M"
  default_message_ttl                     = "PT30M"
  requires_duplicate_detection            = true
  duplicate_detection_history_time_window = "PT15M"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) maxDeliveryCount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name               = "acctestservicebusqueue-%d"
  namespace_id       = azurerm_servicebus_namespace.test.id
  max_delivery_count = 20
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) forwardTo(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "forward_to" {
  name         = "acctestservicebusqueue-forward_to-%d"
  namespace_id = azurerm_servicebus_namespace.test.id
}

resource "azurerm_servicebus_queue" "test" {
  name         = "acctestservicebusqueue-%d"
  namespace_id = azurerm_servicebus_namespace.test.id
  forward_to   = azurerm_servicebus_queue.forward_to.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) forwardDeadLetteredMessagesTo(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "forward_dl_messages_to" {
  name         = "acctestservicebusqueue-forward_dl_messages_to-%d"
  namespace_id = azurerm_servicebus_namespace.test.id
}

resource "azurerm_servicebus_queue" "test" {
  name                              = "acctestservicebusqueue-%d"
  namespace_id                      = azurerm_servicebus_namespace.test.id
  forward_dead_lettered_messages_to = azurerm_servicebus_queue.forward_dl_messages_to.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusQueueResource) status(data acceptance.TestData, status string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name         = "acctestservicebusqueue-%d"
  namespace_id = azurerm_servicebus_namespace.test.id
  status       = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, status)
}
