package netapp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NetAppSnapshotResource struct {
}

func TestAccNetAppSnapshot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot", "test")
	r := NetAppSnapshotResource{}

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

func TestAccNetAppSnapshot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot", "test")
	r := NetAppSnapshotResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_netapp_snapshot"),
		},
	})
}

func TestAccNetAppSnapshot_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot", "test")
	r := NetAppSnapshotResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.FoO").HasValue("BaR"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppSnapshot_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot", "test")
	r := NetAppSnapshotResource{}
	oldVolumeName := fmt.Sprintf("acctest-NetAppVolume-%d", data.RandomInteger)
	newVolumeName := fmt.Sprintf("acctest-updated-NetAppVolume-%d", data.RandomInteger)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("volume_name").HasValue(oldVolumeName),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.FoO").HasValue("BaR"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("volume_name").HasValue(oldVolumeName),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.FoO").HasValue("BaZ"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("volume_name").HasValue(newVolumeName),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func (t NetAppSnapshotResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SnapshotID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.SnapshotClient.Get(ctx, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Netapp Snapshot (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r NetAppSnapshotResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot" "test" {
  name                = "acctest-NetAppSnapshot-%d"
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_name         = azurerm_netapp_volume.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r NetAppSnapshotResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot" "import" {
  name                = azurerm_netapp_snapshot.test.name
  location            = azurerm_netapp_snapshot.test.location
  resource_group_name = azurerm_netapp_snapshot.test.resource_group_name
  account_name        = azurerm_netapp_snapshot.test.account_name
  pool_name           = azurerm_netapp_snapshot.test.pool_name
  volume_name         = azurerm_netapp_snapshot.test.volume_name
}
`, r.basic(data))
}

func (r NetAppSnapshotResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot" "test" {
  name                = "acctest-NetAppSnapshot-%d"
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_name         = azurerm_netapp_volume.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "FoO" = "BaR"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetAppSnapshotResource) updateTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot" "test" {
  name                = "acctest-NetAppSnapshot-%d"
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_name         = azurerm_netapp_volume.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "FoO" = "BaZ"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetAppSnapshotResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "update" {
  name                = "acctest-updated-VirtualNetwork-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "update" {
  name                 = "acctest-updated-Subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.update.name
  address_prefix       = "10.0.2.0/24"

  delegation {
    name = "netapp"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_volume" "update" {
  name                = "acctest-updated-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-updated-unique-file-path-%d"
  service_level       = "Premium"
  subnet_id           = azurerm_subnet.update.id
  storage_quota_in_gb = 100
}

resource "azurerm_netapp_snapshot" "test" {
  name                = "acctest-NetAppSnapshot-%d"
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_name         = azurerm_netapp_volume.update.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (NetAppSnapshotResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-Subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"

  delegation {
    name = "netapp"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Premium"
  size_in_tb          = 4
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%d"
  service_level       = "Premium"
  subnet_id           = azurerm_subnet.test.id
  storage_quota_in_gb = 100
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
