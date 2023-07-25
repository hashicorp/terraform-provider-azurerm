// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maintenance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MaintenanceConfigurationResource struct{}

func TestAccMaintenanceConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_configuration", "test")
	r := MaintenanceConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue("SQLDB"),
				check.That(data.ResourceName).Key("visibility").HasValue("Custom"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMaintenanceConfiguration_basicWithInGuestPatch(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_configuration", "test")
	r := MaintenanceConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_withInGuestPatch(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMaintenanceConfiguration_basicWithOnePatchOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_configuration", "test")
	r := MaintenanceConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_onePatchOnly(data, true),
		},
		data.ImportStep(),
		{
			Config: r.basic_onePatchOnly(data, false),
		},
		data.ImportStep(),
	})
}

func TestAccMaintenanceConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_configuration", "test")
	r := MaintenanceConfigurationResource{}

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

func TestAccMaintenanceConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_configuration", "test")
	r := MaintenanceConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue("SQLDB"),
				check.That(data.ResourceName).Key("visibility").HasValue("Custom"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.enV").HasValue("TesT"),
				check.That(data.ResourceName).Key("window.0.start_date_time").HasValue("5555-12-31 00:00"),
				check.That(data.ResourceName).Key("window.0.expiration_date_time").HasValue("6666-12-31 00:00"),
				check.That(data.ResourceName).Key("window.0.duration").HasValue("06:00"),
				check.That(data.ResourceName).Key("window.0.time_zone").HasValue("Pacific Standard Time"),
				check.That(data.ResourceName).Key("window.0.recur_every").HasValue("2Days"),
				check.That(data.ResourceName).Key("properties.%").HasValue("1"),
				check.That(data.ResourceName).Key("properties.description").HasValue("acceptance test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMaintenanceConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_configuration", "test")
	r := MaintenanceConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue("SQLDB"),
				check.That(data.ResourceName).Key("visibility").HasValue("Custom"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("window.#").HasValue("0"),
				check.That(data.ResourceName).Key("properties.%").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue("SQLDB"),
				check.That(data.ResourceName).Key("visibility").HasValue("Custom"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.enV").HasValue("TesT"),
				check.That(data.ResourceName).Key("window.0.start_date_time").HasValue("5555-12-31 00:00"),
				check.That(data.ResourceName).Key("window.0.expiration_date_time").HasValue("6666-12-31 00:00"),
				check.That(data.ResourceName).Key("window.0.duration").HasValue("06:00"),
				check.That(data.ResourceName).Key("window.0.time_zone").HasValue("Pacific Standard Time"),
				check.That(data.ResourceName).Key("window.0.recur_every").HasValue("2Days"),
				check.That(data.ResourceName).Key("properties.%").HasValue("1"),
				check.That(data.ResourceName).Key("properties.description").HasValue("acceptance test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue("SQLDB"),
				check.That(data.ResourceName).Key("visibility").HasValue("Custom"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("window.#").HasValue("0"),
				check.That(data.ResourceName).Key("properties.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func (MaintenanceConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := maintenanceconfigurations.ParseMaintenanceConfigurationIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Maintenance.ConfigurationsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return utils.Bool(resp.Model != nil && resp.Model.Properties != nil), nil
}

func (MaintenanceConfigurationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-maint-%d"
  location = "%s"
}

resource "azurerm_maintenance_configuration" "test" {
  name                = "acctest-MC%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope               = "SQLDB"
  visibility          = "Custom"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (MaintenanceConfigurationResource) basic_withInGuestPatch(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-maint-%d"
  location = "%s"
}

resource "azurerm_maintenance_configuration" "test" {
  name                = "acctest-MC%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope               = "InGuestPatch"
  visibility          = "Custom"

  window {
    start_date_time      = "5555-12-31 00:00"
    expiration_date_time = "6666-12-31 00:00"
    duration             = "02:00"
    time_zone            = "Pacific Standard Time"
    recur_every          = "2Days"
  }

  install_patches {
    reboot = "IfRequired"
    linux {
      classifications_to_include = ["Critical", "Security"]
    }
    windows {
      classifications_to_include = ["Critical", "Security"]
    }
  }

  in_guest_user_patch_mode = "User"

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r MaintenanceConfigurationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_configuration" "import" {
  name                = azurerm_maintenance_configuration.test.name
  resource_group_name = azurerm_maintenance_configuration.test.resource_group_name
  location            = azurerm_maintenance_configuration.test.location
  scope               = azurerm_maintenance_configuration.test.scope
  visibility          = azurerm_maintenance_configuration.test.visibility
}
`, r.basic(data))
}

func (MaintenanceConfigurationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-maint-%d"
  location = "%s"
}

resource "azurerm_maintenance_configuration" "test" {
  name                = "acctest-MC%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope               = "SQLDB"
  visibility          = "Custom"

  window {
    start_date_time      = "5555-12-31 00:00"
    expiration_date_time = "6666-12-31 00:00"
    duration             = "06:00"
    time_zone            = "Pacific Standard Time"
    recur_every          = "2Days"
  }

  properties = {
    description = "acceptance test"
  }

  tags = {
    enV = "TesT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (MaintenanceConfigurationResource) basic_onePatchOnly(data acceptance.TestData, isLinux bool) string {
	patch := `linux {
      classifications_to_include = ["Critical", "Security"]
    }`
	if !isLinux {
		patch = `windows {
      classifications_to_include = ["Critical", "Security"]
    }`
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-maint-%d"
  location = "%s"
}

resource "azurerm_maintenance_configuration" "test" {
  name                = "acctest-MC%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope               = "InGuestPatch"
  visibility          = "Custom"

  window {
    start_date_time      = "5555-12-31 00:00"
    expiration_date_time = "6666-12-31 00:00"
    duration             = "02:00"
    time_zone            = "Pacific Standard Time"
    recur_every          = "2Days"
  }

  install_patches {
    reboot = "IfRequired"
    %s
  }

  in_guest_user_patch_mode = "User"

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, patch)
}
