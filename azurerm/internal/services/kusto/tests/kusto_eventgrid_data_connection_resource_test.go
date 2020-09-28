package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMKustoEventGridDataConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventgrid_data_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoEventHubDataConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoEventGridDataConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoEventHubDataConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMKustoEventGridDataConnection_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctesteventhubcg-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_eventgrid_event_subscription" "test" {
  name                 = "acctesteg-%d"
  scope                = azurerm_storage_account.test.id
  eventhub_endpoint_id = azurerm_eventhub.test.id
  included_event_types = ["Microsoft.Storage.BlobCreated", "Microsoft.Storage.BlobRenamed"]
  subject_filter {
    subject_begins_with = "/blobServices/default/containers/"
  }

  retry_policy {
    event_time_to_live    = 10
    max_delivery_attempts = 10
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_kusto_eventgrid_data_connection" "test" {
  name                = "acctestkrgdc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
  database_name       = azurerm_kusto_database.test.name

  storage_account_id = azurerm_storage_account.test.id
  eventhub_id        = azurerm_eventhub.test.id
  consumer_group     = azurerm_eventhub_consumer_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
