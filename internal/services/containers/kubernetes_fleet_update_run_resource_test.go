// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/updateruns"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KubernetesFleetUpdateRunTestResource struct{}

func TestAccKubernetesFleetUpdateRun_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_update_run", "test")
	r := KubernetesFleetUpdateRunTestResource{}

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

func TestAccKubernetesFleetUpdateRun_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_update_run", "test")
	r := KubernetesFleetUpdateRunTestResource{}

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

func TestAccKubernetesFleetUpdateRun_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_update_run", "test")
	r := KubernetesFleetUpdateRunTestResource{}

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

func TestAccKubernetesFleetUpdateRun_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_update_run", "test")
	r := KubernetesFleetUpdateRunTestResource{}

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
func (r KubernetesFleetUpdateRunTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := updateruns.ParseUpdateRunID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.FleetUpdateRunsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}
func (r KubernetesFleetUpdateRunTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_update_run" "test" {
  name                        = "acctestfus-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  managed_cluster_update {
    upgrade {
      type = "NodeImageOnly"
    }
  }
  fleet_update_strategy_id = azurerm_kubernetes_fleet_update_strategy.test.id

  lifecycle {
    ignore_changes = [stage]
  }
  depends_on = [azurerm_kubernetes_fleet_member.test]
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetUpdateRunTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_update_run" "import" {
  name                        = azurerm_kubernetes_fleet_update_run.test.name
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_update_run.test.kubernetes_fleet_manager_id
  managed_cluster_update {
    upgrade {
      type = "NodeImageOnly"
    }
  }
  fleet_update_strategy_id = azurerm_kubernetes_fleet_update_strategy.test.id
}
`, r.basic(data))
}

func (r KubernetesFleetUpdateRunTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_update_run" "test" {
  name                        = "acctestfus-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  managed_cluster_update {
    upgrade {
      type               = "Full"
      kubernetes_version = "1.27"
    }
    node_image_selection {
      type = "Latest"
    }
  }
  stage {
    name = "acctestfus-%[2]d"
    group {
      name = "acctestfus-%[2]d"
    }
    after_stage_wait_in_seconds = 21
  }

  depends_on = [azurerm_kubernetes_fleet_member.test]
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetUpdateRunTestResource) template(data acceptance.TestData) string {
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

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestkc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestkc-%[2]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_fleet_member" "test" {
  name                  = "acctestkfm-%[2]d"
  kubernetes_fleet_id   = azurerm_kubernetes_fleet_manager.test.id
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  group                 = "acctestfus-%[2]d"
}

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
`, data.Locations.Primary, data.RandomInteger)
}
