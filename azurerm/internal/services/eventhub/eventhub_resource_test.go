package eventhub_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

type EventHubResource struct {
}

func TestAccEventHubPartitionCount_validation(t *testing.T) {
	cases := []struct {
		Value    int
		ErrCount int
	}{
		{
			Value:    0,
			ErrCount: 1,
		},
		{
			Value:    1,
			ErrCount: 0,
		},
		{
			Value:    2,
			ErrCount: 0,
		},
		{
			Value:    3,
			ErrCount: 0,
		},
		{
			Value:    21,
			ErrCount: 0,
		},
		{
			Value:    1024,
			ErrCount: 0,
		},
		{
			Value:    1025,
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validate.ValidateEventHubPartitionCount(tc.Value, "azurerm_eventhub")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM EventHub Partition Count to trigger a validation error")
		}
	}
}

func TestAccEventHubMessageRetentionCount_validation(t *testing.T) {
	cases := []struct {
		Value    int
		ErrCount int
	}{
		{
			Value:    0,
			ErrCount: 1,
		}, {
			Value:    1,
			ErrCount: 0,
		}, {
			Value:    2,
			ErrCount: 0,
		}, {
			Value:    3,
			ErrCount: 0,
		}, {
			Value:    4,
			ErrCount: 0,
		}, {
			Value:    5,
			ErrCount: 0,
		}, {
			Value:    6,
			ErrCount: 0,
		}, {
			Value:    90,
			ErrCount: 0,
		}, {
			Value:    91,
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validate.ValidateEventHubMessageRetentionCount(tc.Value, "azurerm_eventhub")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM EventHub Message Retention Count to trigger a validation error")
		}
	}
}

func TestAccEventHubArchiveNameFormat_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "",
			ErrCount: 9,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 0,
		},
		{
			Value:    "Prod_{Eventub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validate.ValidateEventHubArchiveNameFormat(tc.Value, "azurerm_eventhub")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q to trigger a validation error", tc.Value)
		}
	}
}

func TestAccEventHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub", "test")
	r := EventHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHub_basicOnePartition(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub", "test")
	r := EventHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("partition_count").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub", "test")
	r := EventHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_eventhub"),
		},
	})
}

func TestAccEventHub_partitionCountUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub", "test")
	r := EventHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("partition_count").HasValue("2"),
			),
		},
		{
			Config: r.partitionCountUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("partition_count").HasValue("10"),
			),
		},
	})
}

func TestAccEventHub_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub", "test")
	r := EventHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standard(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHub_captureDescription(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub", "test")
	r := EventHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.captureDescription(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("capture_description.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("capture_description.0.skip_empty_archives").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHub_captureDescriptionDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub", "test")
	r := EventHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.captureDescription(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("capture_description.0.enabled").HasValue("true"),
			),
		},
		{
			Config: r.captureDescription(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("capture_description.0.enabled").HasValue("false"),
			),
		},
	})
}

func TestAccEventHub_messageRetentionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub", "test")
	r := EventHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standard(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("message_retention").HasValue("7"),
			),
		},
		{
			Config: r.messageRetentionUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("message_retention").HasValue("5"),
			),
		},
	})
}

func (EventHubResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.EventHubID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Eventhub.EventHubsClient.Get(ctx, id.ResourceGroup, id.NamespaceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (EventHubResource) basic(data acceptance.TestData, partitionCount int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
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
  partition_count     = %d
  message_retention   = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, partitionCount)
}

func (EventHubResource) requiresImport(data acceptance.TestData) string {
	template := EventHubResource{}.basic(data, 2)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub" "import" {
  name                = azurerm_eventhub.test.name
  namespace_name      = azurerm_eventhub.test.namespace_name
  resource_group_name = azurerm_eventhub.test.resource_group_name
  partition_count     = azurerm_eventhub.test.partition_count
  message_retention   = azurerm_eventhub.test.message_retention
}
`, template)
}

func (EventHubResource) partitionCountUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
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
  partition_count     = 10
  message_retention   = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (EventHubResource) standard(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 7
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (EventHubResource) captureDescription(data acceptance.TestData, enabled bool) string {
	enabledString := strconv.FormatBool(enabled)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctest"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 7

  capture_description {
    enabled             = %s
    encoding            = "Avro"
    interval_in_seconds = 60
    size_limit_in_bytes = 10485760
    skip_empty_archives = true

    destination {
      name                = "EventHubArchive.AzureBlockBlob"
      archive_name_format = "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}"
      blob_container_name = azurerm_storage_container.test.name
      storage_account_id  = azurerm_storage_account.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, enabledString)
}

func (EventHubResource) messageRetentionUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 5
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
