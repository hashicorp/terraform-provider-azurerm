package media_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MediaPrivateEndpointConnectionResource struct {
}

func TestAccMediaPrivateEndpointConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_private_endpoint_connection", "test")
	r := MediaPrivateEndpointConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("provisioning_state").HasValue("Succeeded"),
				check.That(data.ResourceName).Key("private_link_connection_state.0.status").HasValue("Approved"),
				check.That(data.ResourceName).Key("private_link_connection_state.0.description").Exists(),
				check.That(data.ResourceName).Key("private_link_connection_state.0.actions_required").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r MediaPrivateEndpointConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.PrivateEndpointConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.PrivateEndpointConnectionsClient.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Private Endpoint Connection %s (Media Account %s) (resource group: %s): %v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.PrivateLinkServiceConnectionState != nil), nil
}

func (r MediaPrivateEndpointConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-media-%d"
  location = "%s"
}

resource "azurerm_storage_account" "first" {
  name                     = "acctestsa1%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }

  tags = {
    environment = "staging"
  }
}

resource "azurerm_media_private_endpoint_connection" "test" {
  name                        = "endpoint1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
