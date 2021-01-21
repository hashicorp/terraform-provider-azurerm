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

data "azurerm_search_service" "example" {
  name                = "example-search-service"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this Search Service. Changing this forces a new Search Service to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Search Service should exist. Changing this forces a new Search Service to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Search Service.

* `primary_key` - The Primary Key used for Search Service Administration.

* `secondary_key` - The Secondary Key used for Search Service Administration.

* `query_keys` - A `query_keys` block as defined below.

* `public_network_access_enabled` - Whether or not public network access is enabled for this resource.

* `partition_count` - The number of partitions which have been created.

* `replica_count` - The number of replica's which have been created.

* `identity` - An `identity` block as defined below.

---

A `identity` block supports the following:

* `type` - The Type of Managed Identity which is used for the Search Service.

---

A `query_keys` block exports the following:

* `key` - The value of this Query Key.

* `name` - The name of this Query Key.

---

A `identity` block exports the following:

* `principal_id` - The (Client) ID of the Service Principal.

* `tenant_id` - The ID of the Tenant the Service Principal is assigned in.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Search Service.
