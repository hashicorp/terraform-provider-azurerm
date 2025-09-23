// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2025-02-01/protectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BackupContainerVMAppTestResource struct{}

func TestAccBackupContainerVMAppSequential(t *testing.T) {
	// The dependent SAP server VM requires many complicated manual configurations. So it has to test based on the resource provided by service team.
	if os.Getenv("ARM_TEST_SAP_VIRTUAL_INSTANCE_NAME") == "" {
		t.Skip("Skipping as `ARM_TEST_SAP_VIRTUAL_INSTANCE_NAME` is not specified")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"sapVirtualInstance": {
			"basic":          testAccBackupContainerVMApp_basic,
			"requiresImport": testAccBackupContainerVMApp_requiresImport,
		},
	})
}

func testAccBackupContainerVMApp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_container_vm_app", "test")
	r := BackupContainerVMAppTestResource{}

	sourceVMID := os.Getenv("ARM_TEST_SAP_VM_ID")
	if sourceVMID == "" {
		t.Skip("Skipping as ARM_TEST_SAP_VM_ID is not set")
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, sourceVMID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("workload_type").HasValue("SAPAseDatabase"),
				check.That(data.ResourceName).Key("source_resource_id").HasValue(sourceVMID),
			),
		},
		data.ImportStep(),
	})
}

func testAccBackupContainerVMApp_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_container_vm_app", "test")
	r := BackupContainerVMAppTestResource{}

	sourceVMID := os.Getenv("ARM_TEST_SAP_VM_ID")
	if sourceVMID == "" {
		t.Skip("Skipping as ARM_TEST_SAP_VM_ID is not set")
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, sourceVMID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r BackupContainerVMAppTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := protectioncontainers.ParseProtectionContainerID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.RecoveryServices.BackupProtectionContainersClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r BackupContainerVMAppTestResource) basic(data acceptance.TestData, sourceVMID string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_container_vm_app" "test" {
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  source_resource_id  = "%s"
  resource_group_name = azurerm_resource_group.test.name
  workload_type       = "SAPAseDatabase"
}
`, r.template(data), sourceVMID)
}

func (r BackupContainerVMAppTestResource) requiresImport(data acceptance.TestData) string {
	sourceVMID := os.Getenv("ARM_TEST_SAP_VM_ID")
	return fmt.Sprintf(`
%s

resource "azurerm_backup_container_vm_app" "import" {
  recovery_vault_name = azurerm_backup_container_vm_app.test.recovery_vault_name
  source_resource_id  = azurerm_backup_container_vm_app.test.source_resource_id
  resource_group_name = azurerm_backup_container_vm_app.test.resource_group_name
  workload_type       = azurerm_backup_container_vm_app.test.workload_type
}
`, r.basic(data, sourceVMID))
}

func (r BackupContainerVMAppTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
