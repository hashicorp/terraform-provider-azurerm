// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WebAppActiveSlotResource struct{}

func TestWebAppAccActiveSlot_basicWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_app_active_slot", "test")
	r := WebAppActiveSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWindows(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestWebAppAccActiveSlot_basicLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_app_active_slot", "test")
	r := WebAppActiveSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestWebAppAccActiveSlot_windowsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_app_active_slot", "test")
	r := WebAppActiveSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWindows(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.windowsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestWebAppAccActiveSlot_linuxUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_app_active_slot", "test")
	r := WebAppActiveSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.linuxUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r WebAppActiveSlotResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseWebAppID(state.ID)
	if err != nil {
		return nil, err
	}
	slotId, err := webapps.ParseSlotID(state.Attributes["slot_id"])
	if err != nil {
		return nil, err
	}

	app, err := client.AppService.WebAppsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retreiving Function App %s for slot %s: %+v", id, slotId.SlotName, err)
	}
	if app.Model.Properties == nil || app.Model.Properties.SlotSwapStatus == nil || app.Model.Properties.SlotSwapStatus.SourceSlotName == nil {
		return nil, fmt.Errorf("missing App Slot Properties for %s", id)
	}

	return utils.Bool(*app.Model.Properties.SlotSwapStatus.SourceSlotName == slotId.SlotName), nil
}

func (r WebAppActiveSlotResource) basicWindows(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_web_app_active_slot" "test" {
  slot_id = azurerm_windows_web_app_slot.test.id
}

`, r.templateWindows(data))
}

func (r WebAppActiveSlotResource) basicLinux(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_web_app_active_slot" "test" {
  slot_id = azurerm_linux_web_app_slot.test.id
}

`, r.templateLinux(data))
}

func (r WebAppActiveSlotResource) windowsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "update" {
  name           = "acctestWAS2-%[2]d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}
}

resource "azurerm_web_app_active_slot" "test" {
  slot_id = azurerm_windows_web_app_slot.update.id
}

`, r.templateWindows(data), data.RandomInteger)
}

func (r WebAppActiveSlotResource) linuxUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "update" {
  name           = "acctestWAS2-%[2]d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}
}

resource "azurerm_web_app_active_slot" "test" {
  slot_id = azurerm_linux_web_app_slot.update.id
}

`, r.templateLinux(data), data.RandomInteger)
}

func (WebAppActiveSlotResource) templateLinux(data acceptance.TestData) string {
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
  name           = "acctestWAS-%[1]d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary)
}

func (WebAppActiveSlotResource) templateWindows(data acceptance.TestData) string {
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

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%[1]d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary)
}
