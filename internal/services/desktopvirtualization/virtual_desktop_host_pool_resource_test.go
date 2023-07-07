// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualDesktopHostPoolResource struct{}

func TestAccVirtualDesktopHostPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_host_pool", "test")
	r := VirtualDesktopHostPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccVirtualDesktopHostPool_agentupdates(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_host_pool", "test")
	r := VirtualDesktopHostPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.agentUpdateBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("scheduled_agent_updates.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("scheduled_agent_updates.0.schedule.0.day_of_week").HasValue("Saturday"),
			),
		},
		{
			Config: r.agentUpdateComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("scheduled_agent_updates.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("scheduled_agent_updates.0.use_session_host_timezone").HasValue("true"),
				check.That(data.ResourceName).Key("scheduled_agent_updates.0.schedule.0.day_of_week").HasValue("Saturday"),
				check.That(data.ResourceName).Key("scheduled_agent_updates.0.schedule.1.day_of_week").HasValue("Sunday"),
			),
		},
		{
			Config: r.agentUpdateDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("scheduled_agent_updates.0.enabled").HasValue("false"),
			),
		},
	})
}

func TestAccVirtualDesktopHostPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_host_pool", "test")
	r := VirtualDesktopHostPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func TestAccVirtualDesktopHostPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_host_pool", "test")
	r := VirtualDesktopHostPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccVirtualDesktopHostPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_host_pool", "test")
	r := VirtualDesktopHostPoolResource{}

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

func (VirtualDesktopHostPoolResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := hostpool.ParseHostPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DesktopVirtualization.HostPoolsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (VirtualDesktopHostPoolResource) agentUpdateBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktophp-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                 = "acctestHP%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  type                 = "Pooled"
  validate_environment = true
  load_balancer_type   = "BreadthFirst"
  scheduled_agent_updates {
    enabled = true
    schedule {
      day_of_week = "Saturday"
      hour_of_day = 2
    }
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString)
}

func (VirtualDesktopHostPoolResource) agentUpdateComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktophp-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                 = "acctestHP%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  type                 = "Pooled"
  validate_environment = true
  load_balancer_type   = "BreadthFirst"
  scheduled_agent_updates {
    enabled                   = true
    use_session_host_timezone = true
    schedule {
      day_of_week = "Saturday"
      hour_of_day = 2
    }
    schedule {
      day_of_week = "Sunday"
      hour_of_day = 2
    }
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString)
}

func (VirtualDesktopHostPoolResource) agentUpdateDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktophp-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                 = "acctestHP%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  type                 = "Pooled"
  validate_environment = true
  load_balancer_type   = "BreadthFirst"
  scheduled_agent_updates {
    enabled                   = false
    use_session_host_timezone = true
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString)
}

func (VirtualDesktopHostPoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktophp-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                 = "acctestHP%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  type                 = "Pooled"
  validate_environment = true
  load_balancer_type   = "DepthFirst"
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString)
}

func (VirtualDesktopHostPoolResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktophp-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                     = "acctestHP%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  type                     = "Pooled"
  friendly_name            = "A Friendly Name!"
  description              = "A Description!"
  validate_environment     = true
  start_vm_on_connect      = true
  load_balancer_type       = "BreadthFirst"
  maximum_sessions_allowed = 100
  preferred_app_group_type = "Desktop"
  custom_rdp_properties    = "audiocapturemode:i:1;audiomode:i:0;"

  tags = {
    Purpose = "Acceptance-Testing"
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString)
}

func (r VirtualDesktopHostPoolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_host_pool" "import" {
  name                 = azurerm_virtual_desktop_host_pool.test.name
  location             = azurerm_virtual_desktop_host_pool.test.location
  resource_group_name  = azurerm_virtual_desktop_host_pool.test.resource_group_name
  validate_environment = azurerm_virtual_desktop_host_pool.test.validate_environment
  type                 = azurerm_virtual_desktop_host_pool.test.type
  load_balancer_type   = azurerm_virtual_desktop_host_pool.test.load_balancer_type
}
`, r.basic(data))
}
