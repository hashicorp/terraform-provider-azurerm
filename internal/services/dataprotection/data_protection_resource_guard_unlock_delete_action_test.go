// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

// before_destroy, caller, and on_failure require Terraform 1.16.
// TODO:  terraform-plugin-testing release
var terraformVersion1_16_0 = version.Must(version.NewVersion("1.16.0"))

type DataProtectionResourceGuardUnlockDeleteAction struct{}

func TestAccDataProtectionResourceGuardUnlockDeleteAction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_resource_guard_unlock_delete", "test")
	a := DataProtectionResourceGuardUnlockDeleteAction{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(terraformVersion1_16_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.basic(data),
				Check:  nil, // TODO
			},
			{
				Config:  a.basic(data),
				Destroy: true,
				Check:   nil, // TODO
			},
		},
	})
}

func TestAccDataProtectionResourceGuardUnlockDeleteAction_withoutGuard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_resource_guard_unlock_delete", "test")
	a := DataProtectionResourceGuardUnlockDeleteAction{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(terraformVersion1_16_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.withoutGuard(data),
				Check:  nil, // TODO
			},
			{
				Config:  a.withoutGuard(data),
				Destroy: true,
				Check:   nil, // TODO
			},
		},
	})
}

func (a DataProtectionResourceGuardUnlockDeleteAction) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_resource_guard" "test" {
  name                = "acctest-dprg-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_recovery_services_vault_resource_guard_association" "test" {
  vault_id          = azurerm_recovery_services_vault.test.id
  resource_guard_id = azurerm_data_protection_resource_guard.test.id
}

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  source_vm_id        = azurerm_virtual_machine.test.id
  backup_policy_id    = azurerm_backup_policy_vm.test.id

  include_disk_luns = [0]

  lifecycle {
    action_trigger {
      events     = [before_destroy]
      actions    = [action.azurerm_data_protection_resource_guard_unlock_delete.test]
      on_failure = halt
    }
  }

  depends_on = [azurerm_recovery_services_vault_resource_guard_association.test]
}
`, a.template(data), data.RandomInteger)
}

func (a DataProtectionResourceGuardUnlockDeleteAction) withoutGuard(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  source_vm_id        = azurerm_virtual_machine.test.id
  backup_policy_id    = azurerm_backup_policy_vm.test.id

  include_disk_luns = [0]

  lifecycle {
    action_trigger {
      events     = [before_destroy]
      actions    = [action.azurerm_data_protection_resource_guard_unlock_delete.test]
      on_failure = halt
    }
  }
}
`, a.template(data))
}

func (DataProtectionResourceGuardUnlockDeleteAction) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
terraform {
  required_version = ">= %s"
}

provider "azurerm" {
  features {
    recovery_service {
      vm_backup_stop_protection_and_retain_data_on_destroy = false
      purge_protected_items_from_vault_on_destroy          = true
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet"
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest_subnet"
  virtual_network_name = azurerm_virtual_network.test.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["10.0.10.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctest_nic"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "acctestipconfig"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-ip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "acctestip%d"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctest-datadisk"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1023"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctestvm"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  vm_size               = "Standard_D2s_v3"
  network_interface_ids = [azurerm_network_interface.test.id]

  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "acctest-osdisk"
    managed_disk_type = "Standard_LRS"
    caching           = "ReadWrite"
    create_option     = "FromImage"
  }

  storage_data_disk {
    name              = "acctest-datadisk"
    managed_disk_id   = azurerm_managed_disk.test.id
    managed_disk_type = "Standard_LRS"
    disk_size_gb      = azurerm_managed_disk.test.disk_size_gb
    create_option     = "Attach"
    lun               = 0
  }

  storage_data_disk {
    name              = "acctest-another-datadisk"
    create_option     = "Empty"
    disk_size_gb      = "1"
    lun               = 1
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "acctest"
    admin_username = "vmadmin"
    admin_password = "Password123!@#"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  boot_diagnostics {
    enabled     = true
    storage_uri = azurerm_storage_account.test.primary_blob_endpoint
  }

}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}

action "azurerm_data_protection_resource_guard_unlock_delete" "test" {
  config {
    protected_item_id = caller.id
  }
}
`, terraformVersion1_16_0.String(), data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}
