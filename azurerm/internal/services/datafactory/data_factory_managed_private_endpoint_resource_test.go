package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ManagedPrivateEndpointResource struct{}

func TestAccDataFactoryManagedPrivateEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_managed_private_endpoint", "test")
	r := ManagedPrivateEndpointResource{}

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

func TestAccDataFactoryManagedPrivateEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_managed_private_endpoint", "test")
	r := ManagedPrivateEndpointResource{}

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

func (r ManagedPrivateEndpointResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ManagedPrivateEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	iter, err := client.DataFactory.ManagedPrivateEndpointsClient.ListByFactoryComplete(ctx, id.ResourceGroup, id.FactoryName, id.ManagedVirtualNetworkName)
	if err != nil {
		return nil, fmt.Errorf("listing %s: %+v", id, err)
	}
	for iter.NotDone() {
		managedPrivateEndpoint := iter.Value()
		if managedPrivateEndpoint.Name != nil && *managedPrivateEndpoint.Name == id.Name {
			return utils.Bool(true), nil
		}

		if err := iter.NextWithContext(ctx); err != nil {
			return nil, err
		}
	}
	return utils.Bool(false), nil
}

func (r ManagedPrivateEndpointResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_data_factory_managed_private_endpoint" "test" {
  name               = "acctestEndpoint%d"
  data_factory_id    = azurerm_data_factory.test.id
  target_resource_id = azurerm_storage_account.test.id
  subresource_name   = "blob"
}
`, template, data.RandomInteger)
}

func (r ManagedPrivateEndpointResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_data_factory_managed_private_endpoint" "import" {
  name               = azurerm_data_factory_managed_private_endpoint.test.name
  data_factory_id    = azurerm_data_factory_managed_private_endpoint.test.data_factory_id
  target_resource_id = azurerm_data_factory_managed_private_endpoint.test.target_resource_id
  subresource_name   = azurerm_data_factory_managed_private_endpoint.test.subresource_name
}
`, config)
}

func (r ManagedPrivateEndpointResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-adf-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                            = "acctestdf%d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  managed_virtual_network_enabled = true
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
