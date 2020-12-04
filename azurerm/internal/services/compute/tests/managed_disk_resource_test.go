package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMManagedDisk_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_empty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagedDisk_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_empty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
			{
				Config:      testAccAzureRMManagedDisk_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_managed_disk"),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_zeroGbFromPlatformImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_zeroGbFromPlatformImage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
				ExpectNonEmptyPlan: true, // since the `disk_size_gb` will have changed
			},
		},
	})
}

func TestAccAzureRMManagedDisk_import(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var vm compute.VirtualMachine
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then delete it so we can use the vhd to test import
				Config:             testAccAzureRMVirtualMachine_basicLinuxMachine(data),
				Destroy:            false,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists("azurerm_virtual_machine.test", &vm),
					testDeleteAzureRMVirtualMachine("azurerm_virtual_machine.test"),
				),
			},
			{
				Config: testAccAzureRMManagedDisk_import(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_copy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_copy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_fromPlatformImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_platformImage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_empty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "acctest"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost-center", "ops"),
					resource.TestCheckResourceAttr(data.ResourceName, "disk_size_gb", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_account_type", string(compute.StorageAccountTypesStandardLRS)),
				),
			},
			{
				Config: testAccAzureRMManagedDisk_empty_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "acctest"),
					resource.TestCheckResourceAttr(data.ResourceName, "disk_size_gb", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_account_type", string(compute.StorageAccountTypesPremiumLRS)),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_encryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_encryption(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
					resource.TestCheckResourceAttr(data.ResourceName, "encryption_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "encryption_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "encryption_settings.0.disk_encryption_key.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "encryption_settings.0.disk_encryption_key.0.secret_url"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "encryption_settings.0.disk_encryption_key.0.source_vault_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "encryption_settings.0.key_encryption_key.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "encryption_settings.0.key_encryption_key.0.key_url"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "encryption_settings.0.key_encryption_key.0.source_vault_id"),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_importEmpty_withZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_empty_withZone(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagedDisk_create_withUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_create_withUltraSSD(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagedDisk_update_withUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_create_withUltraSSD(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
					resource.TestCheckResourceAttr(data.ResourceName, "disk_iops_read_write", "101"),
					resource.TestCheckResourceAttr(data.ResourceName, "disk_mbps_read_write", "10"),
				),
			},
			{
				Config: testAccAzureRMManagedDisk_update_withUltraSSD(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
					resource.TestCheckResourceAttr(data.ResourceName, "disk_iops_read_write", "102"),
					resource.TestCheckResourceAttr(data.ResourceName, "disk_mbps_read_write", "11"),
				),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_import_withUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_create_withUltraSSD(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
			{
				Config:      testAccAzureRMManagedDisk_import_withUltraSSD(data),
				ExpectError: acceptance.RequiresImportError("azurerm_managed_disk"),
			},
		},
	})
}

