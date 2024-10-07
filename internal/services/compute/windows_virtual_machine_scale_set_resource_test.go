// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WindowsVirtualMachineScaleSetResource struct{}

func (r WindowsVirtualMachineScaleSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.VirtualMachineScaleSetsClient.Get(ctx, *id, virtualmachinescalesets.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving Windows %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (WindowsVirtualMachineScaleSetResource) vmName(data acceptance.TestData) string {
	// windows VM names can be up to 15 chars, however the prefix can only be 9 chars
	return fmt.Sprintf("acctvm%s", fmt.Sprintf("%d", data.RandomInteger)[0:2])
}

func (r WindowsVirtualMachineScaleSetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  vm_name = "%s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}
`, r.vmName(data), data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
