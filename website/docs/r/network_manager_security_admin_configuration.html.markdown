---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_security_admin_configuration"
description: |-
  Manages a Network Manager Security Admin Configuration.
---

# azurerm_network_manager_security_admin_configuration

Manages a Network Manager Security Admin Configuration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {
}

resource "azurerm_network_manager" "example" {
  name                = "example-network-manager"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity", "SecurityAdmin"]
  description    = "example network manager"
}

resource "azurerm_network_manager_network_group" "example" {
  name               = "example-network-group"
  network_manager_id = azurerm_network_manager.example.id
}

resource "azurerm_network_manager_security_admin_configuration" "example" {
  name                                          = "example-admin-conf"
  network_manager_id                            = azurerm_network_manager.example.id
  description                                   = "example admin conf"
  apply_on_network_intent_policy_based_services = ["None"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Manager Security Admin Configuration. Changing this forces a new Network Manager Security Admin Configuration to be created.

* `network_manager_id` - (Required) Specifies the ID of the Network Manager Security Admin Configuration. Changing this forces a new Network Manager Security Admin Configuration to be created.

* `apply_on_network_intent_policy_based_services` - (Optional) A list of network intent policy based services. Possible values are `All`, `None` and `AllowRulesOnly`. Exactly one value should be set. The `All` option requires `Microsoft.Network/AllowAdminRulesOnNipBasedServices` feature registration to Subscription. Please see [this document](https://learn.microsoft.com/en-us/azure/virtual-network-manager/concept-security-admins#network-intent-policies-and-security-admin-rules) for more information.

* `description` - (Optional) A description of the Security Admin Configuration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Security Admin Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Security Admin Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Security Admin Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Security Admin Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Security Admin Configuration.

## Import

Network Manager Security Admin Configuration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_security_admin_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/securityAdminConfigurations/configuration1
```
