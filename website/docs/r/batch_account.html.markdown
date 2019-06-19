---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_account"
sidebar_current: "docs-azurerm-resource-batch-account"
description: |-
  Manages an Azure Batch account.

---

# azurerm_batch_account

Manages an Azure Batch account.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "testbatch"
  location = "westeurope"
}

resource "azurerm_storage_account" "test" {
  name                     = "teststorage"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testbatchaccount"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
  storage_account_id   = "${azurerm_storage_account.test.id}"

  tags = {
    env = "test"
  }
}
```

## Example Usage with User Subscription mode

It's possible to deploy Azure Batch Account in User Subscription mode. In this mode, all the machines that will be created by batch pools will be created in the user Azure subscription. In this mode, you need to specify a reference to an Azure Key Vault that will be used by Azure Batch to store and retrieve sensitive information. You can read more about User Subscription mode in Azure Batch on [this page](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#account).

~> **NOTE:** the script below uses the also the [AzureAD provider](https://www.terraform.io/docs/providers/azuread/) to retrieve the "Microsoft Azure Batch" service principal information.


```hcl
resource "azurerm_resource_group" "example" {
  name     = "batch-rg"
  location = "westeurope"
}

# Get Microsoft Azure Batch service principal reference using the AzureAD provider
data "azuread_service_principal" "batchsp" {
  display_name = "Microsoft Azure Batch"
}

resource "azurerm_key_vault" "example" {
  name                            = "batchkv"
  location                        = "${azurerm_resource_group.example.location}"
  resource_group_name             = "${azurerm_resource_group.example.name}"
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  tenant_id                       = "00000000-0000-0000-0000-000000000002"

  sku {
    name = "standard"
  }

  access_policy {
    tenant_id = "00000000-0000-0000-0000-000000000002"
    object_id = "${data.azuread_service_principal.batchsp.object_id}"

    secret_permissions = [
      "get",
      "list",
      "set",
      "delete"
    ]

  }
}

resource "azurerm_role_assignment" "contribrole" {
  scope                = "/subscriptions/00000000-0000-0000-0000-000000000001"
  role_definition_name = "Contributor"
  principal_id         = "${data.azuread_service_principal.batchsp.object_id}"
}

resource "azurerm_batch_account" "example" {
  name                 = "batchaccount"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  location             = "${azurerm_resource_group.example.location}"
  storage_account_id   = "${azurerm_storage_account.example.id}"
  pool_allocation_mode = "UserSubscription"
  
  # reference the Azure KeyVault
  key_vault_reference {
    id  = "${azurerm_key_vault.example.id}"
    url = "${azurerm_key_vault.example.vault_uri}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Batch account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Batch account. Changing this forces a new resource to be created.

~> **NOTE:** To work around [a bug in the Azure API](https://github.com/Azure/azure-rest-api-specs/issues/5574) this property is currently treated as case-insensitive. A future version of Terraform will require that the casing is correct.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `pool_allocation_mode` - (Optional) Specifies the mode to use for pool allocation. Possible values are `BatchService` or `UserSubscription`. Defaults to `BatchService`.

~> **NOTE:** When using `UserSubscription` mode, an Azure KeyVault reference has to be specified. See `key_vault_reference` below.

* `key_vault_reference` - (Optional) A `key_vault_reference` block that describes the Azure KeyVault reference to use when deploying the Azure Batch account using the `UserSubscription` pool allocation mode. 

* `storage_account_id` - (Optional) Specifies the storage account to use for the Batch account. If not specified, Azure Batch will manage the storage.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `key_vault_reference` block supports the following:

* `id` - (Required) The Azure identifier of the Azure KeyVault to use.

* `url` - (Required) The HTTPS URL of the Azure KeyVault to use.

---

## Attributes Reference

The following attributes are exported:

* `id` - The Batch account ID.

* `primary_access_key` - The Batch account primary access key.

* `secondary_access_key` - The Batch account secondary access key.

* `account_endpoint` - The account endpoint used to interact with the Batch service.

~> **NOTE:** Primary and secondary access keys are only available when `pool_allocation_mode` is set to `BatchService`. See [documentation](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics) for more information.