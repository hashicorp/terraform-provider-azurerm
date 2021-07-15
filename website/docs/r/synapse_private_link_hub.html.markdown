---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_private_link_hub"
description: |-
  Manages a Synapse Private Link Hub.
---

# azurerm_synapse_private_link_hub

Manages a Synapse Private Link Hub.

## Example Usage

```hcl
resource "azurerm_synapse_private_link_hub" "example" {
  name                = "example-resource"
  resource_group_name = "example-rg"
  location            = "West Europe"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse Private Link Hub. Changing this forces a new Synapse Private Link Hub to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Synapse Private Link Hub. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure location where the Synapse Private Link Hub exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Synapse Private Link Hub.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Synapse Private Link Hub.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Private Link Hub.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Private Link Hub.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Private Link Hub.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Private Link Hub.

## Import

Synapse Private Link Hub can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_private_link_hub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Synapse/privateLinkHubs/privateLinkHub1
```
