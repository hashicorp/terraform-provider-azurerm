package servicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ServiceBusQueueResource struct {
}

func TestAccServiceBusQueue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_express").HasValue("false"),
				check.That(data.ResourceName).Key("enable_partitioning").HasValue("false"),
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
				check.That(data.ResourceName).Key("enable_express").HasValue("false"),
				check.That(data.ResourceName).Key("enable_partitioning").HasValue("false"),
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
				check.That(data.ResourceName).Key("enable_express").HasValue("false"),
				check.That(data.ResourceName).Key("enable_batched_operations").HasValue("true"),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("enable_express").HasValue("true"),
				check.That(data.ResourceName).Key("max_size_in_megabytes").HasValue("2048"),
				check.That(data.ResourceName).Key("enable_batched_operations").HasValue("false"),
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
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_partitioning").HasValue("false"),
			),
		},
		{
			Config: r.enablePartitioningStandard(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("enable_partitioning").HasValue("true"),
				// Ensure size is read back in its original value and not the x16 value returned by Azure
				check.That(data.ResourceName).Key("max_size_in_megabytes").HasValue("5120"),
			),
		},
	})
}

func TestAccServiceBusQueue_defaultEnablePartitioningPremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.Premium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_partitioning").HasValue("false"),
				check.That(data.ResourceName).Key("enable_express").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueue_enableDuplicateDetection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServiceBusQueueResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("requires_duplicate_detection").HasValue("false"),
			),
		},
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
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("requires_session").HasValue("false"),
			),
		},
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
				check.That(data.ResourceName).Key("lock_duration").HasValue("PT40S"),
			),
		},
		{
			Config: r.lockDurationUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("lock_duration").HasValue("PT2M"),
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
				check.That(data.ResourceName).Key("max_delivery_count").HasValue("10"),
			),
		},
		{
			Config: r.maxDeliveryCount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("max_delivery_count").HasValue("20"),
			),
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
	id, err := parse.QueueID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.QueuesClient.Get(ctx, id.ResourceGroup, id.NamespaceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Service Bus NameSpace Queue (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
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
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctestservicebusqueue-%d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ServiceBusQueueResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_queue" "import" {
  name                = azurerm_servicebus_queue.test.name
  resource_group_name = azurerm_servicebus_queue.test.resource_group_name
  namespace_name      = azurerm_servicebus_queue.test.namespace_name
}
`, r.basic(data))
}

func (ServiceBusQueueResource) Premium(data acceptance.TestData) string {
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
  sku                 = "premium"
  capacity            = 1
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctestservicebusqueue-%d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
  enable_partitioning = false
  enable_express      = false
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
  name                      = "acctestservicebusqueue-%d"
  resource_group_name       = azurerm_resource_group.test.name
  namespace_name            = azurerm_servicebus_namespace.test.name
  enable_express            = true
  max_size_in_megabytes     = 2048
  enable_batched_operations = false
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
  resource_group_name   = azurerm_resource_group.test.name
  namespace_name        = azurerm_servicebus_namespace.test.name
  enable_partitioning   = true
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
  resource_group_name          = azurerm_resource_group.test.name
  namespace_name               = azurerm_servicebus_namespace.test.name
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

  sku = "standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctestservicebusqueue-%d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
  requires_session    = true
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

  sku = "standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                                 = "acctestservicebusqueue-%d"
  resource_group_name                  = azurerm_resource_group.test.name
  namespace_name                       = azurerm_servicebus_namespace.test.name
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
  name                = "acctestservicebusqueue-%d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
  lock_duration       = "PT40S"
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
  name                = "acctestservicebusqueue-%d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
  lock_duration       = "PT2M"
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
  resource_group_name                     = azurerm_resource_group.test.name
  namespace_name                          = azurerm_servicebus_namespace.test.name
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
  name                = "acctestservicebusqueue-%d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
  max_delivery_count  = 20
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
  name                = "acctestservicebusqueue-forward_to-%d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctestservicebusqueue-%d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
  forward_to          = azurerm_servicebus_queue.forward_to.name
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
  name                = "acctestservicebusqueue-forward_dl_messages_to-%d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
}

resource "azurerm_servicebus_queue" "test" {
  name                              = "acctestservicebusqueue-%d"
  resource_group_name               = azurerm_resource_group.test.name
  namespace_name                    = azurerm_servicebus_namespace.test.name
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
  name                = "acctestservicebusqueue-%d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
  status              = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, status)
}
