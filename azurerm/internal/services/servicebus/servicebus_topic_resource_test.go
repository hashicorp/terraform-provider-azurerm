package servicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

type ServiceBusTopicResource struct {
}

func TestAccServiceBusTopic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")
	r := ServiceBusTopicResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusTopic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")
	r := ServiceBusTopicResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_service_fabric_cluster"),
		},
	})
}

func TestAccServiceBusTopic_basicDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")
	r := ServiceBusTopicResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicDisabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusTopic_basicDisableEnable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")
	r := ServiceBusTopicResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicDisabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccServiceBusTopic_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")
	r := ServiceBusTopicResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("enable_batched_operations").HasValue("true"),
				check.That(data.ResourceName).Key("enable_express").HasValue("true"),
			),
		},
	})
}

func TestAccServiceBusTopic_enablePartitioningStandard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")
	r := ServiceBusTopicResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.enablePartitioningStandard(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("enable_partitioning").HasValue("true"),
				// Ensure size is read back in its original value and not the x16 value returned by Azure
				check.That(data.ResourceName).Key("max_size_in_megabytes").HasValue("5120"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusTopic_enablePartitioningPremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")
	r := ServiceBusTopicResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicPremium(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.enablePartitioningPremium(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("enable_partitioning").HasValue("false"),
				check.That(data.ResourceName).Key("max_size_in_megabytes").HasValue("81920"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusTopic_enableDuplicateDetection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")
	r := ServiceBusTopicResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.enableDuplicateDetection(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("requires_duplicate_detection").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusTopic_isoTimeSpanAttributes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")
	r := ServiceBusTopicResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.isoTimeSpanAttributes(data),
			Check: resource.ComposeTestCheckFunc(
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

func (t ServiceBusTopicResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.TopicID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.TopicsClient.Get(ctx, id.ResourceGroup, id.NamespaceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Service Bus Topic (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ServiceBusTopicResource) basic(data acceptance.TestData) string {
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
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ServiceBusTopicResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_topic" "import" {
  name                = azurerm_servicebus_topic.test.name
  namespace_name      = azurerm_servicebus_topic.test.namespace_name
  resource_group_name = azurerm_servicebus_topic.test.resource_group_name
}
`, r.basic(data))
}

func (ServiceBusTopicResource) basicDisabled(data acceptance.TestData) string {
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
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  status              = "disabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusTopicResource) update(data acceptance.TestData) string {
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
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                      = "acctestservicebustopic-%d"
  namespace_name            = azurerm_servicebus_namespace.test.name
  resource_group_name       = azurerm_resource_group.test.name
  enable_batched_operations = true
  enable_express            = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusTopicResource) basicPremium(data acceptance.TestData) string {
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
  sku                 = "Premium"
  capacity            = 1
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  enable_partitioning = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusTopicResource) enablePartitioningStandard(data acceptance.TestData) string {
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
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                  = "acctestservicebustopic-%d"
  namespace_name        = azurerm_servicebus_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  enable_partitioning   = true
  max_size_in_megabytes = 5120
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusTopicResource) enablePartitioningPremium(data acceptance.TestData) string {
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
  sku                 = "premium"
  capacity            = 1
}

resource "azurerm_servicebus_topic" "test" {
  name                  = "acctestservicebustopic-%d"
  namespace_name        = azurerm_servicebus_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  enable_partitioning   = false
  max_size_in_megabytes = 81920
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusTopicResource) enableDuplicateDetection(data acceptance.TestData) string {
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
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                         = "acctestservicebustopic-%d"
  namespace_name               = azurerm_servicebus_namespace.test.name
  resource_group_name          = azurerm_resource_group.test.name
  requires_duplicate_detection = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusTopicResource) isoTimeSpanAttributes(data acceptance.TestData) string {
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
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                                    = "acctestservicebustopic-%d"
  namespace_name                          = azurerm_servicebus_namespace.test.name
  resource_group_name                     = azurerm_resource_group.test.name
  auto_delete_on_idle                     = "PT10M"
  default_message_ttl                     = "PT30M"
  requires_duplicate_detection            = true
  duplicate_detection_history_time_window = "PT15M"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
