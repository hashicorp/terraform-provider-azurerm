---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_ip_restriction"
description: |-
  Manages a App Service IP Restriction.
---

# azurerm_app_service_ip_restriction

Manages a App Service IP Restriction.

## Disclaimers

~> **NOTE:** It's possible to define App Service IP Restrictions both within [the `azurerm_app_service` resource](azurerm_app_service.html) via `ip_restriction` blocks inside the `site_config` block and by using [the `azurerm_app_service_ip_restriction` resource](azurerm_app_service_ip_restriction.html). However it's not possible to use both methods to manage IP Restrictions for an App Service, since there'll be conflicts.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "example" {
  name                = "example-appserviceplan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "example-app-service"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id
}

resource "azurerm_app_service_ip_restriction" "example" {
  app_service_id = azurerm_app_service.example.id

  ip_restriction = {
    name       = "example"
    ip_address = "10.10.10.10/32"
    action     = "Allow"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `app_service_id` - (Required) The ID of the TODO. Changing this forces a new App Service IP Restriction to be created.

* `ip_restriction` - (Required) A `ip_restriction` block as defined below.

---

A `ip_restriction` block supports the following:

* `name` - (Required) The name for this IP Restriction.

* `action` - (Optional) Allow or Deny access for this IP range. Defaults to Allow.

* `headers` - (Optional) A `headers` block as defined blow.

* `ip_address` - (Optional) The IP Address used for this IP Restriction in CIDR notation.

* `priority` - (Optional) The priority for this IP Restriction. Restrictions are enforced in priority order. By default, priority is set to 65000 if not specified.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

-> **NOTE:** One of either `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified

---

A `headers` block supports the following:

* `x_azure_fdid` - (Optional) A list of allowed Azure FrontDoor IDs in UUID notation with a maximum of 8.

* `x_fd_health_probe` - (Optional) A list to allow the Azure FrontDoor health probe header. Only allowed value is "1".

* `x_forwarded_for` - (Optional) A list of allowed 'X-Forwarded-For' IPs in CIDR notation with a maximum of 8

* `x_forwarded_host` - (Optional) A list of allowed 'X-Forwarded-Host' domains with a maximum of 8.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the App Service IP Restriction.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service IP Restriction.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service IP Restriction.
* `update` - (Defaults to 30 minutes) Used when updating the App Service IP Restriction.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service IP Restriction.

## Import

App Service IP Restrictions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_ip_restriction.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Web/sites/siteName/ipRestriction/name
```
