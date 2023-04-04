---
subcategory: "storagemover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storagemover_endpoint"
description: |-
  Manages a Storagemover Endpoints.
---

# azurerm_storagemover_endpoint

Manages a Storagemover Endpoints.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storagemover_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_storagemover_endpoint" "example" {
  name                          = "example-se"
  storagemover_storage_mover_id = azurerm_storagemover_storage_mover.test.id
  description                   = ""
  endpoint_type                 = ""

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Storagemover Endpoints. Changing this forces a new Storagemover Endpoints to be created.

* `storagemover_storage_mover_id` - (Required) Specifies the ID of the Storagemover Endpoints. Changing this forces a new Storagemover Endpoints to be created.

* `endpoint_type` - (Required) Specifies the Endpoint resource type. Changing this forces a new Storagemover Endpoints to be created.

* `description` - (Optional) A description for the Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storagemover Endpoints.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storagemover Endpoints.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storagemover Endpoints.
* `update` - (Defaults to 30 minutes) Used when updating the Storagemover Endpoints.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storagemover Endpoints.

## Import

Storagemover Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storagemover_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1/endpoints/endpoint1
```
