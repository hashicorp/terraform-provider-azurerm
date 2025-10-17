---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_resource_anchor"
description: |-
  Manages an Oracle Resource Anchor.
---

# azurerm_oracle_resource_anchor

Manages an oracle resource anchor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "eastus"
}

resource "azurerm_oracle_resource_anchor" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = "eastus"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this oracle resource anchor. Changing this forces a new oracle resource anchor to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the oracle resource anchor should exist. Changing this forces a new oracle resource anchor to be created.

* `location` - (Required) The Azure Region where the oracle resource anchor should exist. Changing this forces a new oracle resource anchor to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the oracle resource anchor.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the oracle resource anchor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the oracle  resource anchor.
* `read` - (Defaults to 5 minutes) Used when retrieving the oracle resource anchor.
* `update` - (Defaults to 10 minutes) Used when updating the oracle resource anchor.
* `delete` - (Defaults to 10 minutes) Used when deleting the oracle resource anchor.

## Import

oracle resource anchors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_resource_anchor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Oracle.Database/resourceanchors/example
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
