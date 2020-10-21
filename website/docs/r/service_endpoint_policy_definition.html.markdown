---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_endpoint_policy_definition"
description: |-
  Manages a Service Endpoint Policy Definition.
---

# azurerm_service_endpoint_policy_definition

Manages a Service Endpoint Policy Definition.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_service_endpoint_policy" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_service_endpoint_policy_definition" "example" {
  name                  = "example"
  policy_id             = azurerm_service_endpoint_policy.example.id
  service_endpoint_name = "Microsoft.Storage"
  service_resources     = ["/subscriptions/00000000-0000-0000-0000-000000000000"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Service Endpoint Policy Definition. Changing this forces a new Service Endpoint Policy Definition to be created.

* `policy_id` - (Required) The ID of the Service Endpoint Policy which this Service Endpoint Policy Definition belongs to. Changing this forces a new Service Endpoint Policy Definition to be created.

* `service_endpoint_name` - (Required) The service endpoint name.

* `service_resources` - (Required) Specifies a list of resources that this Service Endpoint Policy Definition applies to.

---

* `description` - (Optional) The description of this Service Endpoint Policy Definition.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Service Endpoint Policy Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Endpoint Policy Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Endpoint Policy Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Service Endpoint Policy Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Endpoint Policy Definition.

## Import

Service Endpoint Policy Definitions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_service_endpoint_policy_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/serviceEndpointPolicies/policy1/serviceEndpointPolicyDefinitions/definition1
```
