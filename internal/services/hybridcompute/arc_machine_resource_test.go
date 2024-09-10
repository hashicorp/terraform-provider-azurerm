// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hybridcompute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ArcMachineResource struct{}

func TestAccArcMachineResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_machine", "test")
	r := ArcMachineResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccArcMachineResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_machine", "test")
	r := ArcMachineResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r ArcMachineResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := machines.ParseMachineID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HybridCompute.HybridComputeClient_v2024_07_10.Machines.Get(ctx, *id, machines.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ArcMachineResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_arc_machine" "test" {
  name                = "acctest-hcm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "SCVMM"
}
`, r.template(data), data.RandomInteger)
}

func (r ArcMachineResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_arc_machine" "import" {
  name                = azurerm_arc_machine.test.name
  resource_group_name = azurerm_arc_machine.test.resource_group_name
  location            = azurerm_arc_machine.test.location
  kind                = azurerm_arc_machine.test.kind
}
`, r.basic(data))
}

func (r ArcMachineResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-hcm-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
