// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualDesktopScalingPlanResource struct{}

func TestAccVirtualDesktopScalingPlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_scaling_plan", "test")
	r := VirtualDesktopScalingPlanResource{}
	roleAssignmentId := uuid.New().String()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, roleAssignmentId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccVirtualDesktopScalingPlan_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_scaling_plan", "test")
	r := VirtualDesktopScalingPlanResource{}
	roleAssignmentId := uuid.New().String()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, roleAssignmentId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func TestAccVirtualDesktopScalingPlan_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_scaling_plan", "test")
	r := VirtualDesktopScalingPlanResource{}
	roleAssignmentId := uuid.New().String()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, roleAssignmentId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config: r.complete(data, roleAssignmentId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		{
			Config: r.basic(data, roleAssignmentId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccVirtualDesktopScalingPlan_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_scaling_plan", "test")
	r := VirtualDesktopScalingPlanResource{}
	roleAssignmentId := uuid.New().String()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, roleAssignmentId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, roleAssignmentId),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_desktop_scaling_plan"),
		},
	})
}

func (VirtualDesktopScalingPlanResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := scalingplan.ParseScalingPlanID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DesktopVirtualization.ScalingPlansClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (VirtualDesktopScalingPlanResource) basic(data acceptance.TestData, roleAssignmentId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}


resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}


data "azuread_service_principal" "test" {
  display_name = "Windows Virtual Desktop"
}

resource "azurerm_role_assignment" "test" {
  name                             = "%s"
  scope                            = azurerm_resource_group.test.id
  role_definition_name             = "Desktop Virtualization Power On Off Contributor"
  principal_id                     = data.azuread_service_principal.test.object_id
  skip_service_principal_aad_check = true
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                 = "acctestHP%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  type                 = "Pooled"
  validate_environment = true
  load_balancer_type   = "BreadthFirst"
}

