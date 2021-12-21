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

type WebPubsubHubResource struct{}

func TestAccWebPubsubHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_hub", "test")
	r := WebPubsubHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubsubHub_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_hub", "test")
	r := WebPubsubHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func (r WebPubsubHubResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WebPubsubHubID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Webpubsub.WebPubsubHubsClient.Get(ctx, id.HubName, id.ResourceGroup, id.WebPubsubName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Web Pubsub Hub (%q): %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r WebPubsubHubResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_web_pubsub_hub" "test" {
  name                = "acctestwpsh%d"
  web_pubsub_name     = azurerm_web_pubsub.test.name
  resource_group_name = azurerm_resource_group.test.name

  event_handler {
    url_template       = "https://test.com/api/{hub}/{event}"
    user_event_pattern = "event1, event2"
    system_events      = ["connect", "connected"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r WebPubsubHubResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}


resource "azurerm_web_pubsub_hub" "test" {
  name                = "acctestwpsh%d"
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
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r WebPubsubHubResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-wps-%d"
  location = "%s"
}

resource "azurerm_web_pubsub" "test" {
  name                = "acctest-webpubsub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku {
    name = "Standard_S1"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
