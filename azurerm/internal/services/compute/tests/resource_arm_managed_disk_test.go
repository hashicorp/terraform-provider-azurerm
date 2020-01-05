package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

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
				//need to create a vm and then delete it so we can use the vhd to test import
				Config:             testAccAzureRMVirtualMachine_basicLinuxMachine(data),
				Destroy:            false,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists(data.ResourceName, &vm),
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

func TestAccAzureRMManagedDisk_NonStandardCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	var d compute.Disk

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedDiskNonStandardCasing(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedDiskExists(data.ResourceName, &d, true),
				),
			},
			{
				Config:             testAccAzureRMManagedDiskNonStandardCasing(data),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

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

func testCheckAzureRMManagedDiskExists(resourceName string, d *compute.Disk, shouldExist bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		dName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for disk: %s", dName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DisksClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		vmName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual machine: %s", vmName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		future, err := client.Delete(ctx, resourceGroup, vmName)
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
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
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
  name                 = "${azurerm_managed_disk.test.name}"
  location             = "${azurerm_managed_disk.test.location}"
  resource_group_name  = "${azurerm_managed_disk.test.resource_group_name}"
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
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
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
  resource_group_name   = "${azurerm_resource_group.test.name}"
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
  disk_size_gb         = "45"

  tags = {
    environment = "acctest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMManagedDisk_copy(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "source" {
  name                 = "acctestd1-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
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
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Copy"
  source_resource_id   = "${azurerm_managed_disk.source.id}"
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
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = "2"

  tags = {
    environment = "acctest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDiskNonStandardCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "standard_lrs"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDisk_platformImage(data acceptance.TestData) string {
	return fmt.Sprintf(`
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
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  os_type              = "Linux"
  create_option        = "FromImage"
  image_reference_id   = "${data.azurerm_platform_image.test.id}"
  storage_account_type = "Standard_LRS"
}
`, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDisk_zeroGbFromPlatformImage(data acceptance.TestData) string {
	return fmt.Sprintf(`
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
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  os_type              = "Linux"
  create_option        = "FromImage"
  disk_size_gb         = "0"
  image_reference_id   = "${data.azurerm_platform_image.test.id}"
  storage_account_type = "Standard_LRS"
}
`, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDisk_encryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
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

  tags = {
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
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
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
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
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
  name                 = "${azurerm_managed_disk.test.name}"
  location             = "${azurerm_managed_disk.test.location}"
  resource_group_name  = "${azurerm_managed_disk.test.resource_group_name}"
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
