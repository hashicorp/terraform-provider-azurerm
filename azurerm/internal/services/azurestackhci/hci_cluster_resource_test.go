package azurestackhci_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/azurestackhci/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type HCIClusterResource struct{}

func TestAccHCICluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hci_cluster", "test")
	r := HCIClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHCICluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hci_cluster", "test")
	r := HCIClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccHCICluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hci_cluster", "test")
	r := HCIClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHCICluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hci_cluster", "test")
	r := HCIClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r HCIClusterResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	clusterClient := client.AzureStackHCI.ClusterClient
	id, err := parse.HciClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, id.ResourceGroup, id.ClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving HCI Cluster %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.ClusterProperties != nil), nil
}

func (r HCIClusterResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hci_cluster" "test" {
  name                = "acctest-HCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = data.azurerm_client_config.current.client_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
}
`, template, data.RandomInteger)
}

func (r HCIClusterResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hci_cluster" "import" {
  name                = azurerm_hci_cluster.test.name
  resource_group_name = azurerm_hci_cluster.test.resource_group_name
  location            = azurerm_hci_cluster.test.location
  client_id           = azurerm_hci_cluster.test.client_id
  tenant_id           = azurerm_hci_cluster.test.tenant_id
}
`, config)
}

func (r HCIClusterResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hci_cluster" "test" {
  name                = "acctest-HCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = data.azurerm_client_config.current.client_id
  tenant_id           = data.azurerm_client_config.current.tenant_id

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r HCIClusterResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hci_cluster" "test" {
  name                = "acctest-HCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = data.azurerm_client_config.current.client_id
  tenant_id           = data.azurerm_client_config.current.tenant_id

  tags = {
    ENv = "Test2"
  }
}
`, template, data.RandomInteger)
}

func (r HCIClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hci-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
