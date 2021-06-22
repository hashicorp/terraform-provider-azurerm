package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KubernetesClusterMaintenanceConfigurationResource struct{}

func TestAccKubernetesClusterMaintenanceConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_maintenance_configuration", "test")
	r := KubernetesClusterMaintenanceConfigurationResource{}
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

func TestAccKubernetesClusterMaintenanceConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_maintenance_configuration", "test")
	r := KubernetesClusterMaintenanceConfigurationResource{}
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

func TestAccKubernetesClusterMaintenanceConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_maintenance_configuration", "test")
	r := KubernetesClusterMaintenanceConfigurationResource{}
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

func TestAccKubernetesClusterMaintenanceConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_maintenance_configuration", "test")
	r := KubernetesClusterMaintenanceConfigurationResource{}
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
			Config: r.notAllowedTime(data),
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

func TestAccKubernetesClusterMaintenanceConfiguration_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_maintenance_configuration", "first")
	r := KubernetesClusterMaintenanceConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r KubernetesClusterMaintenanceConfigurationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.MaintenanceConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Containers.MaintenanceConfigurationsClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r KubernetesClusterMaintenanceConfigurationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_maintenance_configuration" "test" {
  name                  = "acctest-CMC-%d"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  maintenance_allowed {
    day        = "Monday"
    hour_slots = [1, 2]
  }
}
`, template, data.RandomInteger)
}

func (r KubernetesClusterMaintenanceConfigurationResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_maintenance_configuration" "import" {
  name                  = azurerm_kubernetes_cluster_maintenance_configuration.test.name
  kubernetes_cluster_id = azurerm_kubernetes_cluster_maintenance_configuration.test.kubernetes_cluster_id

  dynamic "maintenance_allowed" {
    for_each = azurerm_kubernetes_cluster_maintenance_configuration.test.maintenance_allowed
    content {
      day        = maintenance_allowed.value.day
      hour_slots = maintenance_allowed.value.hour_slots
    }
  }
}
`, config)
}

func (r KubernetesClusterMaintenanceConfigurationResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_maintenance_configuration" "test" {
  name                  = "acctest-cmc-%d"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id

  maintenance_not_allowed_window {
    end   = "2021-11-30T12:00:00Z"
    start = "2021-11-26T03:00:00Z"
  }

  maintenance_not_allowed_window {
    end   = "2021-12-30T12:00:00Z"
    start = "2021-12-26T03:00:00Z"
  }

  maintenance_allowed {
    day        = "Monday"
    hour_slots = [1, 2]
  }

  maintenance_allowed {
    day        = "Sunday"
    hour_slots = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23]
  }
}
`, template, data.RandomInteger)
}

func (r KubernetesClusterMaintenanceConfigurationResource) notAllowedTime(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_maintenance_configuration" "test" {
  name                  = "acctest-cmc-%d"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id

  maintenance_not_allowed_window {
    end   = "2021-11-30T12:00:00Z"
    start = "2021-11-26T03:00:00Z"
  }
}
`, template, data.RandomInteger)
}

func (r KubernetesClusterMaintenanceConfigurationResource) multiple(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_maintenance_configuration" "first" {
  name                  = "acctest-CMC1-%d"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  maintenance_allowed {
    day        = "Monday"
    hour_slots = [1, 2]
  }
}

resource "azurerm_kubernetes_cluster_maintenance_configuration" "second" {
  name                  = "acctest-CMC2-%d"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  maintenance_allowed {
    day        = "Sunday"
    hour_slots = [20, 21]
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r KubernetesClusterMaintenanceConfigurationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
