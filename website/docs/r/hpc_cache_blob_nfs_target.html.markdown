---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hpc_cache_blob_nfs_target"
description: |-
  Manages a Blob NFSv3 Target within a HPC Cache.
---

# azurerm_hpc_cache_blob_nfs_target

Manages a Blob NFSv3 Target within a HPC Cache.

~> **NOTE:**: By request of the service team the provider no longer automatically registering the `Microsoft.StorageCache` Resource Provider for this resource. To register it you can run `az provider register --namespace 'Microsoft.StorageCache'`.

~> **NOTE:**: This resource depends on the NFSv3 enabled Storage Account, which has some prerequisites need to meet. Please checkout: https://docs.microsoft.com/en-us/azure/storage/blobs/network-file-system-protocol-support-how-to?tabs=azure-powershell.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

provider "azuread" {}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
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
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

data "azuread_service_principal" "example" {
  display_name = "HPC Cache Resource Provider"
}

resource "azurerm_storage_account" "example" {
  name                      = "examplestorageaccount"
  resource_group_name       = azurerm_resource_group.example.name
  location                  = azurerm_resource_group.example.location
  account_tier              = "Standard"
  account_kind              = "StorageV2"
  account_replication_type  = "LRS"
  is_hns_enabled            = true
  nfsv3_enabled             = true
  enable_https_traffic_only = false
  network_rules {
    default_action             = "Deny"
    virtual_network_subnet_ids = [azurerm_subnet.example.id]
  }
}

# Due to https://github.com/terraform-providers/terraform-provider-azurerm/issues/2977 and the fact
# that the NFSv3 enabled storage account can't allow public network access - otherwise the NFSv3 protocol will fail,
# we have to use the ARM template to deploy the storage container as a workaround.
# Once the issue above got resolved, we can instead use the azurerm_storage_container resource.
resource "azurerm_resource_group_template_deployment" "storage-containers" {
  name                = "example-deployment"
  resource_group_name = azurerm_storage_account.example.resource_group_name
  deployment_mode     = "Incremental"

  parameters_content = jsonencode({
    location = {
      value = azurerm_storage_account.example.location
    },
    storageAccountName = {
      value = azurerm_storage_account.example.name
    },
    containerName = {
      value = "example-container"
    }
  })

  template_content = <<EOF
{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "storageAccountName": {
      "type": "String"
    },
    "containerName": {
      "type": "String"
    },
    "location": {
      "type": "String"
    }
  },
  "resources": [
    {
      "type": "Microsoft.Storage/storageAccounts",
      "apiVersion": "2019-06-01",
      "name": "[parameters('storageAccountName')]",
      "location": "[parameters('location')]",
      "sku": {
        "name": "Standard_LRS",
        "tier": "Standard"
      },
      "kind": "StorageV2",
      "properties": {
        "accessTier": "Hot"
      },
      "resources": [
        {
          "type": "blobServices/containers",
          "apiVersion": "2019-06-01",
          "name": "[concat('default/', parameters('containerName'))]",
          "dependsOn": [
            "[parameters('storageAccountName')]"
          ]
        }
      ]
    }
  ],

  "outputs": {
    "id": {
      "type": "String",
      "value": "[resourceId('Microsoft.Storage/storageAccounts/blobServices/containers', parameters('storageAccountName'), 'default', parameters('containerName'))]"
    }
  }
}
EOF
}

resource "azurerm_role_assignment" "example_storage_account_contrib" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Account Contributor"
  principal_id         = data.azuread_service_principal.example.object_id
}

resource "azurerm_role_assignment" "example_storage_blob_data_contrib" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azuread_service_principal.example.object_id
}

resource "azurerm_hpc_cache" "example" {
  name                = "example-hpc-cache"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.example.id
  sku_name            = "Standard_2G"
}

resource "azurerm_hpc_cache_blob_nfs_target" "example" {
  name                 = "example-hpc-target"
  resource_group_name  = azurerm_resource_group.example.name
  cache_name           = azurerm_hpc_cache.example.name
  storage_container_id = jsondecode(azurerm_resource_group_template_deployment.storage-containers.output_content).id.value
  namespace_path       = "/p1"
  usage_model          = "READ_HEAVY_INFREQ"
}
```

## Arguments Reference

The following arguments are supported:

* `cache_name` - (Required) The name of the HPC Cache, which the HPC Cache Blob NFS Target will be added to. Changing this forces a new HPC Cache Blob NFS Target to be created.

* `name` - (Required) The name which should be used for this HPC Cache Blob NFS Target. Changing this forces a new HPC Cache Blob NFS Target to be created.

* `namespace_path` - (Required) The client-facing file path of the HPC Cache Blob NFS Target.

* `resource_group_name` - (Required) The name of the Resource Group where the HPC Cache Blob NFS Target should exist. Changing this forces a new HPC Cache Blob NFS Target to be created.

* `storage_container_id` - (Required) The Resource Manager ID of the Storage Container used as the HPC Cache Blob NFS Target. Changing this forces a new resource to be created.

-> **Note:** This is the Resource Manager ID of the Storage Container, rather than the regular ID - and can be accessed on the `azurerm_storage_container` Data Source/Resource as `resource_manager_id`.

* `usage_model` - (Required) The type of usage of the HPC Cache Blob NFS Target. Possible values are: `READ_HEAVY_INFREQ`, `READ_HEAVY_CHECK_180`, `WRITE_WORKLOAD_15`, `WRITE_AROUND`, `WRITE_WORKLOAD_CHECK_30`, `WRITE_WORKLOAD_CHECK_60` and `WRITE_WORKLOAD_CLOUDWS`.

---

* `access_policy_name` - (Optional) The name of the access policy applied to this target. Defaults to `default`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the HPC Cache Blob NFS Target.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the HPC Cache Blob NFS Target.
* `read` - (Defaults to 5 minutes) Used when retrieving the HPC Cache Blob NFS Target.
* `update` - (Defaults to 30 minutes) Used when updating the HPC Cache Blob NFS Target.
* `delete` - (Defaults to 30 minutes) Used when deleting the HPC Cache Blob NFS Target.

## Import

HPC Cache Blob NFS Targets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hpc_cache_blob_nfs_target.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StorageCache/caches/cache1/storageTargets/target1
```
