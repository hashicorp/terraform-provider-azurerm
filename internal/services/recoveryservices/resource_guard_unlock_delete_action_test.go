// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

type ResourceGuardUnlockDeleteAction struct{}

func TestAccResourceGuardUnlockDeleteAction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_guard_unlock_delete", "test")
	a := ResourceGuardUnlockDeleteAction{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.16.0"))),
		},
		Steps: []resource.TestStep{
			{
				Config: a.basic(data),
				Check:  check.That("azurerm_backup_protected_vm.test").ExistsInAzure(BackupProtectedVmResource{}),
			},
			{
				Config:  a.basic(data),
				Destroy: true,
			},
		},
	})
}

func TestAccResourceGuardUnlockDeleteAction_withoutGuard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_guard_unlock_delete", "test")
	a := ResourceGuardUnlockDeleteAction{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.16.0"))),
		},
		Steps: []resource.TestStep{
			{
				Config: a.withoutGuard(data),
				Check:  check.That("azurerm_backup_protected_vm.test").ExistsInAzure(BackupProtectedVmResource{}),
			},
			{
				Config:  a.withoutGuard(data),
				Destroy: true,
			},
		},
	})
}

func (a ResourceGuardUnlockDeleteAction) basic(data acceptance.TestData) string {
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
      actions    = [action.azurerm_resource_guard_unlock_delete.test]
      on_failure = halt
    }
  }

  depends_on = [azurerm_recovery_services_vault_resource_guard_association.test]
}
`, a.template(data), data.RandomInteger)
}

func (a ResourceGuardUnlockDeleteAction) withoutGuard(data acceptance.TestData) string {
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
      actions    = [action.azurerm_resource_guard_unlock_delete.test]
      on_failure = halt
    }
  }
}
`, a.template(data))
}

func (ResourceGuardUnlockDeleteAction) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
terraform {
  required_version = ">= 1.16.0"
}

provider "azurerm" {
  features {
    recovery_service {
      vm_backup_stop_protection_and_retain_data_on_destroy = false
      purge_protected_items_from_vault_on_destroy          = true
    }
  }
}

%s

action "azurerm_resource_guard_unlock_delete" "test" {
  config {
    protected_item_id = caller.id
  }
}
`, BackupProtectedVmResource{}.baseWithOutProvider(data))
}
