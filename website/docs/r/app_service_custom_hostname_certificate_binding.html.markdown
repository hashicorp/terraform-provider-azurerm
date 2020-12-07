---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_custom_hostname_certificate_binding"
description: |-
  Manages a App Service Custom Hostname Certificate Binding.
---

# azurerm_app_service_custom_hostname_certificate_binding

Manages a App Service Custom Hostname Certificate Binding.

## Example Usage

```hcl
resource "azurerm_app_service_custom_hostname_certificate_binding" "example" {
  hostname_binding_id = "TODO"
  certificate_id = "TODO"
  ssl_state = "TODO"
}
```

## Arguments Reference

The following arguments are supported:

* `certificate_id` - (Required) The ID of the TODO. Changing this forces a new App Service Custom Hostname Certificate Binding to be created.

* `hostname_binding_id` - (Required) The ID of the TODO. Changing this forces a new App Service Custom Hostname Certificate Binding to be created.

* `ssl_state` - (Required) TODO. Changing this forces a new App Service Custom Hostname Certificate Binding to be created.

---

* `resource_group_name` - (Optional) The name of the Resource Group where the App Service Custom Hostname Certificate Binding should exist. Changing this forces a new App Service Custom Hostname Certificate Binding to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the App Service Custom Hostname Certificate Binding.

* `app_service_name` - TODO.

* `hostname` - TODO.

* `thumbprint` - TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Custom Hostname Certificate Binding.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Custom Hostname Certificate Binding.
* `update` - (Defaults to 30 minutes) Used when updating the App Service Custom Hostname Certificate Binding.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Custom Hostname Certificate Binding.

## Import

App Service Custom Hostname Certificate Bindings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_custom_hostname_certificate_binding.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1/hostNameBindings/mywebsite.com
```