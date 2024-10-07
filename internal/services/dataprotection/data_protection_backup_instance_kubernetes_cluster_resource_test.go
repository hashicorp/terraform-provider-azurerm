// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupInstanceKubernatesClusterTestResource struct{}

func TestAccDataProtectionBackupInstanceKubernatesCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_kubernetes_cluster", "test")
	r := DataProtectionBackupInstanceKubernatesClusterTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupInstanceKubernatesCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_kubernetes_cluster", "test")
	r := DataProtectionBackupInstanceKubernatesClusterTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDataProtectionBackupInstanceKubernatesCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_kubernetes_cluster", "test")
	r := DataProtectionBackupInstanceKubernatesClusterTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DataProtectionBackupInstanceKubernatesClusterTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backupinstances.ParseBackupInstanceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupInstanceClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r DataProtectionBackupInstanceKubernatesClusterTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dp-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "snap" {
  name     = "acctest-dp-snap-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-dbv-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  soft_delete         = "Off"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  default_node_pool {
    name                    = "default"
    node_count              = 1
    vm_size                 = "Standard_DS2_v2"
    host_encryption_enabled = true
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_trusted_access_role_binding" "test_aks_cluster_trusted_access" {
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  name                  = "mayankta"
  roles                 = ["Microsoft.DataProtection/backupVaults/backup-operator"]
  source_resource_id    = azurerm_data_protection_backup_vault.test.id
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "testaccsc%[3]s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_kubernetes_cluster_extension" "test" {
  name              = "acctest-kce-%[1]d"
  cluster_id        = azurerm_kubernetes_cluster.test.id
  extension_type    = "Microsoft.DataProtection.Kubernetes"
  release_train     = "stable"
  release_namespace = "dataprotection-microsoft"
  configuration_settings = {
    "configuration.backupStorageLocation.bucket"                = azurerm_storage_container.test.name
    "configuration.backupStorageLocation.config.resourceGroup"  = azurerm_resource_group.test.name
    "configuration.backupStorageLocation.config.storageAccount" = azurerm_storage_account.test.name
    "configuration.backupStorageLocation.config.subscriptionId" = data.azurerm_client_config.current.subscription_id
    "credentials.tenantId"                                      = data.azurerm_client_config.current.tenant_id
  }
}

resource "azurerm_role_assignment" "test_extension_and_storage_account_permission" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Account Contributor"
  principal_id         = azurerm_kubernetes_cluster_extension.test.aks_assigned_identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_msi_read_on_cluster" {
  scope                = azurerm_kubernetes_cluster.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_msi_read_on_snap_rg" {
  scope                = azurerm_resource_group.snap.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_cluster_msi_contributor_on_snap_rg" {
  scope                = azurerm_resource_group.snap.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_kubernetes_cluster.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_data_contributor_on_storage" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_msi_snapshot_contributor_on_snap_rg" {
  scope                = azurerm_resource_group.snap.id
  role_definition_name = "Disk Snapshot Contributor"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_data_operator_on_snap_rg" {
  scope                = azurerm_resource_group.snap.id
  role_definition_name = "Data Operator for Managed Disks"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_data_protection_backup_policy_kubernetes_cluster" "test" {
  name                = "acctest-paks-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  vault_name          = azurerm_data_protection_backup_vault.test.name

  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  retention_rule {
    name     = "Daily"
    priority = 25

    life_cycle {
      duration        = "P84D"
      data_store_type = "OperationalStore"
    }

    criteria {
      days_of_week           = ["Thursday"]
      months_of_year         = ["November"]
      weeks_of_month         = ["First"]
      scheduled_backup_times = ["2021-05-23T02:30:00Z"]
    }
  }

  default_retention_rule {
    life_cycle {
      duration        = "P14D"
      data_store_type = "OperationalStore"
    }
  }

  depends_on = [
    azurerm_role_assignment.test_extension_and_storage_account_permission,
    azurerm_role_assignment.test_vault_msi_read_on_cluster,
    azurerm_role_assignment.test_vault_msi_read_on_snap_rg,
    azurerm_role_assignment.test_cluster_msi_contributor_on_snap_rg,
    azurerm_role_assignment.test_vault_msi_snapshot_contributor_on_snap_rg,
    azurerm_role_assignment.test_vault_data_operator_on_snap_rg,
    azurerm_role_assignment.test_vault_data_contributor_on_storage,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r DataProtectionBackupInstanceKubernatesClusterTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_kubernetes_cluster" "import" {
  name                         = azurerm_data_protection_backup_instance_kubernetes_cluster.test.name
  location                     = azurerm_data_protection_backup_instance_kubernetes_cluster.test.location
  vault_id                     = azurerm_data_protection_backup_instance_kubernetes_cluster.test.vault_id
  backup_policy_id             = azurerm_data_protection_backup_instance_kubernetes_cluster.test.backup_policy_id
  kubernetes_cluster_id        = azurerm_data_protection_backup_instance_kubernetes_cluster.test.kubernetes_cluster_id
  snapshot_resource_group_name = azurerm_data_protection_backup_instance_kubernetes_cluster.test.snapshot_resource_group_name
}
`, config)
}

func (r DataProtectionBackupInstanceKubernatesClusterTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_instance_kubernetes_cluster" "test" {
  name                         = "acctest-iaks-%[2]d"
  location                     = azurerm_resource_group.test.location
  vault_id                     = azurerm_data_protection_backup_vault.test.id
  backup_policy_id             = azurerm_data_protection_backup_policy_kubernetes_cluster.test.id
  kubernetes_cluster_id        = azurerm_kubernetes_cluster.test.id
  snapshot_resource_group_name = azurerm_resource_group.snap.name
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupInstanceKubernatesClusterTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_instance_kubernetes_cluster" "test" {
  name                         = "acctest-iaks-%[2]d"
  location                     = azurerm_resource_group.test.location
  vault_id                     = azurerm_data_protection_backup_vault.test.id
  backup_policy_id             = azurerm_data_protection_backup_policy_kubernetes_cluster.test.id
  kubernetes_cluster_id        = azurerm_kubernetes_cluster.test.id
  snapshot_resource_group_name = azurerm_resource_group.snap.name

  backup_datasource_parameters {
    excluded_namespaces              = ["test-excluded-namespaces"]
    excluded_resource_types          = ["exvolumesnapshotcontents.snapshot.storage.k8s.io"]
    cluster_scoped_resources_enabled = true
    included_namespaces              = ["test-included-namespaces"]
    included_resource_types          = ["involumesnapshotcontents.snapshot.storage.k8s.io"]
    label_selectors                  = ["kubernetes.io/metadata.name:test"]
    volume_snapshot_enabled          = false
  }

  depends_on = [
    azurerm_role_assignment.test_extension_and_storage_account_permission,
    azurerm_role_assignment.test_vault_msi_read_on_cluster,
    azurerm_role_assignment.test_vault_msi_read_on_snap_rg,
    azurerm_role_assignment.test_cluster_msi_contributor_on_snap_rg,
    azurerm_role_assignment.test_vault_msi_snapshot_contributor_on_snap_rg,
    azurerm_role_assignment.test_vault_data_operator_on_snap_rg,
    azurerm_role_assignment.test_vault_data_contributor_on_storage,
  ]
}
`, template, data.RandomInteger)
}
