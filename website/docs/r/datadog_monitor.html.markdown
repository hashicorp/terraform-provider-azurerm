---
subcategory: "Datadog"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_datadog_monitor"
description: |-
  Manages a Datadog Monitor.
---

# azurerm_datadog_monitor

Manages a datadog Monitor.

## Example Usage

### Monitor creation with linking to Datadog organization

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-datadog"
  location = "West US 2"
}
resource "azurerm_datadog_monitor" "example" {
  name                = "example-monitor"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datadog_organization {
    api_key         = "XXXX"
    application_key = "XXXX"
  }
  user {
    name  = "Example"
    email = "abc@xyz.com"
  }
  sku_name = "Linked"
  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the user that will be associated with the Datadog Monitor. Changing this forces a new Datadog Monitor to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Datadog Monitor should exist. Changing this forces a new Datadog Monitor to be created.

* `location` - (Required) The Azure Region where the Datadog Monitor should exist. Changing this forces a new Datadog Monitor to be created.

* `sku_name` - (Required) The name which should be used for this sku.

* `identity` - (Optional) A `identity` block as defined below.

* `user` - (Required) A `user` block as defined below.

* `datadog_organization` - (Required) A `datadog_organization` block as defined below.

* `monitoring_enabled` - (Optional) Is monitoring enabled? Defaults to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Datadog Monitor.

---

A `datadog_organization` block exports the following:

* `api_key` - (Required) Api key associated to the Datadog organization. Changing this forces a new Datadog Monitor to be created.

* `application_key` - (Required) Application key associated to the Datadog organization. Changing this forces a new Datadog Monitor to be created.

* `enterprise_app_id` - (Optional) The ID of the enterprise_app. Changing this forces a new resource to be created.

* `linking_auth_code` - (Optional) The auth code used to linking to an existing Datadog organization. Changing this forces a new Datadog Monitor to be created.

* `linking_client_id` - (Optional) The ID of the linking_client. Changing this forces a new Datadog Monitor to be created.

* `redirect_uri` - (Optional) The redirect uri for linking. Changing this forces a new Datadog Monitor to be created.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Datadog Monitor. At this time the only allowed value is `SystemAssigned`.

-> **Note:** The assigned `principal_id` and `tenant_id` can be retrieved after the identity `type` has been set to `SystemAssigned` and the Datadog Monitor has been created. More details are available below.

---

An `user` block exports the following:

* `name` - (Required) The name which should be used for this user_info. Changing this forces a new resource to be created.

* `email` - (Required) Email of the user used by Datadog for contacting them if needed. Changing this forces a new Datadog Monitor to be created.

* `phone_number` - (Optional) Phone number of the user used by Datadog for contacting them if needed. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Datadog Monitor.

* `identity` - A `identity` block as defined below.

* `marketplace_subscription_status` - Flag specifying the Marketplace Subscription Status of the resource. If payment is not made in time, the resource will go in Suspended state.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this Datadog Monitor.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this Datadog Monitor.

-> **Note:** You can access the Principal ID via `${azurerm_datadog_monitor.example.identity[0].principal_id}` and the Tenant ID via `${azurerm_datadog_monitor.example.identity[0].tenant_id}`

## Role Assignment

To enable metrics flow, perform role assignment on the identity created above. `Monitoring reader(43d0d8ad-25c7-4714-9337-8ba259a9fe05)` role is required .

### Role assignment on the monitor created

```hcl
data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "monitoring_reader" {
  name = "Monitoring Reader"
}

resource "azurerm_role_assignment" "example" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = data.azurerm_role_definition.monitoring_reader.role_definition_id
  principal_id       = azurerm_datadog_monitor.example.identity[0].principal_id
}
```

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Datadog Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Datadog Monitor.
* `update` - (Defaults to 30 minutes) Used when updating the Datadog Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the Datadog Monitor.

## Import

Datadog Monitors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_datadog_monitor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Datadog/monitors/monitor1
```
