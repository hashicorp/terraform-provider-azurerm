// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maintenance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MaintenanceDynamicScopeResource struct{}

func TestAccMaintenanceAssignmentDynamicScope_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_dynamic_scope", "test")
	r := MaintenanceDynamicScopeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("filter.0.locations.0").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("filter.0.resource_types.0").HasValue("Microsoft.Compute/virtualMachines"),
			),
		},
		data.ImportStep("location"),
	})
}

func TestAccMaintenanceAssignmentDynamicScope_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_dynamic_scope", "test")
	r := MaintenanceDynamicScopeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("filter.#").HasValue("1"),
				check.That(data.ResourceName).Key("filter.0.locations.0").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("filter.0.os_types.0").HasValue("windows"),
				check.That(data.ResourceName).Key("filter.0.resource_types.0").HasValue("Microsoft.Compute/virtualMachines"),
				check.That(data.ResourceName).Key("filter.0.tag_filter").HasValue("Any"),
				check.That(data.ResourceName).Key("filter.0.tags.0.tag").HasValue("foo"),
				check.That(data.ResourceName).Key("filter.0.tags.0.values.0").HasValue("barbar"),
			),
		},
		data.ImportStep("location"),
	})
}

func TestAccMaintenanceAssignmentDynamicScope_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_dynamic_scope", "test")
	r := MaintenanceDynamicScopeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("filter.0.locations.0").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("filter.0.resource_types.0").HasValue("Microsoft.Compute/virtualMachines"),
			),
		},
		data.ImportStep("location"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("filter.#").HasValue("1"),
				check.That(data.ResourceName).Key("filter.0.locations.0").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("filter.0.os_types.0").HasValue("windows"),
				check.That(data.ResourceName).Key("filter.0.resource_types.0").HasValue("Microsoft.Compute/virtualMachines"),
				check.That(data.ResourceName).Key("filter.0.tag_filter").HasValue("Any"),
				check.That(data.ResourceName).Key("filter.0.tags.0.tag").HasValue("foo"),
				check.That(data.ResourceName).Key("filter.0.tags.0.values.0").HasValue("barbar"),
			),
		},
		data.ImportStep("location"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("filter.0.locations.0").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("filter.0.resource_types.0").HasValue("Microsoft.Compute/virtualMachines"),
			),
		},
		data.ImportStep("location"),
	})
}

func TestAccMaintenanceAssignmentDynamicScope_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_dynamic_scope", "test")
	r := MaintenanceDynamicScopeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("filter.0.locations.0").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("filter.0.resource_types.0").HasValue("Microsoft.Compute/virtualMachines"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (MaintenanceDynamicScopeResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := configurationassignments.ParseConfigurationAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Maintenance.ConfigurationAssignmentsClient.ForSubscriptionsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r MaintenanceDynamicScopeResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_maintenance_assignment_dynamic_scope" "test" {
  name                         = "acctest-ca%[2]d"
  maintenance_configuration_id = azurerm_maintenance_configuration.test.id
  subscription_id              = format("/subscriptions/%%s", data.azurerm_client_config.current.subscription_id)

  filter {
    locations      = ["%[3]s"]
    resource_types = ["Microsoft.Compute/virtualMachines"]
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r MaintenanceDynamicScopeResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_maintenance_assignment_dynamic_scope" "test" {
  name                         = "acctest-complete-%[2]d"
  maintenance_configuration_id = azurerm_maintenance_configuration.test.id
  subscription_id              = format("/subscriptions/%%s", data.azurerm_client_config.current.subscription_id)

  filter {
    locations       = ["%[3]s"]
    os_types        = ["windows"]
    resource_groups = [azurerm_resource_group.test.name]
    resource_types  = ["Microsoft.Compute/virtualMachines"]
    tag_filter      = "Any"
    tags {
      tag    = "foo"
      values = ["barbar"]
    }
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r MaintenanceDynamicScopeResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_maintenance_assignment_dynamic_scope" "import" {
  name                         = "acctest-ca%[2]d"
  maintenance_configuration_id = azurerm_maintenance_configuration.test.id
  subscription_id              = format("/subscriptions/%%s", data.azurerm_client_config.current.subscription_id)

  filter {
    locations = ["%[3]s"]
  }
}
`, r.basic(data), data.RandomInteger, data.Locations.Primary)
}

func (MaintenanceDynamicScopeResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-maint-%[1]d"
  location = "%[2]s"
}

resource "azurerm_maintenance_configuration" "test" {
  name                     = "acctest-MC%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  scope                    = "InGuestPatch"
  in_guest_user_patch_mode = "User"

  window {
    start_date_time = formatdate("YYYY-MM-DD hh:mm", timestamp())
    time_zone       = "Greenwich Standard Time"
    recur_every     = "1Day"
  }

  install_patches {
    reboot = "Always"

    windows {
      classifications_to_include = ["Critical"]
      kb_numbers_to_exclude      = []
      kb_numbers_to_include      = []
    }
  }

  lifecycle {
    ignore_changes = [
      window[0].start_date_time,
      window[0].duration
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
