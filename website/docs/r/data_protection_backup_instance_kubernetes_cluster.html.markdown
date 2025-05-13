---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_instance_kubernetes_cluster"
description: |-
  Manages a Backup Instance to back up a Kubernetes Cluster.
---

# azurerm_data_protection_backup_instance_kubernetes_cluster

Manages a Backup Instance to back up a Kubernetes Cluster.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_resource_group" "snap" {
  name     = "example-snap"
  location = "West Europe"
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "dns"

  default_node_pool {
    name                    = "default"
    node_count              = 1
    vm_size                 = "Standard_DS2_v2"
    host_encryption_enabled = true
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_trusted_access_role_binding" "aks_cluster_trusted_access" {
  kubernetes_cluster_id = azurerm_kubernetes_cluster.example.id
  name                  = "example"
  roles                 = ["Microsoft.DataProtection/backupVaults/backup-operator"]
  source_resource_id    = azurerm_data_protection_backup_vault.example.id
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "example"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_kubernetes_cluster_extension" "example" {
  name              = "example"
  cluster_id        = azurerm_kubernetes_cluster.example.id
  extension_type    = "Microsoft.DataProtection.Kubernetes"
  release_train     = "stable"
  release_namespace = "dataprotection-microsoft"
  configuration_settings = {
    "configuration.backupStorageLocation.bucket"                = azurerm_storage_container.example.name
    "configuration.backupStorageLocation.config.resourceGroup"  = azurerm_resource_group.example.name
    "configuration.backupStorageLocation.config.storageAccount" = azurerm_storage_account.example.name
    "configuration.backupStorageLocation.config.subscriptionId" = data.azurerm_client_config.current.subscription_id
    "credentials.tenantId"                                      = data.azurerm_client_config.current.tenant_id
  }
}

resource "azurerm_role_assignment" "test_extension_and_storage_account_permission" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Account Contributor"
  principal_id         = azurerm_kubernetes_cluster_extension.example.aks_assigned_identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_msi_read_on_cluster" {
  scope                = azurerm_kubernetes_cluster.example.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_msi_read_on_snap_rg" {
  scope                = azurerm_resource_group.snap.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_msi_snapshot_contributor_on_snap_rg" {
  scope                = azurerm_resource_group.snap.id
  role_definition_name = "Disk Snapshot Contributor"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_data_operator_on_snap_rg" {
  scope                = azurerm_resource_group.snap.id
  role_definition_name = "Data Operator for Managed Disks"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_vault_data_contributor_on_storage" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_cluster_msi_contributor_on_snap_rg" {
  scope                = azurerm_resource_group.snap.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_kubernetes_cluster.example.identity[0].principal_id
}

resource "azurerm_data_protection_backup_policy_kubernetes_cluster" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  vault_name          = azurerm_data_protection_backup_vault.example.name

  backup_repeating_time_intervals = ["R/2023-05-23T02:30:00+00:00/P1W"]

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
      scheduled_backup_times = ["2023-05-23T02:30:00Z"]
    }
  }

  default_retention_rule {
    life_cycle {
      duration        = "P14D"
      data_store_type = "OperationalStore"
    }
  }
}

resource "azurerm_data_protection_backup_instance_kubernetes_cluster" "example" {
  name                         = "example"
  location                     = azurerm_resource_group.example.location
  vault_id                     = azurerm_data_protection_backup_vault.example.id
  kubernetes_cluster_id        = azurerm_kubernetes_cluster.example.id
  snapshot_resource_group_name = azurerm_resource_group.snap.name
  backup_policy_id             = azurerm_data_protection_backup_policy_kubernetes_cluster.example.id

  backup_datasource_parameters {
    excluded_namespaces              = ["test-excluded-namespaces"]
    excluded_resource_types          = ["exvolumesnapshotcontents.snapshot.storage.k8s.io"]
    cluster_scoped_resources_enabled = true
    included_namespaces              = ["test-included-namespaces"]
    included_resource_types          = ["involumesnapshotcontents.snapshot.storage.k8s.io"]
    label_selectors                  = ["kubernetes.io/metadata.name:test"]
    volume_snapshot_enabled          = true
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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Backup Instance Kubernetes Cluster. Changing this forces a new resource to be created.

* `location` - (Required) The location of the Backup Instance Kubernetes Cluster. Changing this forces a new resource to be created.

* `vault_id` - (Required) The ID of the Backup Vault within which the Backup Instance Kubernetes Cluster should exist. Changing this forces a new resource to be created.

* `backup_policy_id` - (Required) The ID of the Backup Policy. Changing this forces a new resource to be created.

* `kubernetes_cluster_id` - (Required) The ID of the Kubernetes Cluster. Changing this forces a new resource to be created.

* `snapshot_resource_group_name` - (Required) The name of the Resource Group where snapshots are stored. Changing this forces a new resource to be created.

* `backup_datasource_parameters` - (Optional)  A `backup_datasource_parameters` block as defined below.

---

A `backup_datasource_parameters` block supports the following:

* `excluded_namespaces` - (Optional) Specifies the namespaces to be excluded during backup. Changing this forces a new resource to be created.

* `excluded_resource_types` - (Optional) Specifies the resource types to be excluded during backup. Changing this forces a new resource to be created.

* `cluster_scoped_resources_enabled` - (Optional) Whether to include cluster scope resources during backup. Default to `false`. Changing this forces a new resource to be created.

* `included_namespaces` - (Optional) Specifies the namespaces to be included during backup. Changing this forces a new resource to be created.

* `included_resource_types` - (Optional)  Specifies the resource types to be included during backup. Changing this forces a new resource to be created.

* `label_selectors` - (Optional) Specifies the resources with such label selectors to be included during backup. Changing this forces a new resource to be created.

* `volume_snapshot_enabled` - (Optional) Whether to take volume snapshots during backup. Default to `false`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Instance Kubernetes Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Instance Kubernetes Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Instance Kubernetes Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Instance Kubernetes Cluster.

## Import

Backup Instance Kubernetes Cluster can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_instance_kubernetes_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupInstances/backupInstance1
```
