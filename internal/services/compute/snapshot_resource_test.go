// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SnapshotResource struct{}

func TestAccSnapshot_fromManagedDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")
	r := SnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fromManagedDisk(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("source_uri"),
	})
}

func TestAccSnapshot_networkAccessPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")
	r := SnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkAccessPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccSnapshot_networkAccessPolicyAllowPrivate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")
	r := SnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkAccessPolicyAllowPrivate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccSnapshot_publicNetworkAccess(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")
	r := SnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccSnapshot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")
	r := SnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fromManagedDisk(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_snapshot"),
		},
	})
}

func TestAccSnapshot_encryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")
	r := SnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("source_uri"),
		{
			Config: r.encryptionUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("source_uri"),
	})
}

func TestAccSnapshot_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")
	r := SnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fromManagedDisk(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.fromManagedDiskUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccSnapshot_extendingManagedDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")
	r := SnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extendingManagedDisk(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccSnapshot_fromExistingSnapshot(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "second")
	r := SnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fromExistingSnapshot(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccSnapshot_fromUnmanagedDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")
	r := SnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fromUnmanagedDisk(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccSnapshot_incrementalEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")
	r := SnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.incrementalEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("source_uri"),
	})
}

func TestAccSnapshot_trustedLaunch(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshot", "test")
	r := SnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.trustedLaunch(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("trusted_launch_enabled").HasValue("true"),
			),
		},
		data.ImportStep("source_uri"),
	})
}

func (t SnapshotResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := snapshots.ParseSnapshotID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.SnapshotsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Shared Image Gallery %q", id)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (SnapshotResource) fromManagedDisk(data acceptance.TestData) string {
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
  name                          = "acctestss_%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  create_option                 = "Copy"
  source_uri                    = azurerm_managed_disk.test.id
  public_network_access_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r SnapshotResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_snapshot" "import" {
  name                = azurerm_snapshot.test.name
  location            = azurerm_snapshot.test.location
  resource_group_name = azurerm_snapshot.test.resource_group_name
  create_option       = azurerm_snapshot.test.create_option
  source_uri          = azurerm_snapshot.test.source_uri
}
`, r.fromManagedDisk(data))
}

func (SnapshotResource) fromManagedDiskUpdated(data acceptance.TestData) string {
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
  name                          = "acctestss_%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  create_option                 = "Copy"
  source_uri                    = azurerm_managed_disk.test.id
  network_access_policy         = "AllowAll"
  public_network_access_enabled = true

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (SnapshotResource) encryptionTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults       = false
      purge_soft_delete_on_destroy          = false
      purge_soft_deleted_keys_on_destroy    = false
      purge_soft_deleted_secrets_on_destroy = false
    }
  }
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
  name                     = "acctestkv%s"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  tenant_id                = "${data.azurerm_client_config.current.tenant_id}"
  purge_protection_enabled = true

  sku_name = "standard"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.object_id}"

    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Delete",
      "Get",
      "Set",
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r SnapshotResource) encryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_snapshot" "test" {
  name                = "acctestss_%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  create_option       = "Copy"
  source_uri          = "${azurerm_managed_disk.test.id}"
  disk_size_gb        = "20"

  encryption_settings {
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
`, r.encryptionTemplate(data), data.RandomInteger)
}

func (r SnapshotResource) encryptionUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault" "test2" {
  name                     = "acctestkv2%[2]s"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  tenant_id                = "${data.azurerm_client_config.current.tenant_id}"
  purge_protection_enabled = true

  sku_name = "standard"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.object_id}"

    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]
  }

  enabled_for_disk_encryption = true
}

resource "azurerm_key_vault_key" "test2" {
  name         = "generated-certificate"
  key_vault_id = azurerm_key_vault.test2.id
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

resource "azurerm_key_vault_secret" "test2" {
  name         = "secret-sauce"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.test2.id
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss_%[3]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  create_option       = "Copy"
  source_uri          = "${azurerm_managed_disk.test.id}"
  disk_size_gb        = "20"

  encryption_settings {
    disk_encryption_key {
      secret_url      = "${azurerm_key_vault_secret.test2.id}"
      source_vault_id = "${azurerm_key_vault.test2.id}"
    }

    key_encryption_key {
      key_url         = "${azurerm_key_vault_key.test2.id}"
      source_vault_id = "${azurerm_key_vault.test2.id}"
    }
  }
}
`, r.encryptionTemplate(data), data.RandomString, data.RandomInteger)
}

func (SnapshotResource) extendingManagedDisk(data acceptance.TestData) string {
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
  name                          = "acctestss_%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  create_option                 = "Copy"
  source_uri                    = azurerm_managed_disk.test.id
  disk_size_gb                  = "20"
  public_network_access_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (SnapshotResource) fromExistingSnapshot(data acceptance.TestData) string {
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

func (SnapshotResource) fromUnmanagedDisk(data acceptance.TestData) string {
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
  address_prefixes     = ["10.0.2.0/24"]
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
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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
  storage_account_id  = "${azurerm_storage_account.test.id}"
  depends_on = [
    azurerm_virtual_machine.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (SnapshotResource) incrementalEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestmd-%[2]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss_%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
  incremental_enabled = true
}
`, data.Locations.Primary, data.RandomInteger)
}

func (SnapshotResource) trustedLaunch(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_platform_image" "test" {
  location  = "%[2]s"
  publisher = "Canonical"
  offer     = "0001-com-ubuntu-server-jammy"
  sku       = "22_04-lts-gen2"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_managed_disk" "test" {
  name                   = "acctestd-%[1]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  os_type                = "Linux"
  create_option          = "FromImage"
  image_reference_id     = data.azurerm_platform_image.test.id
  storage_account_type   = "Standard_LRS"
  hyper_v_generation     = "V2"
  trusted_launch_enabled = true
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss_%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (SnapshotResource) networkAccessPolicy(data acceptance.TestData) string {
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
  name                  = "acctestss_%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  create_option         = "Copy"
  source_uri            = azurerm_managed_disk.test.id
  network_access_policy = "AllowAll"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (SnapshotResource) networkAccessPolicyAllowPrivate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_disk_access" "test" {
  name                = "acctestDA-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
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
  name                  = "acctestss_%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  disk_access_id        = azurerm_disk_access.test.id
  create_option         = "Copy"
  source_uri            = azurerm_managed_disk.test.id
  network_access_policy = "AllowPrivate"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (SnapshotResource) publicNetworkAccess(data acceptance.TestData) string {
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
  name                          = "acctestss_%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  create_option                 = "Copy"
  source_uri                    = azurerm_managed_disk.test.id
  public_network_access_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
