// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/standbypool/2025-03-01/standbyvirtualmachinepools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StandbyPoolStandbyVirtualMachinePoolResource struct{}

func TestAccStandbyPoolStandbyVirtualMachinePool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_standby_pool", "test")
	r := StandbyPoolStandbyVirtualMachinePoolResource{}
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

func TestAccStandbyPoolStandbyVirtualMachinePool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_standby_pool", "test")
	r := StandbyPoolStandbyVirtualMachinePoolResource{}
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

func TestAccStandbyPoolStandbyVirtualMachinePool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_standby_pool", "test")
	r := StandbyPoolStandbyVirtualMachinePoolResource{}
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

func TestAccStandbyPoolStandbyVirtualMachinePool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_standby_pool", "test")
	r := StandbyPoolStandbyVirtualMachinePoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
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
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStandbyPoolStandbyVirtualMachinePool_minCapacityExceedMaxCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_standby_pool", "test")
	r := StandbyPoolStandbyVirtualMachinePoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.minCapacityExceedMaxCapacity(data),
			ExpectError: regexp.MustCompile("`min_ready_capacity` cannot exceed `max_ready_capacity`"),
		},
	})
}

func (r StandbyPoolStandbyVirtualMachinePoolResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := standbyvirtualmachinepools.ParseStandbyVirtualMachinePoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.StandbyVirtualMachinePoolsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StandbyPoolStandbyVirtualMachinePoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "vm-contributor" {
  name = "Virtual Machine Contributor"
}

data "azurerm_role_definition" "nw-contributor" {
  name = "Network Contributor"
}

data "azurerm_role_definition" "mi-contributor" {
  name = "Managed Identity Contributor"
}

data "azuread_service_principal" "test" {
  display_name = "Standby Pool Resource Provider"
}

resource "azurerm_role_assignment" "vm-contributor" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.vm-contributor.id}"
  principal_id       = data.azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "nw-contributor" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.nw-contributor.id}"
  principal_id       = data.azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "mi-contributor" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.mi-contributor.id}"
  principal_id       = data.azuread_service_principal.test.object_id
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_fault_domain_count = 1

  zones = ["1"]

  depends_on = [
    azurerm_role_assignment.vm-contributor, azurerm_role_assignment.nw-contributor, azurerm_role_assignment.mi-contributor
  ]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r StandbyPoolStandbyVirtualMachinePoolResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_virtual_machine_scale_set_standby_pool" "test" {
  name                                  = "acctest-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  location                              = "%s"
  attached_virtual_machine_scale_set_id = azurerm_orchestrated_virtual_machine_scale_set.test.id
  virtual_machine_state                 = "Running"

  elasticity_profile {
    max_ready_capacity = 10
    min_ready_capacity = 5
  }
}
`, template, data.RandomIntOfLength(16), data.Locations.Primary)
}

func (r StandbyPoolStandbyVirtualMachinePoolResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_virtual_machine_scale_set_standby_pool" "import" {
  name                                  = azurerm_virtual_machine_scale_set_standby_pool.test.name
  resource_group_name                   = azurerm_resource_group.test.name
  location                              = "%s"
  attached_virtual_machine_scale_set_id = azurerm_virtual_machine_scale_set_standby_pool.test.attached_virtual_machine_scale_set_id
  virtual_machine_state                 = azurerm_virtual_machine_scale_set_standby_pool.test.virtual_machine_state

  elasticity_profile {
    max_ready_capacity = 10
    min_ready_capacity = 5
  }
}
`, config, data.Locations.Primary)
}

func (r StandbyPoolStandbyVirtualMachinePoolResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_virtual_machine_scale_set_standby_pool" "test" {
  name                                  = "acctest-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  location                              = "%s"
  attached_virtual_machine_scale_set_id = azurerm_orchestrated_virtual_machine_scale_set.test.id
  virtual_machine_state                 = "Running"

  elasticity_profile {
    max_ready_capacity = 10
    min_ready_capacity = 5
  }

  tags = {
    key = "value"
  }
}
`, template, data.RandomIntOfLength(16), data.Locations.Primary)
}

func (r StandbyPoolStandbyVirtualMachinePoolResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_orchestrated_virtual_machine_scale_set" "update" {
  name                = "acctestOVMSS-update-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_fault_domain_count = 1

  zones = ["1"]
}

resource "azurerm_virtual_machine_scale_set_standby_pool" "test" {
  name                                  = "acctest-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  location                              = "%s"
  attached_virtual_machine_scale_set_id = azurerm_orchestrated_virtual_machine_scale_set.update.id
  virtual_machine_state                 = "Deallocated"

  elasticity_profile {
    max_ready_capacity = 12
    min_ready_capacity = 6
  }

  tags = {
    key = "updatedValue"
  }
}
`, template, data.RandomInteger, data.RandomIntOfLength(16), data.Locations.Primary)
}

func (r StandbyPoolStandbyVirtualMachinePoolResource) minCapacityExceedMaxCapacity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_virtual_machine_scale_set_standby_pool" "test" {
  name                                  = "acctest-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  location                              = "%s"
  attached_virtual_machine_scale_set_id = azurerm_orchestrated_virtual_machine_scale_set.test.id
  virtual_machine_state                 = "Running"

  elasticity_profile {
    max_ready_capacity = 5
    min_ready_capacity = 10
  }
}
`, template, data.RandomIntOfLength(16), data.Locations.Primary)
}