func TestAccAzureRMManagedDisk_diskEncryptionSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_diskEncryptionSetEncrypted(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagedDisk_diskEncryptionSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_diskEncryptionSetUnencrypted(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMManagedDisk_diskEncryptionSetEncrypted(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagedDisk_attachedDiskUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_managedDiskAttached(data, 10),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMManagedDisk_managedDiskAttached(data, 20),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
					resource.TestCheckResourceAttr(data.ResourceName, "disk_size_gb", "20"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagedDisk_attachedStorageTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDisk_storageTypeUpdateWhilstAttached(data, "Standard_LRS"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMManagedDisk_storageTypeUpdateWhilstAttached(data, "Premium_LRS"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
			data.ImportStep(),
		},
	})
}

// nolint unparam
func testCheckAzureRMManagedDiskExists(resourceName string, d *compute.Disk, shouldExist bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DisksClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		dName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for disk: %s", dName)
		}

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
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DisksClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testDeleteAzureRMVirtualMachine(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		vmName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual machine: %s", vmName)
		}

		future, err := client.Delete(ctx, resourceGroup, vmName, utils.Bool(false))
		if err != nil {
			return fmt.Errorf("Bad: Delete on vmClient: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Bad: Delete on vmClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMManagedDisk_empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDisk_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMManagedDisk_empty(data)
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "import" {
  name                 = azurerm_managed_disk.test.name
  location             = azurerm_managed_disk.test.location
  resource_group_name  = azurerm_managed_disk.test.resource_group_name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, template)
}

func testAccAzureRMManagedDisk_empty_withZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"
  zones                = ["1"]

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDisk_import(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Import"
  source_uri           = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
  storage_account_id   = azurerm_storage_account.test.id
  disk_size_gb         = "45"

  tags = {
    environment = "acctest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMManagedDisk_copy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "source" {
  name                 = "acctestd1-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd2-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Copy"
  source_resource_id   = azurerm_managed_disk.source.id
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMManagedDisk_empty_updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = "2"

  tags = {
    environment = "acctest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDisk_platformImage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "UbuntuServer"
  sku       = "16.04-LTS"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  os_type              = "Linux"
  create_option        = "FromImage"
  image_reference_id   = data.azurerm_platform_image.test.id
  storage_account_type = "Standard_LRS"
}
`, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDisk_zeroGbFromPlatformImage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "UbuntuServer"
  sku       = "16.04-LTS"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  os_type              = "Linux"
  create_option        = "FromImage"
  disk_size_gb         = 0
  image_reference_id   = data.azurerm_platform_image.test.id
  storage_account_type = "Standard_LRS"
}
`, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDisk_encryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"
  sku_name            = "premium"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.object_id}"

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

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "EC"
  key_size     = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

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

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString, data.RandomInteger)
}

func testAccAzureRMManagedDisk_create_withUltraSSD(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "UltraSSD_LRS"
  create_option        = "Empty"
  disk_size_gb         = "4"
  disk_iops_read_write = "101"
  disk_mbps_read_write = "10"
  zones                = ["1"]

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDisk_update_withUltraSSD(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "UltraSSD_LRS"
  create_option        = "Empty"
  disk_size_gb         = "4"
  disk_iops_read_write = "102"
  disk_mbps_read_write = "11"
  zones                = ["1"]

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDisk_import_withUltraSSD(data acceptance.TestData) string {
	template := testAccAzureRMManagedDisk_create_withUltraSSD(data)
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "import" {
  name                 = azurerm_managed_disk.test.name
  location             = azurerm_managed_disk.test.location
  resource_group_name  = azurerm_managed_disk.test.resource_group_name
  storage_account_type = "UltraSSD_LRS"
  create_option        = "Empty"
  disk_size_gb         = "4"
  disk_iops_read_write = "101"
  disk_mbps_read_write = "10"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, template)
}

func testAccAzureRMManagedDisk_diskEncryptionSetDependencies(data acceptance.TestData) string {
	// whilst this is in Preview it's only supported in: West Central US, Canada Central, North Europe
	// TODO: switch back to default location
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkv%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "premium"
  enabled_for_disk_encryption = true
  soft_delete_enabled         = true
  purge_protection_enabled    = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "create",
    "delete",
    "get",
    "update",
  ]

  secret_permissions = [
    "get",
    "delete",
    "set",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = ["azurerm_key_vault_access_policy.service-principal"]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestdes-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "disk-encryption" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "get",
    "wrapkey",
    "unwrapkey",
  ]

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_role_assignment" "disk-encryption-read-keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_disk_encryption_set.test.identity.0.principal_id
}
`, data.RandomInteger, location, data.RandomString, data.RandomInteger)
}

func testAccAzureRMManagedDisk_diskEncryptionSetEncrypted(data acceptance.TestData) string {
	template := testAccAzureRMManagedDisk_diskEncryptionSetDependencies(data)
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "test" {
  name                   = "acctestd-%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  storage_account_type   = "Standard_LRS"
  create_option          = "Empty"
  disk_size_gb           = 1
  disk_encryption_set_id = azurerm_disk_encryption_set.test.id

  depends_on = [
    "azurerm_role_assignment.disk-encryption-read-keyvault",
    "azurerm_key_vault_access_policy.disk-encryption",
  ]
}
`, template, data.RandomInteger)
}

func testAccAzureRMManagedDisk_diskEncryptionSetUnencrypted(data acceptance.TestData) string {
	template := testAccAzureRMManagedDisk_diskEncryptionSetDependencies(data)

	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 1

  depends_on = [
    "azurerm_role_assignment.disk-encryption-read-keyvault",
    "azurerm_key_vault_access_policy.disk-encryption",
  ]
}
`, template, data.RandomInteger)
}

func testAccAzureRMManagedDisk_managedDiskAttached(data acceptance.TestData, diskSize int) string {
	template := testAccAzureRMManagedDisk_templateAttached(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_managed_disk" "test" {
  name                 = "%d-disk1"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = %d
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  lun                = "0"
  caching            = "None"
}
`, template, data.RandomInteger, diskSize)
}

func testAccAzureRMManagedDisk_storageTypeUpdateWhilstAttached(data acceptance.TestData, storageAccountType string) string {
	template := testAccAzureRMManagedDisk_templateAttached(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_managed_disk" "test" {
  name                 = "acctestdisk-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "%s"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  lun                = "0"
  caching            = "None"
}
`, template, data.RandomInteger, storageAccountType)
}

func testAccAzureRMManagedDisk_templateAttached(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestvm-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_D2s_v3"
  admin_username                  = "adminuser"
  admin_password                  = "Password1234!"
  disable_password_authentication = false

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
