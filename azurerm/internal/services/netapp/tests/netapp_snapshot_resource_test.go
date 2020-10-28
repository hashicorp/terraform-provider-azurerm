package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetAppSnapshot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppSnapshot_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetAppSnapshot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppSnapshot_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetAppSnapshot_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_netapp_snapshot"),
			},
		},
	})
}

func TestAccAzureRMNetAppSnapshot_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppSnapshot_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.FoO", "BaR"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetAppSnapshot_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot", "test")
	oldVolumeName := fmt.Sprintf("acctest-NetAppVolume-%d", data.RandomInteger)
	newVolumeName := fmt.Sprintf("acctest-updated-NetAppVolume-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppSnapshot_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "volume_name", oldVolumeName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.FoO", "BaR"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetAppSnapshot_updateTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "volume_name", oldVolumeName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.FoO", "BaZ"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetAppSnapshot_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "volume_name", newVolumeName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMNetAppSnapshotExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).NetApp.SnapshotClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NetApp Snapshot not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		poolName := rs.Primary.Attributes["pool_name"]
		volumeName := rs.Primary.Attributes["volume_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, accountName, poolName, volumeName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: NetApp Snapshot %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on netapp.SnapshotClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetAppSnapshotDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).NetApp.SnapshotClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_netapp_snapshot" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		poolName := rs.Primary.Attributes["pool_name"]
		volumeName := rs.Primary.Attributes["volume_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, accountName, poolName, volumeName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on netapp.SnapshotClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMNetAppSnapshot_basic(data acceptance.TestData) string {
	template := testAccAzureRMNetAppSnapshot_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMNetAppSnapshot_requiresImport(data acceptance.TestData) string {
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
`, testAccAzureRMNetAppSnapshot_basic(data))
}

func testAccAzureRMNetAppSnapshot_complete(data acceptance.TestData) string {
	template := testAccAzureRMNetAppSnapshot_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMNetAppSnapshot_updateTags(data acceptance.TestData) string {
	template := testAccAzureRMNetAppSnapshot_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMNetAppSnapshot_update(data acceptance.TestData) string {
	template := testAccAzureRMNetAppSnapshot_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "update" {
  name                = "acctest-updated-VirtualNetwork-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "update" {
  name               = "acctest-updated-Subnet-%d"
  virtual_network_id = azurerm_virtual_network.update.id
  address_prefix     = "10.0.2.0/24"

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
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetAppSnapshot_template(data acceptance.TestData) string {
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
  name               = "acctest-Subnet-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.2.0/24"

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
