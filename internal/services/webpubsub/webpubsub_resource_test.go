package webpubsub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/webpubsub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WebpubsubResource struct{}

func TestAccWebpubsub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub", "test")
	r := WebpubsubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebpubsub_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub", "test")
	r := WebpubsubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebpubsub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub", "test")
	r := WebpubsubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccWebpubsub_free(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub", "test")
	r := WebpubsubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.FreeWithCapacity(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r WebpubsubResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WebPubsubID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Webpubsub.WebPubsubClient.Get(ctx, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r WebpubsubResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_web_pubsub" "test" {
  name                = "acctestWebPubsub-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_S1"
  capacity            = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r WebpubsubResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-wps-%d"
  location = "%s"
}

resource "azurerm_web_pubsub" "test" {
  name                = "acctestWebPubsub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku      = "Standard_S1"
  capacity = 1

  public_network_access_enabled = false

  live_trace_configuration {
    enabled = false
    categories {
      name    = "MessagingLogs"
      enabled = true
    }
  }

  local_auth_enabled  = true
  aad_auth_enabled = true

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r WebpubsubResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-wps-%d"
  location = "%s"
}

resource "azurerm_web_pubsub" "import" {
  name                = azurerm_web_pubsub.test.name
  location            = azurerm_web_pubsub.test.location
  resource_group_name = azurerm_web_pubsub.test.resource_group_name

  sku      = "Standard_S1"
  capacity = 1
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r WebpubsubResource) FreeWithCapacity(data acceptance.TestData, capacity int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_web_pubsub" "test" {
  name                = "acctestWebPubsub-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku      = "Free_F1"
  capacity = %d
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, capacity)
}
