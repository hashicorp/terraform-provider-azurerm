---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_resource_anchor"
description: |-
  Gets information about an existing oracle resource anchor.
---

# Data Source: azurerm_oracle_resource_anchor

Use this data source to access information about an existing oracle resource anchor.

## Example Usage

```hcl
data "azurerm_oracle_resource_anchor" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_oracle_resource_anchor.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this oracle resource anchor.

* `resource_group_name` - (Required) The name of the Resource Group where the oracle resource anchor exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the oracle resource anchor.

* `linked_compartment_id` - Oracle Cloud Infrastructure compartment Id (ocid) which was created or linked by customer with resource anchor.

* `location` - The Azure Region where the oracle resource anchor exists.

* `provisioning_state` - The provisioning state of the resource anchor.

* `tags` - A mapping of tags assigned to the oracle resource anchor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the oracle resource anchor.
