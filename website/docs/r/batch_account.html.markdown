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

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Batch account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Batch account. Changing this forces a new resource to be created.

~> **NOTE:** To work around [a bug in the Azure API](https://github.com/Azure/azure-rest-api-specs/issues/5574) this property is currently treated as case-insensitive. A future version of Terraform will require that the casing is correct.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `pool_allocation_mode` - (Optional) Specifies the mode to use for pool allocation. Possible values are `BatchService` or `UserSubscription`. Defaults to `BatchService`.

* `storage_account_id` - (Optional) Specifies the storage account to use for the Batch account. If not specified, Azure Batch will manage the storage.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Batch account ID.

* `primary_access_key` - The Batch account primary access key.

* `secondary_access_key` - The Batch account secondary access key.

* `account_endpoint` - The account endpoint used to interact with the Batch service.

~> **NOTE:** Primary and secondary access keys are only available when `pool_allocation_mode` is set to `BatchService`. See [documentation](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics) for more information.