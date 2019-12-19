package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetAppSnapshot_basic(t *testing.T) {
	resourceName := "azurerm_netapp_snapshot.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppSnapshot_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMNetAppSnapshot_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_netapp_snapshot.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppSnapshot_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMNetAppSnapshot_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_netapp_snapshot"),
			},
		},
	})
}

func TestAccAzureRMNetAppSnapshot_complete(t *testing.T) {
	resourceName := "azurerm_netapp_snapshot.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppSnapshot_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMNetAppSnapshot_update(t *testing.T) {
	resourceName := "azurerm_netapp_snapshot.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	oldVolumeName := fmt.Sprintf("acctest-NetAppVolume-%d", ri)
	newVolumeName := fmt.Sprintf("acctest-updated-NetAppVolume-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppSnapshot_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "volume_name", oldVolumeName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMNetAppSnapshot_update(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "volume_name", newVolumeName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMNetAppSnapshotExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NetApp Snapshot not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		poolName := rs.Primary.Attributes["pool_name"]
		volumeName := rs.Primary.Attributes["volume_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).NetApp.SnapshotClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	client := testAccProvider.Meta().(*ArmClient).NetApp.SnapshotClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMNetAppSnapshot_basic(rInt int, location string) string {
	template := testAccAzureRMNetAppSnapshot_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot" "test" {
  name                = "acctest-NetAppSnapshot-%d"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  volume_name         = "${azurerm_netapp_volume.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template, rInt)
}

func testAccAzureRMNetAppSnapshot_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot" "import" {
  name                = "${azurerm_netapp_snapshot.test.name}"
  location            = "${azurerm_netapp_snapshot.test.location}"
  resource_group_name = "${azurerm_netapp_snapshot.test.name}"
}
`, testAccAzureRMNetAppSnapshot_basic(rInt, location))
}

func testAccAzureRMNetAppSnapshot_complete(rInt int, location string) string {
	template := testAccAzureRMNetAppSnapshot_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot" "test" {
  name                = "acctest-NetAppSnapshot-%d"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  volume_name         = "${azurerm_netapp_volume.test.name}"
  file_system_id      = "${azurerm_netapp_volume.test.file_system_id}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template, rInt)
}

func testAccAzureRMNetAppSnapshot_update(rInt int, location string) string {
	template := testAccAzureRMNetAppSnapshot_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "update" {
  name                = "acctest-updated-VirtualNetwork-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "update" {
  name                 = "acctest-updated-Subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.update.name}"
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  volume_path         = "my-updated-unique-file-path-%d"
  service_level       = "Premium"
  subnet_id           = "${azurerm_subnet.update.id}"
  storage_quota_in_gb = 100
}

resource "azurerm_netapp_snapshot" "test" {
  name                = "acctest-NetAppSnapshot-%d"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  volume_name         = "${azurerm_netapp_volume.update.name}"
  file_system_id      = "${azurerm_netapp_volume.update.file_system_id}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMNetAppSnapshot_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-Subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_netapp_account.test.name}"
  service_level       = "Premium"
  size_in_tb          = 4
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  volume_path         = "my-unique-file-path-%d"
  service_level       = "Premium"
  subnet_id           = "${azurerm_subnet.test.id}"
  storage_quota_in_gb = 100
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}
