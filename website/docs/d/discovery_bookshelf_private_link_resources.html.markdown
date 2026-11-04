---
subcategory: "Discovery"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_discovery_bookshelf_private_link_resources"
description: |-
  Gets information about the Private Link Resources supported by a Microsoft Discovery Bookshelf.
---

# Data Source: azurerm_discovery_bookshelf_private_link_resources

Use this data source to access information about the Private Link Resources supported by a Microsoft Discovery Bookshelf.

This is most commonly used to look up the `group_id` values which can be passed to the `subresource_names` property of an `azurerm_private_endpoint`.

## Example Usage

```hcl
data "azurerm_discovery_bookshelf_private_link_resources" "example" {
  bookshelf_id = azurerm_discovery_bookshelf.example.id
}

output "group_ids" {
  value = data.azurerm_discovery_bookshelf_private_link_resources.example.private_link_resources[*].group_id
}
```

## Arguments Reference

The following arguments are supported:

* `bookshelf_id` - (Required) The ID of the Microsoft Discovery Bookshelf.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Microsoft Discovery Bookshelf.

* `private_link_resources` - A list of `private_link_resources` blocks as defined below.

---

A `private_link_resources` block exports the following:

* `name` - The name of the Private Link Resource.

* `group_id` - The group ID of the Private Link Resource. This is the value which should be used for the `subresource_names` property of an `azurerm_private_endpoint`.

* `required_members` - A list of members which are required to establish a Private Endpoint Connection to this Private Link Resource.

* `required_zone_names` - A list of Private DNS Zone names which are required to resolve this Private Link Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft Discovery Bookshelf Private Link Resources.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Discovery` - 2026-06-01
