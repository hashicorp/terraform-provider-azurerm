---
subcategory: "Azure Managed Lustre File System"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_lustre_file_system_auto_export_job"
description: |-
  Manages an Azure Managed Lustre File System Auto Export Job.
---

# azurerm_managed_lustre_file_system_auto_export_job

Manages an Azure Managed Lustre File System Auto Export Job.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-identity"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_key_vault" "example" {
  name                     = "example-keyvault"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "server" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.example.principal_id

  key_permissions = ["Get", "List", "WrapKey", "UnwrapKey", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_key" "example" {
  name         = "example"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.server,
  ]
}

resource "azurerm_storage_account" "example" {
  name                            = "example-storage-account"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "example" {
  name                  = "storagecontainer"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_storage_container" "example2" {
  name                  = "storagecontainer2"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

data "azuread_service_principal" "example" {
  display_name = "HPC Cache Resource Provider"
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Account Contributor"
  principal_id         = data.azuread_service_principal.example.object_id
}

resource "azurerm_role_assignment" "example2" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azuread_service_principal.example.object_id
}

resource "azurerm_managed_lustre_file_system" "example" {
  name                   = "example-amlfs"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  sku_name               = "AMLFS-Durable-Premium-250"
  subnet_id              = azurerm_subnet.example.id
  storage_capacity_in_tb = 8
  zones                  = ["2"]

  maintenance_window {
    day_of_week     = "Friday"
    time_of_day_utc = "22:00"
  }

  identity {
    type = "UserAssigned"

    identity_ids = [
      azurerm_user_assigned_identity.example.id
    ]
  }

  encryption_key {
    key_url         = azurerm_key_vault_key.example.id
    source_vault_id = azurerm_key_vault.example.id
  }

  hsm_setting {
    container_id         = azurerm_storage_container.example.resource_manager_id
    logging_container_id = azurerm_storage_container.example2.resource_manager_id
    import_prefix        = "/"
  }

  tags = {
    Env = "Test"
  }

  depends_on = [azurerm_role_assignment.example, azurerm_role_assignment.example2]
}

resource "azurerm_managed_lustre_file_system_auto_export_job" "example" {
  name                          = "acctest-amlfs-auto-export-job"
  managed_lustre_file_system_id = azurerm_managed_lustre_file_system.example.id
  location                      = azurerm_resource_group.example.location

  auto_export_prefixes = ["/"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Managed Lustre File System Auto Export Job. Changing this forces a new resource to be created.

* `managed_lustre_file_system_id` - (Required) The ID of the Azure Managed Lustre File System to which this Auto Export Job belongs. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Azure Managed Lustre File System Auto Export Job should exist. Changing this forces a new resource to be created.

* `auto_export_prefixes` - (Required) A list of prefixes that get auto exported to the cluster namespace. Changing this forces a new resource to be created.

* `admin_status_enabled` - (Optional) Whether the administrative status of the Auto Export Job is enabled. Defaults to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Managed Lustre File System Auto Export Job.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The Azure Managed Lustre File System Auto Export Job ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Azure Managed Lustre File System Auto Export Job.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Managed Lustre File System Auto Export Job.
* `update` - (Defaults to 1 hour) Used when updating the Azure Managed Lustre File System Auto Export Job.
* `delete` - (Defaults to 90 minutes) Used when deleting the Azure Managed Lustre File System Auto Export Job.

## Import

Azure Managed Lustre File Systems Auto Export Job can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_lustre_file_system.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageCache/amlFilesystems/amlFilesystem1/autoExportJobs/autoexportjob1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.StorageCache` - 2024-07-01
