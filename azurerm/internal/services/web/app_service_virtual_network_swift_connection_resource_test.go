package web_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AppServiceVirtualNetworkSwiftConnectionResource struct {
}

func TestAccAppServiceVirtualNetworkSwiftConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_virtual_network_swift_connection", "test")
	r := AppServiceVirtualNetworkSwiftConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceVirtualNetworkSwiftConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_virtual_network_swift_connection", "test")
	r := AppServiceVirtualNetworkSwiftConnectionResource{}

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

func TestAccAppServiceVirtualNetworkSwiftConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_virtual_network_swift_connection", "test")
	r := AppServiceVirtualNetworkSwiftConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
			),
		},
	})
}

func TestAccAppServiceVirtualNetworkSwiftConnection_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_virtual_network_swift_connection", "test")
	r := AppServiceVirtualNetworkSwiftConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				data.CheckWithClient(r.disappears),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func (r AppServiceVirtualNetworkSwiftConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VirtualNetworkSwiftConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Web.AppServicesClient.GetSwiftVirtualNetworkConnection(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id.String(), err)
	}

	return utils.Bool(resp.SwiftVirtualNetworkProperties != nil), nil
}

func (t AppServiceVirtualNetworkSwiftConnectionResource) disappears(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := parse.VirtualNetworkSwiftConnectionID(state.ID)
	if err != nil {
		return err
	}

	resp, err := clients.Web.AppServicesClient.DeleteSwiftVirtualNetwork(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", id.String(), err)
		}
	}

	return nil
}

func TestAccAppServiceVirtualNetworkSwiftConnection_functionAppBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_virtual_network_swift_connection", "test")
	r := AppServiceVirtualNetworkSwiftConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.functionAppBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (AppServiceVirtualNetworkSwiftConnectionResource) base(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appservice-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lifecycle {
    ignore_changes = [ddos_protection_plan]
  }
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestSubnet1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestSubnet2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctest-ASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctest-AS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (AppServiceVirtualNetworkSwiftConnectionResource) basic(data acceptance.TestData) string {
	template := AppServiceVirtualNetworkSwiftConnectionResource{}.base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_virtual_network_swift_connection" "test" {
  app_service_id = azurerm_app_service.test.id
  subnet_id      = azurerm_subnet.test1.id
}
`, template)
}

func (AppServiceVirtualNetworkSwiftConnectionResource) update(data acceptance.TestData) string {
	template := AppServiceVirtualNetworkSwiftConnectionResource{}.base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_virtual_network_swift_connection" "test" {
  app_service_id = azurerm_app_service.test.id
  subnet_id      = azurerm_subnet.test2.id
}
`, template)
}

func (AppServiceVirtualNetworkSwiftConnectionResource) requiresImport(data acceptance.TestData) string {
	template := AppServiceVirtualNetworkSwiftConnectionResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_virtual_network_swift_connection" "import" {
  app_service_id = azurerm_app_service_virtual_network_swift_connection.test.app_service_id
  subnet_id      = azurerm_app_service_virtual_network_swift_connection.test.subnet_id
}
`, template)
}

func (AppServiceVirtualNetworkSwiftConnectionResource) functionAppBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appservice-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lifecycle {
    ignore_changes = [ddos_protection_plan]
  }
}

resource "azurerm_subnet" "test" {
  name                 = "acctestSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctest-ASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-FA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_app_service_virtual_network_swift_connection" "test" {
  app_service_id = azurerm_function_app.test.id
  subnet_id      = azurerm_subnet.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}
