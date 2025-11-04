---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_resource_anchor"
description: |-
  Gets information about an existing Oracle Resource Anchor.
---

# Data Source: azurerm_oracle_resource_anchor

Use this data source to access information about an existing Oracle Resource Anchor.

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

* `name` - (Required) The name of this Oracle Resource Anchor.

* `resource_group_name` - (Required) The name of the Resource Group where the Oracle Resource Anchor exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Oracle Resource Anchor.

* `linked_compartment_id` - Oracle Cloud Infrastructure compartment [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) which was created or linked by customer with Resource Anchor.

* `location` - The Azure Region where the Oracle Resource Anchor exists.

* `tags` - A mapping of tags assigned to the Oracle Resource Anchor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Oracle Resource Anchor.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
