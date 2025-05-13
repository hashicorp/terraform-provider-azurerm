---
subcategory: "Dev Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_test_lab"
description: |-
  Manages a Dev Test Lab.
---

# azurerm_dev_test_lab

Manages a Dev Test Lab.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dev_test_lab" "example" {
  name                = "example-devtestlab"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    "Sydney" = "Australia"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Dev Test Lab. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the Dev Test Lab resource has to be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Dev Test Lab should exist. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Test Lab.

* `artifacts_storage_account_id` - The ID of the Storage Account used for Artifact Storage.

* `default_storage_account_id` - The ID of the Default Storage Account for this Dev Test Lab.

* `default_premium_storage_account_id` - The ID of the Default Premium Storage Account for this Dev Test Lab.

* `key_vault_id` - The ID of the Key used for this Dev Test Lab.

* `premium_data_disk_storage_account_id` - The ID of the Storage Account used for Storage of Premium Data Disk.

* `unique_identifier` - The unique immutable identifier of the Dev Test Lab.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DevTest Lab.
* `read` - (Defaults to 5 minutes) Used when retrieving the DevTest Lab.
* `update` - (Defaults to 30 minutes) Used when updating the DevTest Lab.
* `delete` - (Defaults to 30 minutes) Used when deleting the DevTest Lab.

## Import

Dev Test Labs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dev_test_lab.lab1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1
```
