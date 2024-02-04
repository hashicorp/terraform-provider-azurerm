package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15/fleetmembers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesFleetMemberTestResource struct{}

func TestAccKubernetesFleetMember_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_member", "test")
	r := KubernetesFleetMemberTestResource{}

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

func TestAccKubernetesFleetMember_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_member", "test")
	r := KubernetesFleetMemberTestResource{}

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

func TestAccKubernetesFleetMember_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_member", "test")
	r := KubernetesFleetMemberTestResource{}

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

func TestAccKubernetesFleetMember_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_member", "test")
	r := KubernetesFleetMemberTestResource{}

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
func (r KubernetesFleetMemberTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fleetmembers.ParseMemberID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.FleetMembersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}
func (r KubernetesFleetMemberTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_member" "test" {
  name                        = "acctestfus-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  kubernetes_cluster_id       = azurerm_kubernetes_cluster.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetMemberTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_member" "import" {
  name                        = azurerm_kubernetes_fleet_member.test.name
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_member.test.kubernetes_fleet_manager_id
  kubernetes_cluster_id       = azurerm_kubernetes_fleet_member.test.kubernetes_cluster_id
}
`, r.basic(data))
}

func (r KubernetesFleetMemberTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_member" "test" {
  name                        = "acctestfus-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  kubernetes_cluster_id       = azurerm_kubernetes_cluster.test.id
  group                       = "acctestgroup-%[2]d"
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetMemberTestResource) template(data acceptance.TestData) string {
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
    vm_size    = "Standard_D2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}
