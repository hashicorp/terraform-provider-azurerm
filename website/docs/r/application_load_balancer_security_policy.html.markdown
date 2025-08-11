---
subcategory: "Service Networking"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_load_balancer_security_policy"
description: |-
  Manages an Application Load Balancer Security Policy.
---

# azurerm_application_load_balancer_security_policy

Manages an Application Load Balancer Security Policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_application_load_balancer" "example" {
  name                = "example-alb"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_web_application_firewall_policy" "example" {
  name                = "example-wafpolicy"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  managed_rules {
    managed_rule_set {
      type    = "Microsoft_DefaultRuleSet"
      version = "2.1"
    }
  }
  policy_settings {
    enabled = true
    mode    = "Detection"
  }
}

resource "azurerm_application_load_balancer_security_policy" "example" {
  name                               = "example-albsp"
  application_load_balancer_id       = azurerm_application_load_balancer.example.id
  location                           = azurerm_resource_group.example.location
  web_application_firewall_policy_id = azurerm_web_application_firewall_policy.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Application Load Balancer Security Policy. Changing this forces a new Application Load Balancer Security Policy to be created.

* `application_load_balancer_id` - (Required) The ID of the Application Load Balancer. Changing this forces a new Application Gateway for Containers Security Policy to be created.

* `web_application_firewall_policy_id` - (Required) The ID of the Web Application Firewall Policy. Changing this forces a new Application Load Balancer Security Policy to be created.

* `location` - (Required) The Azure Region where the Application Load Balancer Security Policy should exist. Changing this forces a new resource to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Application Load Balancer Security Policy.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Application Load Balancer Security Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Application Load Balancer Security Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Load Balancer Security Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Application Load Balancer Security Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Load Balancer Security Policy.

## Import

Application Load Balancer Security Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_load_balancer_security_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/alb/securityPolicies/sp1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ServiceNetworking` - 2025-01-01
