---
subcategory: "Lab Service"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lab_service_lab_plan"
description: |-
  Manages a Lab Service Lab Plan.
---

# azurerm_lab_service_lab_plan

Manages a Lab Service Lab Plan.

-> **Note:** Before using this resource, it's required to submit the request of registering the provider with Azure CLI `az provider register --namespace Microsoft.LabServices`.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_lab_service_lab_plan" "example" {
  name                = "example-lp"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allowed_regions     = [azurerm_resource_group.example.location]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Lab Service Lab Plan. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Lab Service Lab Plan should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Lab Service Lab Plan should exist. Changing this forces a new resource to be created.

* `allowed_regions` - (Required) The allowed regions for the lab creator to use when creating labs using this Lab Service Lab Plan. The allowed region's count must be between `1` and `28`.

* `default_auto_shutdown` - (Optional) A `default_auto_shutdown` block as defined below.

* `default_connection` - (Optional) A `default_connection` block as defined below.

* `default_network` - (Optional) A `default_network` block as defined below.

* `shared_gallery_id` - (Optional) The resource ID of the Shared Image Gallery attached to this Lab Service Lab Plan. When saving a lab template virtual machine image it will be persisted in this gallery. The shared images from the gallery can be made available to use when creating new labs.

~> **NOTE:** The built-in `Azure Lab Services` Service Principal needs to be assigned to the Shared Image Gallery while using this property.

* `support` - (Optional) A `support` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Lab Service Lab Plan.

---

A `default_auto_shutdown` block supports the following:

* `disconnect_delay` - (Required) The amount of time a VM will stay running after a user disconnects if this behavior is enabled. This value must be formatted as an ISO 8601 string.

* `idle_delay` - (Required) The amount of time a VM will idle before it is shutdown if this behavior is enabled. This value must be formatted as an ISO 8601 string.

* `no_connect_delay` - (Required) The amount of time a VM will stay running before it is shutdown if no connection is made and this behavior is enabled. This value must be formatted as an ISO 8601 string.

* `shutdown_on_disconnect_enabled` - (Required) Is shutdown on disconnect enabled? Possible values are `true` and `false`.

* `shutdown_on_idle` - (Required) Will a VM get shutdown when it has idled for a period of time? Possible values are `UserAbsence`, `LowUsage` and `None`.

* `shutdown_when_not_connected_enabled` - (Required) Will a VM get shutdown when it hasn't been connected to after a period of time? Possible values are `true` and `false`.

---

A `default_connection` block supports the following:

* `client_rdp_access` - (Required) The enabled access level for Client Access over RDP. Possible values are `None`, `Private` and `Public`.

* `client_ssh_access` - (Required) The enabled access level for Client Access over SSH. Possible values are `None`, `Private` and `Public`.

* `web_rdp_access` - (Required) The enabled access level for Web Access over RDP. Possible values are `None`, `Private` and `Public`.

* `web_ssh_access` - (Required) The enabled access level for Web Access over SSH. Possible values are `None`, `Private` and `Public`.

---

A `default_network` block supports the following:

* `subnet_id` - (Optional) The resource ID of the Subnet for the Lab Service Lab Plan network profile.

---

A `support` block supports the following:

* `email` - (Optional) The email address for the support contact.

* `instructions` - (Optional) The instructions for users of the Lab Service Lab Plan.

* `phone` - (Optional) The phone number for the support contact.

* `url` - (Optional) The web address for users of the Lab Service Lab Plan.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Lab Service Lab Plan.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Lab Service Lab Plan.
* `read` - (Defaults to 5 minutes) Used when retrieving the Lab Service Lab Plan.
* `update` - (Defaults to 30 minutes) Used when updating the Lab Service Lab Plan.
* `delete` - (Defaults to 30 minutes) Used when deleting the Lab Service Lab Plan.

## Import

Lab Service Lab Plans can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lab_service_lab_plan.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.LabServices/labPlans/labPlan1
```
