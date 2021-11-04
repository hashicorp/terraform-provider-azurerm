package storage_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"testing"
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

func (s StorageDisksPoolResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.StorageDisksPoolID(state.ID)
	if err != nil {
		return nil, err
	}
	client := clients.Storage.DisksPoolsClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.DiskPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r StorageDisksPoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_disks_pool" "test" {
  availability_zones  = ["1"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctest-diskspool-%s"
  sku_name            = "Basic_B1"
  subnet_id           = azurerm_subnet.test.id
  tags                = {
    "env" = "qa"
  }
}
`, r.template(data), data.RandomString)
}

func (r StorageDisksPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  location = "%[1]s"
  name     = "acctestRG-diskspool-%[2]d"
}

resource "azurerm_virtual_network" "test" {
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestvnet%[2]d"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[2]d"
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
