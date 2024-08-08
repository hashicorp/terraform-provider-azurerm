---
subcategory: "Extended Location"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_custom_location"
description: |-
  Gets information about an existing Custom Location.
---

# Data Source: azurerm_custom_location

Use this data source to access information about an existing Custom Location.

## Example Usage

```hcl
data "azurerm_custom_location" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_custom_location.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Custom Location.

* `resource_group_name` - (Required) The name of the Resource Group where the Custom Location exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Custom Location.

* `cluster_extension_ids` - A `cluster_extension_ids` block as defined below.

* `display_name` - The display name of the Custom Location.

* `host_resource_id` - The ID of the host resource of the Custom Location.

* `identities` - A `identities` block as defined below.

* `location` - The geo-location where the Custom Location exists.

* `namespace` - Kubernetes namespace that is created on the specified cluster for the Custom Location.

* `tags` - A mapping of tags assigned to the Custom Location.

---

A `identities` block exports the following:

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Custom Location.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Custom Location.

* `type` - The type of Managed Service Identity that is configured on this Custom Location.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Custom Location.
