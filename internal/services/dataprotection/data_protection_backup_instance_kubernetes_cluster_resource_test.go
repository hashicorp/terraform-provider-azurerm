// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupInstanceKubernetesClusterResource struct{}

/*
After creating the backup instance Azure will take a minute to actually configure the protection.
While the protection state is still "ConfiguringProtection" any attempt to delete the instance
will result in "406 Not Acceptable" preventing a clean post-test destroy and leaving dangling resources.
This step delays the destroy a bit and also has the added benefit of telling us if the protection is actually working.
*/
func (r DataProtectionBackupInstanceKubernetesClusterResource) waitForProtectionConfiguredStep(data acceptance.TestData, config string) acceptance.TestStep {
	return acceptance.TestStep{
		PreConfig: func() { time.Sleep(5 * time.Minute) },
		Config:    config,
		Check: acceptance.ComposeTestCheckFunc(
			check.That(data.ResourceName).Key("protection_state").HasValue("ProtectionConfigured"),
		),
	}
}

func TestAccDataProtectionBackupInstanceKubernetesCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_kubernetes_cluster", "test")
	r := DataProtectionBackupInstanceKubernetesClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		r.waitForProtectionConfiguredStep(data, r.basic(data)),
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupInstanceKubernetesCluster_crossSubscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_kubernetes_cluster", "test")

	if data.Subscriptions.Secondary == "" {
		t.Skip("Skipping: Test requires `ARM_SUBSCRIPTION_ID_ALT` environment variable to be specified")
	}

	r := DataProtectionBackupInstanceKubernetesClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossSubscription(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		r.waitForProtectionConfiguredStep(data, r.crossSubscription(data)),
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupInstanceKubernetesCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_kubernetes_cluster", "test")
	r := DataProtectionBackupInstanceKubernetesClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		r.waitForProtectionConfiguredStep(data, r.basic(data)),
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDataProtectionBackupInstanceKubernetesCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_kubernetes_cluster", "test")
	r := DataProtectionBackupInstanceKubernetesClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		r.waitForProtectionConfiguredStep(data, r.complete(data)),
		data.ImportStep(),
	})
}

