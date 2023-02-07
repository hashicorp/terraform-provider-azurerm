---
subcategory: "Search"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_search_service"
description: |-
  Manages a Search Service.
---

# azurerm_search_service

Manages a Search Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_search_service" "example" {
  name                = "example-search-service"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "standard"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Search Service should exist. Changing this forces a new Search Service to be created.

* `name` - (Required) The Name which should be used for this Search Service. Changing this forces a new Search Service to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Search Service should exist. Changing this forces a new Search Service to be created.

* `sku` - (Required) The SKU which should be used for this Search Service. Possible values are `basic`, `free`, `standard`, `standard2`, `standard3`, `storage_optimized_l1` and `storage_optimized_l2`. Changing this forces a new Search Service to be created.

-> The `basic` and `free` SKUs provision the Search Service in a Shared Cluster - the `standard` SKUs use a Dedicated Cluster.

~> **Note:** The SKUs `standard2` and `standard3` are only available when enabled on the backend by Microsoft.

---

* `public_network_access_enabled` - (Optional) Whether or not public network access is allowed for this resource. Defaults to `true`.

* `partition_count` - (Optional) The number of partitions which should be created.

* `replica_count` - (Optional) The number of replica's which should be created.

-> **Note:** `partition_count` and `replica_count` can only be configured when using a `standard` sku.

* `allowed_ips` - (Optional) A list of IPv4 addresses or CIDRs that are allowed access to the search service endpoint.

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Search Service.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Search Service. The only possible value is `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Search Service.

* `primary_key` - The Primary Key used for Search Service Administration.

* `query_keys` - A `query_keys` block as defined below.

* `secondary_key` - The Secondary Key used for Search Service Administration.

---

A `query_keys` block exports the following:

* `key` - The value of this Query Key.

* `name` - The name of this Query Key.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Search Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Search Service.
* `update` - (Defaults to 60 minutes) Used when updating the Search Service.
* `delete` - (Defaults to 60 minutes) Used when deleting the Search Service.

## Import

Search Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_search_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Search/searchServices/service1
```
