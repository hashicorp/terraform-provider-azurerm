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

type StackHCIClusterResource struct{}

func TestAccStackHCICluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterResource{}

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

func TestAccStackHCICluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterResource{}

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

func TestAccStackHCICluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterResource{}

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

func TestAccStackHCICluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterResource{}

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

func (r StackHCIClusterResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	clusterClient := client.AzureStackHCI.ClusterClient
	id, err := parse.ClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving Azure Stack HCI Cluster %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.ClusterProperties != nil), nil
}

func (r StackHCIClusterResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_cluster" "test" {
  name                = "acctest-StackHCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = azuread_application.test.application_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
}
`, template, data.RandomInteger)
}

func (r StackHCIClusterResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_cluster" "import" {
  name                = azurerm_stack_hci_cluster.test.name
  resource_group_name = azurerm_stack_hci_cluster.test.resource_group_name
  location            = azurerm_stack_hci_cluster.test.location
  client_id           = azurerm_stack_hci_cluster.test.client_id
  tenant_id           = azurerm_stack_hci_cluster.test.tenant_id
}
`, config)
}

func (r StackHCIClusterResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_cluster" "test" {
  name                = "acctest-StackHCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = azuread_application.test.application_id
  tenant_id           = data.azurerm_client_config.current.tenant_id

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r StackHCIClusterResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_cluster" "test" {
  name                = "acctest-StackHCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = azuread_application.test.application_id
  tenant_id           = data.azurerm_client_config.current.tenant_id

  tags = {
    ENv = "Test2"
  }
}
`, template, data.RandomInteger)
}

func (r StackHCIClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "current" {}

resource "azuread_application" "test" {
  name = "acctestspa-%d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hci-%d"
  location = "%s"
}
`, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}
