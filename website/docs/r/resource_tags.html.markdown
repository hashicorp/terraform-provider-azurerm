---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_tags"
description: |-
  Manages Resource Tags.
---

# azurerm_resource_tags

Manages Azure Resource Tags.
~> **Note** You must ignore the tags property on the resource you are trying to manage the tags on. Below is an example of this.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                = "tagsstorage132428"
  resource_group_name = azurerm_resource_group.example.name

  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  lifecycle {
	ignore_changes = [tags]
  }
}

resource "azurerm_resource_tags" "example" {
  resource_id = azurerm_storage_account.example.id
  
  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `resource_id` - (Required) The ID of the resource you want to manage the tags of. Changing this forces a new Resource Tags resource to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Resource Tags.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Resource Tags resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource Tags.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Tags.
* `update` - (Defaults to 30 minutes) Used when updating the Resource Tags.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Tags.

## Import

Resource Tagss can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_tags.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-rg/providers/Microsoft.Compute/virtualMachines/test-vm/providers/Microsoft.Resources/tags/default
```