resource "azurerm_virtual_desktop_scaling_plan" "test" {
  name                = "scalingPlan%x"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  friendly_name       = "Scaling Plan Test"
  description         = "Test Scaling Plan"
  time_zone           = "GMT Standard Time"
  schedule {
    name                                 = "Weekdays"
    days_of_week                         = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    ramp_up_start_time                   = "06:00"
    ramp_up_load_balancing_algorithm     = "BreadthFirst"
    ramp_up_minimum_hosts_percent        = 20
    ramp_up_capacity_threshold_percent   = 10
    peak_start_time                      = "09:00"
    peak_load_balancing_algorithm        = "BreadthFirst"
    ramp_down_start_time                 = "18:00"
    ramp_down_load_balancing_algorithm   = "BreadthFirst"
    ramp_down_minimum_hosts_percent      = 10
    ramp_down_force_logoff_users         = false
    ramp_down_wait_time_minutes          = 45
    ramp_down_notification_message       = "Please log of in the next 45 minutes..."
    ramp_down_capacity_threshold_percent = 5
    ramp_down_stop_hosts_when            = "ZeroSessions"
    off_peak_start_time                  = "22:00"
    off_peak_load_balancing_algorithm    = "BreadthFirst"
  }

  depends_on = [azurerm_role_assignment.test]


}
`, data.RandomInteger, data.Locations.Primary, roleAssignmentId, data.RandomString, data.RandomString)
}

func (VirtualDesktopScalingPlanResource) complete(data acceptance.TestData, roleAssignmentId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}

data "azuread_service_principal" "test" {
  display_name = "Windows Virtual Desktop"
}

resource "azurerm_role_assignment" "test" {
  name                             = "%s"
  scope                            = azurerm_resource_group.test.id
  role_definition_name             = "Desktop Virtualization Power On Off Contributor"
  principal_id                     = data.azuread_service_principal.test.object_id
  skip_service_principal_aad_check = true
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                 = "acctestHP%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  type                 = "Pooled"
  validate_environment = true
  load_balancer_type   = "BreadthFirst"
}

resource "azurerm_virtual_desktop_scaling_plan" "test" {
  name                = "scalingPlan%x"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  friendly_name       = "Scaling Plan Test"
  description         = "Test Scaling Plan"
  time_zone           = "GMT Standard Time"

  schedule {
    name                                 = "Weekdays"
    days_of_week                         = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    ramp_up_start_time                   = "05:00"
    ramp_up_load_balancing_algorithm     = "BreadthFirst"
    ramp_up_minimum_hosts_percent        = 20
    ramp_up_capacity_threshold_percent   = 10
    peak_start_time                      = "09:00"
    peak_load_balancing_algorithm        = "BreadthFirst"
    ramp_down_start_time                 = "19:00"
    ramp_down_load_balancing_algorithm   = "BreadthFirst"
    ramp_down_minimum_hosts_percent      = 10
    ramp_down_force_logoff_users         = false
    ramp_down_wait_time_minutes          = 45
    ramp_down_notification_message       = "Please log of in the next 45 minutes..."
    ramp_down_capacity_threshold_percent = 5
    ramp_down_stop_hosts_when            = "ZeroSessions"
    off_peak_start_time                  = "22:00"
    off_peak_load_balancing_algorithm    = "BreadthFirst"
  }

  schedule {
    name                                 = "Weekend"
    days_of_week                         = ["Saturday", "Sunday"]
    ramp_up_start_time                   = "09:00"
    ramp_up_load_balancing_algorithm     = "BreadthFirst"
    ramp_up_minimum_hosts_percent        = 30
    ramp_up_capacity_threshold_percent   = 10
    peak_start_time                      = "10:00"
    peak_load_balancing_algorithm        = "BreadthFirst"
    ramp_down_start_time                 = "16:00"
    ramp_down_load_balancing_algorithm   = "BreadthFirst"
    ramp_down_minimum_hosts_percent      = 10
    ramp_down_force_logoff_users         = false
    ramp_down_wait_time_minutes          = 45
    ramp_down_notification_message       = "Please log of in the next 45 minutes..."
    ramp_down_capacity_threshold_percent = 5
    ramp_down_stop_hosts_when            = "ZeroSessions"
    off_peak_start_time                  = "20:00"
    off_peak_load_balancing_algorithm    = "BreadthFirst"
  }


  tags = {
    Acceptance = "Test"
  }

  depends_on = [azurerm_role_assignment.test]
}

`, data.RandomInteger, data.Locations.Primary, roleAssignmentId, data.RandomString, data.RandomString)
}

func (r VirtualDesktopScalingPlanResource) requiresImport(data acceptance.TestData, roleAssignmentId string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_scaling_plan" "import" {
  name                = azurerm_virtual_desktop_scaling_plan.test.name
  location            = azurerm_virtual_desktop_scaling_plan.test.location
  resource_group_name = azurerm_virtual_desktop_scaling_plan.test.resource_group_name
  friendly_name       = azurerm_virtual_desktop_scaling_plan.test.friendly_name
  description         = azurerm_virtual_desktop_scaling_plan.test.description
  time_zone           = azurerm_virtual_desktop_scaling_plan.test.time_zone

  schedule {
    name                                 = "Weekdays"
    days_of_week                         = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    ramp_up_start_time                   = "06:00"
    ramp_up_load_balancing_algorithm     = "BreadthFirst"
    ramp_up_minimum_hosts_percent        = 20
    ramp_up_capacity_threshold_percent   = 10
    peak_start_time                      = "09:00"
    peak_load_balancing_algorithm        = "BreadthFirst"
    ramp_down_start_time                 = "18:00"
    ramp_down_load_balancing_algorithm   = "DepthFirst"
    ramp_down_minimum_hosts_percent      = 10
    ramp_down_force_logoff_users         = false
    ramp_down_wait_time_minutes          = 45
    ramp_down_notification_message       = "Please log of in the next 45 minutes..."
    ramp_down_capacity_threshold_percent = 5
    ramp_down_stop_hosts_when            = "ZeroSessions"
    off_peak_start_time                  = "22:00"
    off_peak_load_balancing_algorithm    = "DepthFirst"
  }
}
`, r.basic(data, roleAssignmentId))
}
