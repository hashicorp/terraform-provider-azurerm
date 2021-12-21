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

type WebPubsubHubDataSource struct{}

func TestAccDataSourceWebPubsubHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_web_pubsub_hub", "test")
	r := WebPubsubHubDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r WebPubsubHubDataSource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WebPubsubHubID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Webpubsub.WebPubsubHubsClient.Get(ctx, id.HubName, id.ResourceGroup, id.WebPubsubName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Web Pubsub Hub (%q): %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r WebPubsubHubDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-wps-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_web_pubsub" "test" {
  name                = "acctest-wps-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku {
    name = "Standard_S1"
  }
}

resource "azurerm_web_pubsub_hub" "test" {
  name                = "acctestwpshub%[1]d"
  web_pubsub_name     = azurerm_web_pubsub.test.name
  resource_group_name = azurerm_resource_group.test.name
  event_handler {
    url_template       = "https://test.com/api/{hub}/{event}"
    user_event_pattern = "event1, event2"
    system_events      = ["connect", "connected"]
    auth {
      type                      = "ManagedIdentity"
      managed_identity_resource = azurerm_user_assigned_identity.test.id
    }
  }
  depends_on = [azurerm_web_pubsub.test]
}

data "azurerm_web_pubsub_hub" "test" {
  name                = azurerm_web_pubsub_hub.test.name
  web_pubsub_name     = azurerm_web_pubsub.test.name
  resource_group_name = azurerm_resource_group.test.name
}

`, data.RandomInteger, data.Locations.Primary)
}
