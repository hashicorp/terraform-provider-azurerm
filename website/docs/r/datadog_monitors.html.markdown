---
subcategory: "Datadog"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_datadog_monitor"
description: |-
  Manages a datadog Monitor.
---

# azurerm_datadog_monitor

Manages a datadog Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-datadog"
  location = "West Europe"
}

resource "azurerm_datadog_monitor" "example" {
  name                = "example-monitor"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this datadog Monitor. Changing this forces a new datadog Monitor to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the datadog Monitor should exist. Changing this forces a new datadog Monitor to be created.

* `location` - (Required) The Azure Region where the datadog Monitor should exist. Changing this forces a new datadog Monitor to be created.

---

* `datadog_organization_properties` - (Optional) A `datadog_organization_properties` block as defined below.

* `identity` - (Optional) A `identity` block as defined below.

* `sku` - (Optional) A `sku` block as defined below.

* `user_info` - (Optional) A `user_info` block as defined below.

* `monitoring_status` - (Optional) Flag specifying if the resource monitoring is enabled or disabled. Possible values are "true" and "false" is allowed.

* `tags` - (Optional) A mapping of tags which should be assigned to the datadog Monitor.

---

An `datadog_organization_properties` block exports the following:

* `api_key` - (Optional) Api key associated to the Datadog organization. Changing this forces a new datadog Monitor to be created.

* `application_key` - (Optional) Application key associated to the Datadog organization. Changing this forces a new datadog Monitor to be created.

* `enterprise_app_id` - (Optional) The ID of the enterprise_app. Changing this forces a new datadog Monitor to be created.

* `linking_auth_code` - (Optional) The auth code used to linking to an existing datadog organization. Changing this forces a new datadog Monitor to be created.

* `linking_client_id` - (Optional) The ID of the linking_client. Changing this forces a new datadog Monitor to be created.

* `redirect_uri` - (Optional) The redirect uri for linking. Changing this forces a new datadog Monitor to be created.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the datadog Monitor. At this time the only allowed value is `SystemAssigned`.

~> **NOTE:** The assigned `principal_id` and `tenant_id` can be retrieved after the identity `type` has been set to `SystemAssigned` and the datadog Monitor has been created. More details are available below.

---

An `sku` block exports the following:

* `name` - (Required) The name which should be used for this sku. Changing this forces a new datadog Monitor to be created.

---

An `user_info` block exports the following:

* `name` - (Optional) The name which should be used for this user_info. Changing this forces a new datadog Monitor to be created.

* `email_address` - (Optional) Email of the user used by Datadog for contacting them if needed. Changing this forces a new datadog Monitor to be created.

* `phone_number` - (Optional) Phone number of the user used by Datadog for contacting them if needed. Changing this forces a new datadog Monitor to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the datadog Monitor.

* `identity` - A `identity` block as defined below.

* `liftr_resource_category` - .

* `liftr_resource_preference` - The priority of the resource.

* `marketplace_subscription_status` - Flag specifying the Marketplace Subscription Status of the resource. If payment is not made in time, the resource will go in Suspended state.

* `type` - The type of the monitor resource.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this datadog Monitor.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this datadog Monitor.

-> You can access the Principal ID via `${azurerm_datadog_monitor.example.identity.0.principal_id}` and the Tenant ID via `${azurerm_datadog_monitor.example.identity.0.tenant_id}`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the datadog Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the datadog Monitor.
* `update` - (Defaults to 30 minutes) Used when updating the datadog Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the datadog Monitor.

## Import

datadog Monitors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_datadog_monitor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Datadog/monitors/monitor1
```
