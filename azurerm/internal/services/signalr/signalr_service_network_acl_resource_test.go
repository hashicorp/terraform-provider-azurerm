package signalr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/signalr/sdk/signalr"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SignalRServiceNetworkACLResource struct{}

func TestAccSignalRServiceNetworkACL_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service_network_acl", "test")
	r := SignalRServiceNetworkACLResource{}

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

func TestAccSignalRServiceNetworkACL_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service_network_acl", "test")
	r := SignalRServiceNetworkACLResource{}

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

func TestAccSignalRServiceNetworkACL_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service_network_acl", "test")
	r := SignalRServiceNetworkACLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSignalRServiceNetworkACL_updateMultiplePrivateEndpoints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service_network_acl", "test")
	r := SignalRServiceNetworkACLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiplePrivateEndpoints(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SignalRServiceNetworkACLResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := signalr.ParseSignalRID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SignalR.Client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	isDefaultConfiguration := false
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if acls := props.NetworkACLs; acls != nil {
				hasDefaultAction := false
				if acls.DefaultAction != nil {
					hasDefaultAction = *acls.DefaultAction == signalr.ACLActionDeny
				}

				hasDefaultMatches := false
				if acls.PublicNetwork != nil && acls.PublicNetwork.Allow != nil {
					hasDefaultMatches = len(*acls.PublicNetwork.Allow) == 4
				}

				isDefaultConfiguration = hasDefaultAction && hasDefaultMatches
			}
		}
	}

	return utils.Bool(!isDefaultConfiguration), nil
}

func (r SignalRServiceNetworkACLResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_signalr_service_network_acl" "test" {
  signalr_service_id = azurerm_signalr_service.test.id
  default_action     = "Deny"

  public_network {
    allowed_request_types = ["ClientConnection"]
  }
}
`, r.template(data))
}

func (r SignalRServiceNetworkACLResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-pe-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.test.id

  private_service_connection {
    name                           = "psc-sig-test"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_signalr_service.test.id
    subresource_names              = ["signalr"]
  }
}

resource "azurerm_signalr_service_network_acl" "test" {
  signalr_service_id = azurerm_signalr_service.test.id
  default_action     = "Allow"

  public_network {
    denied_request_types = ["ClientConnection"]
  }

  private_endpoint {
    id                   = azurerm_private_endpoint.test.id
    denied_request_types = ["ClientConnection"]
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r SignalRServiceNetworkACLResource) multiplePrivateEndpoints(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-pe-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.test.id

  private_service_connection {
    name                           = "psc-sig-test"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_signalr_service.test.id
    subresource_names              = ["signalr"]
  }
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctest-vnet2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctest-subnet2-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_private_endpoint" "test2" {
  name                = "acctest-pe2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.test2.id

  private_service_connection {
    name                           = "psc-sig-test2"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_signalr_service.test.id
    subresource_names              = ["signalr"]
  }
}

resource "azurerm_signalr_service_network_acl" "test" {
  signalr_service_id = azurerm_signalr_service.test.id
  default_action     = "Allow"

  public_network {
    denied_request_types = ["ClientConnection"]
  }

  private_endpoint {
    id                   = azurerm_private_endpoint.test.id
    denied_request_types = ["ClientConnection"]
  }

  private_endpoint {
    id                   = azurerm_private_endpoint.test2.id
    denied_request_types = ["ServerConnection"]
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r SignalRServiceNetworkACLResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-signalr-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctest-signalr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_S1"
    capacity = 1
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
