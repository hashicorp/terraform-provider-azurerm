package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PrivateEndpointConnectionResource struct{}

func (a PrivateEndpointConnectionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PrivateEndpointConnectionID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automation.PrivateEndpointClient.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", id, err)
	}
	return utils.Bool(resp.PrivateEndpointConnectionProperties != nil), nil
}

func (a PrivateEndpointConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-pe-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_automation_account.test.name
    is_manual_connection           = true
    private_connection_resource_id = azurerm_automation_account.test.id
    subresource_names              = ["Webhook"]
    request_message                = "abc"
  }
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (a PrivateEndpointConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`




%s

data "azurerm_automation_account" "test" {
  name                = azurerm_automation_account.test.name
  resource_group_name = azurerm_resource_group.test.name
  depends_on          = [azurerm_private_endpoint.test]
}

resource "azurerm_automation_private_endpoint_connection" "test" {
  name                    = data.azurerm_automation_account.test.private_endpoint_connection[0].name
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  link_status             = "Approved"
  link_description        = "test description"

  depends_on = [azurerm_private_endpoint.test]
}
`, a.template(data))
}

func (a PrivateEndpointConnectionResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

data "azurerm_automation_account" "test" {
  name                = azurerm_automation_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_automation_private_endpoint_connection" "test" {
  name                    = azurerm_automation_account.test.private_endpoint_connection[0].name
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  link_status             = "Rejected"
  link_description        = "approved 2"
}
`, a.template(data))
}

func TestAccPrivateEndpointConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.PrivateEndpointConnectionResource{}.ResourceType(), "test")
	r := PrivateEndpointConnectionResource{}
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

func TestAccPrivateEndpointConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.PrivateEndpointConnectionResource{}.ResourceType(), "test")
	r := PrivateEndpointConnectionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}
