package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestSnapshotName_validation(t *testing.T) {
	str := acctest.RandString(80)
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "ab",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "cosmosDBAccount1",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 0,
		},
		{
			Value:    "hello+world",
			ErrCount: 1,
		},
		{
			Value:    str,
			ErrCount: 0,
		},
		{
			Value:    str + "a",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := compute.ValidateSnapshotName(tc.Value, "azurerm_snapshot")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Snapshot Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAccAzureRMSnapshot_fromManagedDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSnapshot_fromManagedDisk(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSnapshotExists(data.ResourceName),
				),
			},
			data.ImportStep("source_uri"),
		},
	})
}

func TestAccAzureRMSnapshot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSnapshot_fromManagedDisk(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSnapshotExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMSnapshot_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_snapshot"),
			},
		},
	})
}

func TestAccAzureRMSnapshot_encryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSnapshot_encryption(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSnapshotExists(data.ResourceName),
				),
			},
			data.ImportStep("source_uri"),
		},
	})
}

func TestAccAzureRMSnapshot_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSnapshot_fromManagedDisk(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSnapshotExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMSnapshot_fromManagedDiskUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSnapshotExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMSnapshot_extendingManagedDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSnapshot_extendingManagedDisk(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSnapshotExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMSnapshot_fromExistingSnapshot(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "second")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSnapshot_fromExistingSnapshot(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSnapshotExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMSnapshot_fromUnmanagedDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSnapshot_fromUnmanagedDisk(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSnapshotExists(data.ResourceName),
				),
			},
		},
	})
}

func testCheckAzureRMSnapshotDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.SnapshotsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_snapshot" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Snapshot still exists:\n%+v", resp)
	}

	return nil
}

func testCheckAzureRMSnapshotExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.SnapshotsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Snapshot: %q", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Topic %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on snapshotsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMSnapshot_fromManagedDisk(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestmd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss_%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMSnapshot_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_snapshot" "import" {
  name                = azurerm_snapshot.test.name
  location            = azurerm_snapshot.test.location
  resource_group_name = azurerm_snapshot.test.resource_group_name
  create_option       = azurerm_snapshot.test.create_option
  source_uri          = azurerm_snapshot.test.source_uri
}
`, testAccAzureRMSnapshot_fromManagedDisk(data))
}

func testAccAzureRMSnapshot_fromManagedDiskUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestmd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss_%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMSnapshot_encryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestmd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "premium"

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
}

resource "azurerm_key_vault_key" "test" {
  name         = "generated-certificate"
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
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-sauce"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss_%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  create_option       = "Copy"
  source_uri          = "${azurerm_managed_disk.test.id}"
  disk_size_gb        = "20"

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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger)
}

func testAccAzureRMSnapshot_extendingManagedDisk(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestmd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss_%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
  disk_size_gb        = "20"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMSnapshot_fromExistingSnapshot(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "original" {
  name                 = "acctestmd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
}

resource "azurerm_snapshot" "first" {
  name                = "acctestss1_%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.original.id
}

resource "azurerm_snapshot" "second" {
  name                = "acctestss2_%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_resource_id  = azurerm_snapshot.first.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMSnapshot_fromUnmanagedDisk(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine" "test" {
  name                          = "acctestvm-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  network_interface_ids         = [azurerm_network_interface.test.id]
  vm_size                       = "Standard_F2"
  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "acctestvm-%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "staging"
  }
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss_%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Import"
  source_uri          = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
  depends_on = [
    azurerm_virtual_machine,
    test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
