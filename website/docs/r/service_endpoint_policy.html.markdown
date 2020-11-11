---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_endpoint_policy"
description: |-
  Manages a Service Endpoint Policy.
---

# azurerm_service_endpoint_policy

Manages a Service Endpoint Policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_service_endpoint_policy" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Service Endpoint Policy. Changing this forces a new Service Endpoint Policy to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Service Endpoint Policy should exist. Changing this forces a new Service Endpoint Policy to be created.

* `location` - (Required) The Azure Region where the Service Endpoint Policy should exist. Changing this forces a new Service Endpoint Policy to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Service Endpoint Policy.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Service Endpoint Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Endpoint Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Endpoint Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Service Endpoint Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Endpoint Policy.

## Import

Service Endpoint Policys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_service_endpoint_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/serviceEndpointPolicies/policy1
```
