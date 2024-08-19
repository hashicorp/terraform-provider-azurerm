// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleetupdatestrategies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesFleetUpdateStrategyTestResource struct{}

func TestAccKubernetesFleetUpdateStrategy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_update_strategy", "test")
	r := KubernetesFleetUpdateStrategyTestResource{}

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

func TestAccKubernetesFleetUpdateStrategy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_update_strategy", "test")
	r := KubernetesFleetUpdateStrategyTestResource{}

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

func TestAccKubernetesFleetUpdateStrategy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_update_strategy", "test")
	r := KubernetesFleetUpdateStrategyTestResource{}

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

func TestAccKubernetesFleetUpdateStrategy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_update_strategy", "test")
	r := KubernetesFleetUpdateStrategyTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}
func (r KubernetesFleetUpdateStrategyTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fleetupdatestrategies.ParseUpdateStrategyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.FleetUpdateStrategiesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}
func (r KubernetesFleetUpdateStrategyTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_update_strategy" "test" {
  name                        = "acctestfus-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  stage {
    name = "acctestfus-%[2]d"
    group {
      name = "acctestfus-%[2]d"
    }
  }

}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetUpdateStrategyTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_update_strategy" "import" {
  name                        = azurerm_kubernetes_fleet_update_strategy.test.name
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_update_strategy.test.kubernetes_fleet_manager_id
  stage {
    name = "acctestfus-%[2]d"
    group {
      name = "acctestfus-%[2]d"
    }
  }
}
`, r.basic(data), data.RandomInteger)
}

func (r KubernetesFleetUpdateStrategyTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_update_strategy" "test" {
  name                        = "acctestfus-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  stage {
    name = "acctestfus-%[2]d-complte"
    group {
      name = "acctestfus-%[2]d-complete"
    }
    after_stage_wait_in_seconds = 21
  }
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetUpdateStrategyTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_fleet_manager" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestkfm-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  hub_profile {
    dns_prefix = "val-%[2]d"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}
