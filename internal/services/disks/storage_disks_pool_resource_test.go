package disks_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/sdk/2021-08-01/diskpools"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageDisksPoolResource struct{}

func TestAccStorageDisksPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_disks_pool", "test")
	r := StorageDisksPoolResource{}
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

func TestAccStorageDisksPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_disks_pool", "test")
	r := StorageDisksPoolResource{}
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

func TestAccStorageDisksPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_disks_pool", "test")
	r := StorageDisksPoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
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
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageDisksPoolResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := diskpools.ParseDiskPoolID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.Disks.DiskPoolsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r StorageDisksPoolResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_disks_pool" "test" {
  name                = "acctest-diskspool-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  availability_zones  = ["1"]
  sku_name            = "Basic_B1"
  subnet_id           = azurerm_subnet.test.id
}
`, template, data.RandomString)
}

func (r StorageDisksPoolResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_disks_pool" "import" {
  name                = azurerm_storage_disks_pool.test.name
  resource_group_name = azurerm_storage_disks_pool.test.resource_group_name
  location            = azurerm_storage_disks_pool.test.location
  availability_zones  = azurerm_storage_disks_pool.test.availability_zones
  sku_name            = azurerm_storage_disks_pool.test.sku_name
  subnet_id           = azurerm_storage_disks_pool.test.subnet_id
}
`, template)
}

func (r StorageDisksPoolResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_disks_pool" "test" {
  name                = "acctest-diskspool-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  availability_zones  = ["1"]
  sku_name            = "Basic_B1"
  subnet_id           = azurerm_subnet.test.id
}
`, template, data.RandomString)
}

func (r StorageDisksPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-diskspool-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/24"]
  delegation {
    name = "diskspool"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/read"]
      name    = "Microsoft.StoragePool/diskPools"
    }
  }
}
`, data.Locations.Primary, data.RandomInteger)
}
