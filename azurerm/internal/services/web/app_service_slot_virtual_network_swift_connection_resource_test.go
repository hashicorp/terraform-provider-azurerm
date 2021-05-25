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

type AppServiceSlotVirtualNetworkSwiftConnectionResource struct {
}

func TestAccAppServiceSlotVirtualNetworkSwiftConnection_app_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_virtual_network_swift_connection", "test")
	r := AppServiceSlotVirtualNetworkSwiftConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.app_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceSlotVirtualNetworkSwiftConnection_app_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_virtual_network_swift_connection", "test")
	r := AppServiceSlotVirtualNetworkSwiftConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.app_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.app_requiresImport),
	})
}

func TestAccAppServiceSlotVirtualNetworkSwiftConnection_app_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_virtual_network_swift_connection", "test")
	r := AppServiceSlotVirtualNetworkSwiftConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.app_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
			),
		},
		{
			Config: r.app_update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
			),
		},
	})
}

func TestAccAppServiceSlotVirtualNetworkSwiftConnection_app_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_virtual_network_swift_connection", "test")
	r := AppServiceSlotVirtualNetworkSwiftConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.app_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				data.CheckWithClient(r.disappears),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccFunctionAppSlotVirtualNetworkSwiftConnection_function_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_virtual_network_swift_connection", "test")
	r := AppServiceSlotVirtualNetworkSwiftConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.function_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionAppSlotVirtualNetworkSwiftConnection_function_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_virtual_network_swift_connection", "test")
	r := AppServiceSlotVirtualNetworkSwiftConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.function_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.function_requiresImport),
	})
}

func TestAccFunctionAppSlotVirtualNetworkSwiftConnection_function_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_virtual_network_swift_connection", "test")
	r := AppServiceSlotVirtualNetworkSwiftConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.function_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.function_update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccFunctionAppSlotVirtualNetworkSwiftConnection_function_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_virtual_network_swift_connection", "test")
	r := AppServiceSlotVirtualNetworkSwiftConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.function_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.disappears),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func (r AppServiceSlotVirtualNetworkSwiftConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SlotVirtualNetworkSwiftConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Web.AppServicesClient.GetSwiftVirtualNetworkConnectionSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id.String(), err)
	}

	return utils.Bool(resp.SwiftVirtualNetworkProperties != nil), nil
}

func (t AppServiceSlotVirtualNetworkSwiftConnectionResource) disappears(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := parse.SlotVirtualNetworkSwiftConnectionID(state.ID)
	if err != nil {
		return err
	}

	resp, err := clients.Web.AppServicesClient.DeleteSwiftVirtualNetworkSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", id.String(), err)
		}
	}

	return nil
}

func (AppServiceSlotVirtualNetworkSwiftConnectionResource) app_base(data acceptance.TestData) string {
	return fmt.Sprintf(`
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
    ignore_changes = ["ddos_protection_plan"]
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

resource "azurerm_app_service_slot" "test-staging" {
  name                = "acctest-AS-%d-staging"
  app_service_name    = azurerm_app_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceSlotVirtualNetworkSwiftConnectionResource) app_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot_virtual_network_swift_connection" "test" {
  slot_name      = azurerm_app_service_slot.test-staging.name
  app_service_id = azurerm_app_service.test.id
  subnet_id      = azurerm_subnet.test1.id
}
`, r.app_base(data))
}

func (r AppServiceSlotVirtualNetworkSwiftConnectionResource) app_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot_virtual_network_swift_connection" "test" {
  slot_name      = azurerm_app_service_slot.test-staging.name
  app_service_id = azurerm_app_service.test.id
  subnet_id      = azurerm_subnet.test2.id
}
`, r.app_base(data))
}

func (r AppServiceSlotVirtualNetworkSwiftConnectionResource) app_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot_virtual_network_swift_connection" "import" {
  slot_name      = azurerm_app_service_slot_virtual_network_swift_connection.test.slot_name
  app_service_id = azurerm_app_service_slot_virtual_network_swift_connection.test.app_service_id
  subnet_id      = azurerm_app_service_slot_virtual_network_swift_connection.test.subnet_id
}
`, r.app_basic(data))
}

func (AppServiceSlotVirtualNetworkSwiftConnectionResource) function_base(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-functionapp-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lifecycle {
    ignore_changes = ["ddos_protection_plan"]
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
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test-staging" {
  name                       = "acctest-FA-%d-staging"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceSlotVirtualNetworkSwiftConnectionResource) function_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot_virtual_network_swift_connection" "test" {
  slot_name      = azurerm_function_app_slot.test-staging.name
  app_service_id = azurerm_function_app.test.id
  subnet_id      = azurerm_subnet.test1.id
}
`, r.function_base(data))
}

func (r AppServiceSlotVirtualNetworkSwiftConnectionResource) function_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot_virtual_network_swift_connection" "test" {
  slot_name      = azurerm_function_app_slot.test-staging.name
  app_service_id = azurerm_function_app.test.id
  subnet_id      = azurerm_subnet.test2.id
}
`, r.function_base(data))
}

func (r AppServiceSlotVirtualNetworkSwiftConnectionResource) function_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot_virtual_network_swift_connection" "import" {
  slot_name      = azurerm_app_service_slot_virtual_network_swift_connection.test.slot_name
  app_service_id = azurerm_app_service_slot_virtual_network_swift_connection.test.app_service_id
  subnet_id      = azurerm_app_service_slot_virtual_network_swift_connection.test.subnet_id
}
`, r.function_basic(data))
}