func (r DataProtectionBackupInstanceKubernetesClusterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backupinstanceresources.ParseBackupInstanceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupInstanceClient.BackupInstancesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r DataProtectionBackupInstanceKubernetesClusterResource) template(data acceptance.TestData, crossSubscription bool) string {
	altSubscriptionId := data.Subscriptions.Primary
	if crossSubscription {
		altSubscriptionId = data.Subscriptions.Secondary
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azurerm-alt" {
  subscription_id = "%[4]s"
  features {}
}

data "azurerm_client_config" "current_alt" {
  provider = azurerm-alt
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dp-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test_alt" {
  provider = azurerm-alt
  name     = "acctest-dp-alt-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "snap" {
  provider = azurerm-alt
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
  provider            = azurerm-alt
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test_alt.location
  resource_group_name = azurerm_resource_group.test_alt.name
  dns_prefix          = "acctestaks%[1]d"

  default_node_pool {
    name                    = "default"
    node_count              = 1
    vm_size                 = "Standard_D2s_v5"
    host_encryption_enabled = false
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  node_provisioning_profile {
    mode               = "Manual"
    default_node_pools = "Auto"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_trusted_access_role_binding" "test_aks_cluster_trusted_access" {
  provider              = azurerm-alt
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  name                  = "mayankta"
  roles                 = ["Microsoft.DataProtection/backupVaults/backup-operator"]
  source_resource_id    = azurerm_data_protection_backup_vault.test.id
}

resource "azurerm_storage_account" "test" {
  provider                 = azurerm-alt
  name                     = "acctest%[3]s"
  resource_group_name      = azurerm_resource_group.test_alt.name
  location                 = azurerm_resource_group.test_alt.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  provider              = azurerm-alt
  name                  = "testaccsc%[3]s"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}

resource "azurerm_kubernetes_cluster_extension" "test" {
  provider          = azurerm-alt
  name              = "acctest-kce-%[1]d"
  cluster_id        = azurerm_kubernetes_cluster.test.id
  extension_type    = "Microsoft.DataProtection.Kubernetes"
  release_train     = "stable"
  release_namespace = "dataprotection-microsoft"
  configuration_settings = {
    "configuration.backupStorageLocation.bucket"                = azurerm_storage_container.test.name
    "configuration.backupStorageLocation.config.resourceGroup"  = azurerm_resource_group.test_alt.name
    "configuration.backupStorageLocation.config.storageAccount" = azurerm_storage_account.test.name
    "configuration.backupStorageLocation.config.subscriptionId" = data.azurerm_client_config.current_alt.subscription_id
    "credentials.tenantId"                                      = data.azurerm_client_config.current_alt.tenant_id
  }
}

resource "azurerm_role_assignment" "test_extension_and_storage_account_permission" {
  provider             = azurerm-alt
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Account Contributor"
  principal_id         = azurerm_kubernetes_cluster_extension.test.aks_assigned_identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_msi_read_on_cluster" {
  provider             = azurerm-alt
  scope                = azurerm_kubernetes_cluster.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_msi_read_on_snap_rg" {
  provider             = azurerm-alt
  scope                = azurerm_resource_group.snap.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_cluster_msi_contributor_on_snap_rg" {
  provider             = azurerm-alt
  scope                = azurerm_resource_group.snap.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_kubernetes_cluster.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_data_contributor_on_storage" {
  provider             = azurerm-alt
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_msi_snapshot_contributor_on_snap_rg" {
  provider             = azurerm-alt
  scope                = azurerm_resource_group.snap.id
  role_definition_name = "Disk Snapshot Contributor"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_data_operator_on_snap_rg" {
  provider             = azurerm-alt
  scope                = azurerm_resource_group.snap.id
  role_definition_name = "Data Operator for Managed Disks"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

// role propagation needs to complete before we can attempt to configure the protection
resource "time_sleep" "wait_5_minutes" {
  create_duration = "5m"

  depends_on = [
    azurerm_kubernetes_cluster_trusted_access_role_binding.test_aks_cluster_trusted_access,
    azurerm_role_assignment.test_extension_and_storage_account_permission,
    azurerm_role_assignment.test_vault_msi_read_on_cluster,
    azurerm_role_assignment.test_vault_msi_read_on_snap_rg,
    azurerm_role_assignment.test_cluster_msi_contributor_on_snap_rg,
    azurerm_role_assignment.test_vault_msi_snapshot_contributor_on_snap_rg,
    azurerm_role_assignment.test_vault_data_operator_on_snap_rg,
    azurerm_role_assignment.test_vault_data_contributor_on_storage,
  ]
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
    time_sleep.wait_5_minutes
  ]
}
	`, data.RandomInteger, data.Locations.Primary, data.RandomString, altSubscriptionId)
}

func (r DataProtectionBackupInstanceKubernetesClusterResource) requiresImport(data acceptance.TestData) string {
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

  backup_datasource_parameters {}
}
`, config)
}

func (r DataProtectionBackupInstanceKubernetesClusterResource) basic(data acceptance.TestData) string {
	template := r.template(data, false)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_instance_kubernetes_cluster" "test" {
  name                         = "acctest-iaks-%[2]d"
  location                     = azurerm_resource_group.test.location
  vault_id                     = azurerm_data_protection_backup_vault.test.id
  backup_policy_id             = azurerm_data_protection_backup_policy_kubernetes_cluster.test.id
  kubernetes_cluster_id        = azurerm_kubernetes_cluster.test.id
  snapshot_resource_group_name = azurerm_resource_group.snap.name

  backup_datasource_parameters {}
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupInstanceKubernetesClusterResource) crossSubscription(data acceptance.TestData) string {
	template := r.template(data, true)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_instance_kubernetes_cluster" "test" {
  name                         = "acctest-iaks-%[2]d"
  location                     = azurerm_resource_group.test.location
  vault_id                     = azurerm_data_protection_backup_vault.test.id
  backup_policy_id             = azurerm_data_protection_backup_policy_kubernetes_cluster.test.id
  kubernetes_cluster_id        = azurerm_kubernetes_cluster.test.id
  snapshot_subscription_id     = "%[3]s"
  snapshot_resource_group_name = azurerm_resource_group.snap.name

  backup_datasource_parameters {}
}
`, template, data.RandomInteger, data.Subscriptions.Secondary)
}

func (r DataProtectionBackupInstanceKubernetesClusterResource) complete(data acceptance.TestData) string {
	template := r.template(data, false)
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
    label_selectors                  = ["kubernetes.io/metadata.name=test"]
    volume_snapshot_enabled          = false
  }
}
`, template, data.RandomInteger)
}
