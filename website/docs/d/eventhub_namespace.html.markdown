---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_namespace"
description: |-
  Gets information about an existing EventHub Namespace.
---

# Data Source: azurerm_eventhub_namespace

Use this data source to access information about an existing EventHub Namespace.

## Example Usage

```hcl
data "azurerm_eventhub_namespace" "example" {
  name                = "search-eventhubns"
  resource_group_name = "search-service"
}

output "eventhub_namespace_id" {
  value = data.azurerm_eventhub_namespace.example.id
}
```

## Argument Reference

* `name` - The name of the EventHub Namespace.
* `resource_group_name` - The Name of the Resource Group where the EventHub Namespace exists.

## Attributes Reference

* `id` - The ID of the EventHub Namespace.

* `location` - The Azure location where the EventHub Namespace exists

* `sku` - Defines which tier to use.

* `capacity` - The Capacity / Throughput Units for a `Standard` SKU namespace.

* `auto_inflate_enabled` - Is Auto Inflate enabled for the EventHub Namespace?

* `maximum_throughput_units` -  Specifies the maximum number of throughput units when Auto Inflate is Enabled.

* `dedicated_cluster_id` - The ID of the EventHub Dedicated Cluster where this Namespace exists.

* `local_authentication_enabled` - Is this EventHub Namespace SAS authentication enabled?

* `tags` - A mapping of tags to assign to the EventHub Namespace.

The following attributes are exported only if there is an authorization rule named
`RootManageSharedAccessKey` which is created automatically by Azure.

* `default_primary_connection_string` - The primary connection string for the authorization
    rule `RootManageSharedAccessKey`.

* `default_primary_connection_string_alias` - The alias of the primary connection string for the authorization
    rule `RootManageSharedAccessKey`.

* `default_primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_connection_string` - The secondary connection string for the
    authorization rule `RootManageSharedAccessKey`.

* `default_secondary_connection_string_alias` - The alias of the secondary connection string for the
    authorization rule `RootManageSharedAccessKey`.

* `default_secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Namespace.
