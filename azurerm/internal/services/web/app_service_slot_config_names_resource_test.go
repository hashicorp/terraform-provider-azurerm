package web_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SlotConfigNamesResource struct{}

func TestAccSlotConfigNames_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_config_names", "test")
	r := SlotConfigNamesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("slot_config_names.0.app_setting_names.0").HasValue("key2"),
				check.That(data.ResourceName).Key("slot_config_names.0.connection_string_names.0").HasValue("Database2"),
			),
		},
	})
}

// actual slot config sticky effect will need to be tested on azurerm_app_service_slot resource
func TestAccSlotConfigNames_slotConfigSticky(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test002")
	r := AppServiceSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.slotConfigSticky(data),
			// after slot swap its config strings will be updated
			ExpectNonEmptyPlan: true,
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.key2").DoesNotExist(),
			),
		},
	})
}

func (r SlotConfigNamesResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id := state.ID

	if id == "" {
		return nil, fmt.Errorf("cannot read slot config name resource id")
	}

	resp, err := clients.Web.AppServicesClient.ListSlotConfigurationNames(ctx, state.Attributes["resource_group_name"], id)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving slot config names from %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r SlotConfigNamesResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-001"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  app_settings = {
    "key1" = "value1"
    "key2" = "value2"
  }

  connection_string {
    name  = "Database1"
    type  = "SQLServer"
    value = "Server=some-server1.mydomain.com;Integrated Security=SSPI"
  }

  connection_string {
    name  = "Database2"
    type  = "SQLServer"
    value = "Server=some-server2.mydomain.com;Integrated Security=SSPI"
  }
}

resource "azurerm_app_service_slot_config_names" "test" {
  app_service_name    = azurerm_app_service.test.name
  resource_group_name = azurerm_resource_group.test.name
  slot_config_names {
    app_setting_names       = ["key2"]
    connection_string_names = ["Database2"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r AppServiceSlotResource) slotConfigSticky(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-002"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  app_settings = {
    "key1" = "value1"
    "key2" = "value2"
  }

  connection_string {
    name  = "Database1"
    type  = "SQLServer"
    value = "Server=some-server1.mydomain.com;Integrated Security=SSPI"
  }

  connection_string {
    name  = "Database2"
    type  = "SQLServer"
    value = "Server=some-server2.mydomain.com;Integrated Security=SSPI"
  }
}

resource "azurerm_app_service_slot_config_names" "test" {
  app_service_name    = azurerm_app_service.test.name
  resource_group_name = azurerm_resource_group.test.name
  slot_config_names {
    app_setting_names       = ["key2"]
    connection_string_names = ["Database2"]
  }
  depends_on = [
    azurerm_app_service_slot.test001,
    azurerm_app_service_slot.test002
  ]
}

resource "azurerm_app_service_slot" "test001" {
  name                = "testslot001"
  app_service_name    = azurerm_app_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  app_settings = {
    "key1" = "value1"
    "key2" = "value2"
  }

  connection_string {
    name  = "Database1"
    type  = "SQLServer"
    value = "Server=some-server1.mydomain.com;Integrated Security=SSPI"
  }

  connection_string {
    name  = "Database2"
    type  = "SQLServer"
    value = "Server=some-server2.mydomain.com;Integrated Security=SSPI"
  }
}

resource "azurerm_app_service_slot" "test002" {
  name                = "testslot002"
  app_service_name    = azurerm_app_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}
resource "azurerm_app_service_active_slot" "test001" {
  resource_group_name   = azurerm_resource_group.test.name
  app_service_name      = azurerm_app_service.test.name
  app_service_slot_name = azurerm_app_service_slot.test001.name

  depends_on = [
    azurerm_app_service_slot_config_names.test
  ]
}

resource "azurerm_app_service_active_slot" "test002" {
  resource_group_name   = azurerm_resource_group.test.name
  app_service_name      = azurerm_app_service.test.name
  app_service_slot_name = azurerm_app_service_slot.test002.name

  depends_on = [
    azurerm_app_service_active_slot.test001
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
