package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ActiveSlotResource struct{}

func TestAccActiveSlot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_app_active_slot", "test")
	r := ActiveSlotResource{}

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

func (a ActiveSlotResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WebAppSlotID(state.ID)
	if err != nil {
		return nil, err
	}

	app, err := client.AppService.WebAppsClient.Get(ctx, id.ResourceGroup, id.SiteName)
	if app.SiteProperties == nil || app.SiteProperties.SlotSwapStatus == nil || app.SiteProperties.SlotSwapStatus.SourceSlotName == nil {
		return nil, fmt.Errorf("missing App Slot Properties for %s", id)
	}

	return utils.Bool(*app.SiteProperties.SlotSwapStatus.SourceSlotName == id.SlotName), nil
}

func (r ActiveSlotResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_web_app_active_slot" "test" {
  slot_id = azurerm_linux_web_app_slot.test.id
}

`, r.template(data))
}

func (ActiveSlotResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-WAS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "S1"
}

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}

resource "azurerm_linux_web_app_slot" "test" {
  name                = "acctestWAS-%[1]d"
  app_service_name    = azurerm_linux_web_app.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary)
}
