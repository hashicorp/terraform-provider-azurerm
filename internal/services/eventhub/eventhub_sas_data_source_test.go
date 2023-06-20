package eventhub_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type EventHubSharedAccessSignatureDataSource struct{}

func TestAccEventHubSharedAccessSignatureDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_sas", "test")
	r := EventHubSharedAccessSignatureDataSource{}
	utcNow := time.Now().UTC()
	endDate := utcNow.Add(time.Hour * 24).Format(time.RFC3339)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data, endDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sas").Exists(),
			),
		},
	})
}

func (EventHubSharedAccessSignatureDataSource) basic(data acceptance.TestData, endDate string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ehn-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-ehn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-eh-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  name                = "acctest-ehar-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = true
  send   = true
  manage = true
}

data "azurerm_eventhub_authorization_rule" "test" {
  name                = azurerm_eventhub_authorization_rule.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_eventhub_sas" "test" {
  connection_string = data.azurerm_eventhub_authorization_rule.test.primary_connection_string
  expiry            = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, endDate)
}
