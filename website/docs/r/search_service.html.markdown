---
subcategory: "Search"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_search_service"
sidebar_current: "docs-azurerm-resource-search-service"
description: |-
  Manages a Search Service.
---

# azurerm_search_service

Allows you to manage an Azure Search Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_search_service" "example" {
  name                = "acceptanceTestSearchService1"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  sku                 = "standard"

  tags = {
    environment = "staging"
    database    = "test"
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Search Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Search Service. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) Valid values are `basic`, `free` and `standard`. `standard2` and `standard3` are also valid, but can only be used when it's enabled on the backend by Microsoft support. `free` provisions the service in shared clusters. `standard` provisions the service in dedicated clusters.  Changing this forces a new resource to be created.

* `replica_count` - (Optional) Default is 1. Valid values include 1 through 12. Valid only when `sku` is `standard`. Changing this forces a new resource to be created.

* `partition_count` - (Optional) Default is 1. Valid values include 1, 2, 3, 4, 6, or 12. Valid only when `sku` is `standard`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Search Service ID.

* `primary_key` - The Search Service Administration primary key.

* `secondary_key` - The Search Service Administration secondary key.

* `query_keys` - A `query_keys` block as defined below.

---

A `query_keys` block exports the following:

* `name` - The name of the query key.

* `key` - The value of the query key.

---

## Import

Search Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_search_service.service1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Search/searchServices/service1
```
