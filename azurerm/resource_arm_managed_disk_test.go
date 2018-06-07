package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMManagedDisk_empty(t *testing.T) {
	var d compute.Disk
	ri := acctest.RandInt()
	config := testAccAzureRMManagedDisk_empty(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists("azurerm_managed_disk.test", &d, true),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_zeroGbFromPlatformImage(t *testing.T) {
	var d compute.Disk
	ri := acctest.RandInt()
	config := testAccAzureRMManagedDisk_zeroGbFromPlatformImage(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists("azurerm_managed_disk.test", &d, true),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_import(t *testing.T) {
	var d compute.Disk
	var vm compute.VirtualMachine
	ri := acctest.RandInt()
	location := testLocation()
	vmConfig := testAccAzureRMVirtualMachine_basicLinuxMachine(ri, location)
	config := testAccAzureRMManagedDisk_import(ri, location)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				//need to create a vm and then delete it so we can use the vhd to test import
				Config:             vmConfig,
				Destroy:            false,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists("azurerm_virtual_machine.test", &vm),
					testDeleteAzureRMVirtualMachine("azurerm_virtual_machine.test"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists("azurerm_managed_disk.test", &d, true),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_copy(t *testing.T) {
	var d compute.Disk
	ri := acctest.RandInt()
	config := testAccAzureRMManagedDisk_copy(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists("azurerm_managed_disk.test", &d, true),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_fromPlatformImage(t *testing.T) {
	var d compute.Disk
	ri := acctest.RandInt()
	config := testAccAzureRMManagedDisk_platformImage(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists("azurerm_managed_disk.test", &d, true),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_update(t *testing.T) {
	var d compute.Disk

	resourceName := "azurerm_managed_disk.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMManagedDisk_empty(ri, testLocation())
	postConfig := testAccAzureRMManagedDisk_empty_updated(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(resourceName, &d, true),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "acctest"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost-center", "ops"),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_account_type", string(compute.StandardLRS)),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(resourceName, &d, true),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "acctest"),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "2"),
					resource.TestCheckResourceAttr(resourceName, "storage_account_type", string(compute.PremiumLRS)),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_encryption(t *testing.T) {
	var d compute.Disk

	resourceName := "azurerm_managed_disk.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMManagedDisk_encryption(ri, rs, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(resourceName, &d, true),
					resource.TestCheckResourceAttr(resourceName, "encryption_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "encryption_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "encryption_settings.0.disk_encryption_key.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "encryption_settings.0.disk_encryption_key.0.secret_url"),
					resource.TestCheckResourceAttrSet(resourceName, "encryption_settings.0.disk_encryption_key.0.source_vault_id"),
					resource.TestCheckResourceAttr(resourceName, "encryption_settings.0.key_encryption_key.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "encryption_settings.0.key_encryption_key.0.key_url"),
					resource.TestCheckResourceAttrSet(resourceName, "encryption_settings.0.key_encryption_key.0.source_vault_id"),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_NonStandardCasing(t *testing.T) {
	var d compute.Disk
	ri := acctest.RandInt()
	config := testAccAzureRMManagedDiskNonStandardCasing(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists("azurerm_managed_disk.test", &d, true),
				),
			},
			{
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testCheckAzureRMManagedDiskExists(name string, d *compute.Disk, shouldExist bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		dName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for disk: %s", dName)
		}

		client := testAccProvider.Meta().(*ArmClient).diskClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, dName)
		if err != nil {
			return fmt.Errorf("Bad: Get on diskClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound && shouldExist {
			return fmt.Errorf("Bad: ManagedDisk %q (resource group %q) does not exist", dName, resourceGroup)
		}
		if resp.StatusCode != http.StatusNotFound && !shouldExist {
			return fmt.Errorf("Bad: ManagedDisk %q (resource group %q) still exists", dName, resourceGroup)
		}

		*d = resp

		return nil
	}
}

func testCheckAzureRMManagedDiskDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).diskClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_managed_disk" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Managed Disk still exists: \n%#v", resp.DiskProperties)
		}
	}

	return nil
}

func testDeleteAzureRMVirtualMachine(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		vmName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual machine: %s", vmName)
		}

		client := testAccProvider.Meta().(*ArmClient).vmClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		future, err := client.Delete(ctx, resourceGroup, vmName)
		if err != nil {
			return fmt.Errorf("Bad: Delete on vmClient: %+v", err)
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Bad: Delete on vmClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMManagedDisk_empty(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_managed_disk" "test" {
    name = "acctestd-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_type = "Standard_LRS"
    create_option = "Empty"
    disk_size_gb = "1"

    tags {
        environment = "acctest"
        cost-center = "ops"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMManagedDisk_empty_withZone(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_managed_disk" "test" {
    name = "acctestd-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_type = "Standard_LRS"
    create_option = "Empty"
    disk_size_gb = "1"
    zones = ["1"]

    tags {
        environment = "acctest"
        cost-center = "ops"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMManagedDisk_import(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_storage_account" "test" {
    name                     = "accsa%d"
    resource_group_name      = "${azurerm_resource_group.test.name}"
    location                 = "${azurerm_resource_group.test.location}"
    account_tier             = "Standard"
    account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_managed_disk" "test" {
    name = "acctestd-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_type = "Standard_LRS"
    create_option = "Import"
    source_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    disk_size_gb = "45"

    tags {
        environment = "acctest"
    }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMManagedDisk_copy(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_managed_disk" "source" {
    name = "acctestd1-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_type = "Standard_LRS"
    create_option = "Empty"
    disk_size_gb = "1"

    tags {
        environment = "acctest"
        cost-center = "ops"
    }
}

resource "azurerm_managed_disk" "test" {
    name = "acctestd2-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_type = "Standard_LRS"
    create_option = "Copy"
    source_resource_id = "${azurerm_managed_disk.source.id}"
    disk_size_gb = "1"

    tags {
        environment = "acctest"
        cost-center = "ops"
    }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMManagedDisk_empty_updated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_managed_disk" "test" {
    name = "acctestd-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_type = "Premium_LRS"
    create_option = "Empty"
    disk_size_gb = "2"

    tags {
        environment = "acctest"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMManagedDiskNonStandardCasing(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}
resource "azurerm_managed_disk" "test" {
    name = "acctestd-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_type = "standard_lrs"
    create_option = "Empty"
    disk_size_gb = "1"
    tags {
        environment = "acctest"
        cost-center = "ops"
    }
}`, rInt, location, rInt)
}

func testAccAzureRMManagedDisk_platformImage(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "UbuntuServer"
  sku       = "16.04-LTS"
}

resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_managed_disk" "test" {
  name = "acctestd-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  os_type = "Linux"
  create_option = "FromImage"
  image_reference_id = "${data.azurerm_platform_image.test.id}"
  storage_account_type = "Standard_LRS"
}
`, location, rInt, location, rInt)
}

func testAccAzureRMManagedDisk_zeroGbFromPlatformImage(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "UbuntuServer"
  sku       = "16.04-LTS"
}

resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_managed_disk" "test" {
  name = "acctestd-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  os_type = "Linux"
  create_option = "FromImage"
  disk_size_gb = "0" 
  image_reference_id = "${data.azurerm_platform_image.test.id}"
  storage_account_type = "Standard_LRS"
}
`, location, rInt, location, rInt)
}

func testAccAzureRMManagedDisk_encryption(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "create",
      "delete",
      "get",
    ]

    secret_permissions = [
      "delete",
      "get",
      "set",
    ]
  }

  enabled_for_disk_encryption = true

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name      = "secret-%s"
  value     = "szechuan"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
}

resource "azurerm_key_vault_key" "test" {
  name      = "key-%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
  key_type  = "EC"
  key_size  = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}

resource "azurerm_managed_disk" "test" {
    name = "acctestd-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_type = "Standard_LRS"
    create_option = "Empty"
    disk_size_gb = "1"

    encryption_settings {
      enabled = true
      disk_encryption_key {
	    secret_url      = "${azurerm_key_vault_secret.test.id}"
	    source_vault_id = "${azurerm_key_vault.test.id}"
      }
      key_encryption_key {
        key_url         = "${azurerm_key_vault_key.test.id}"
        source_vault_id = "${azurerm_key_vault.test.id}"
      }
    }

    tags {
        environment = "acctest"
        cost-center = "ops"
    }
}
`, rInt, location, rString, rString, rString, rInt)
}
