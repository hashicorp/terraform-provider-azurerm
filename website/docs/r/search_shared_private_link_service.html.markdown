---
subcategory: "Search"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_search_shared_private_link_service"
description: |-
  Manages the Shared Private Link Service for an Azure Search Service.
---

# azurerm_search_shared_private_link_service

Manages the Shared Private Link Service for an Azure Search Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resourceGroup"
  location = "east us"
}
resource "azurerm_search_service" "test" {
  name                = "example-search"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
}
resource "azurerm_storage_account" "test" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
resource "azurerm_search_shared_private_link_service" "test" {
  name               = "example-spl"
  search_service_id  = azurerm_search_service.test.id
  subresource_name   = "blob"
  target_resource_id = azurerm_storage_account.test.id
  request_message    = "please approve"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specify the name of the Azure Search Shared Private Link Resource. Changing this forces a new resource to be created.

* `search_service_id` - (Required) Specify the id of the Azure Search Service. Changing this forces a new resource to be created.

* `subresource_name` - (Required) Specify the sub resource name which the Azure Search Private Endpoint is able to connect to. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) Specify the ID of the Shared Private Link Enabled Remote Resource which this Azure Search Private Endpoint should be connected to. Changing this forces a new resource to be created.

-> **Note:** The sub resource name should match with the type of the target resource id that's being specified.

* `request_message` - (Optional) Specify the request message for requesting approval of the Shared Private Link Enabled Remote Resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Search Shared Private Link resource.

* `status` - The status of a private endpoint connection. Possible values are Pending, Approved, Rejected or Disconnected.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Azure Search Shared Private Link Resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Search Shared Private Link Resource.
* `update` - (Defaults to 1 hour) Used when updating the Azure Search Shared Private Link Resource.
* `delete` - (Defaults to 1 hour) Used when deleting the Azure Search Shared Private Link Resource.

## Import

Azure Search Shared Private Link Resource can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_search_shared_private_link_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Search/searchServices/service1/sharedPrivateLinkResources/resource1
```
