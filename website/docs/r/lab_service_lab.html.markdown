---
subcategory: "Lab Service"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lab_service_lab"
description: |-
  Manages a Lab Service Lab.
---

# azurerm_lab_service_lab

Manages a Lab Service Lab.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_lab_services_lab" "example" {
  name                = "example-lab"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  title               = "Test Title"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Lab Service Lab. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Lab Service Lab should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Lab Service Lab should exist. Changing this forces a new resource to be created.

* `title` - (Required) The title of the Lab Service Lab.

* `auto_shutdown_profile` - (Required) An `auto_shutdown_profile` block as defined below.

* `connection_profile` - (Required) A `connection_profile` block as defined below.

* `security_profile` - (Required) A `security_profile` block as defined below.

* `virtual_machine_profile` - (Required) A `virtual_machine_profile` block as defined below.

* `description` - (Optional) The description of the Lab Service Lab.

* `lab_plan_id` - (Optional) The resource ID of the Lab Plan that is used during resource creation to provide defaults and acts as a permission container when creating a Lab Service Lab via `labs.azure.com`.

* `network_profile` - (Optional) A `network_profile` block as defined below.

* `roster_profile` - (Optional) A `roster_profile` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Lab Service Lab.

---

An `auto_shutdown_profile` block supports the following:

* `disconnect_delay` - (Required) The amount of time a VM will stay running after a user disconnects if this behavior is enabled. This value must be formatted as an ISO 8601 string.

* `idle_delay` - (Required) The amount of time a VM will idle before it is shutdown if this behavior is enabled. This value must be formatted as an ISO 8601 string.

* `no_connect_delay` - (Required) The amount of time a VM will stay running before it is shutdown if no connection is made and this behavior is enabled. This value must be formatted as an ISO 8601 string.

* `shutdown_on_disconnect_enabled` - (Required) Is shutdown on disconnect is enabled?

* `shutdown_on_idle` - (Required) A VM will get shutdown when it has idled for a period of time. Possible values are `LowUsage`, `None` and `UserAbsence`.

* `shutdown_enabled_when_not_connected` - (Required) Is a VM shutdown when it hasn't been connected to after a period of time?

---

A `connection_profile` block supports the following:

* `client_rdp_access` - (Required) The enabled access level for Client Access over RDP. Possible values are `None`, `Private` and `Public`.

* `client_ssh_access` - (Required) The enabled access level for Client Access over SSH. Possible values are `None`, `Private` and `Public`.

* `web_rdp_access` - (Required) The enabled access level for Web Access over RDP. Possible values are `None`, `Private` and `Public`.

* `web_ssh_access` - (Required) The enabled access level for Web Access over SSH. Possible values are `None`, `Private` and `Public`.

---

A `security_profile` block supports the following:

* `open_access` - (Required) Is open access enabled to allow any user or only specified users to register to a Lab Service Lab?

---

A `virtual_machine_profile` block supports the following:

* `admin_user` - (Required) An `admin_user` block as defined below.

* `create_option` - (Required) The create option to indicate what Lab Service Lab VMs are created from. Possible values are `Image` and `TemplateVM`. Changing this forces a new resource to be created.

* `image_reference` - (Required) An `image_reference` block as defined below.

* `os_type` - (Optional) .

* `sku` - (Required) A `sku` block as defined below.

* `usage_quota` - (Required) The initial quota allocated to each Lab Service Lab user. This value must be formatted as an ISO 8601 string.

* `shared_password_enabled` - (Required) Is the shared password enabled with the same password for all user VMs? Changing this forces a new resource to be created.

* `additional_capability` - (Optional) An `additional_capability` block as defined below.

* `non_admin_user` - (Optional) A `non_admin_user` block as defined below.

---

An `additional_capability` block supports the following:

* `gpu_drivers_installed` - (Optional) Is flagged to pre-install dedicated GPU drivers? Defaults to `false`. Changing this forces a new resource to be created.

---

An `admin_user` block supports the following:

* `username` - (Required) The username to use when signing in to Lab Service Lab VMs. Changing this forces a new resource to be created.

* `password` - (Optional) The password for the Lab user. Changing this forces a new resource to be created.

---

An `image_reference` block supports the following:

* `id` - (Optional) The resource ID of the image. Changing this forces a new resource to be created.

* `offer` - (Optional) The image offer if applicable. Changing this forces a new resource to be created.

* `publisher` - (Optional) The image publisher. Changing this forces a new resource to be created.

* `sku` - (Optional) The image SKU. Changing this forces a new resource to be created.

* `version` - (Optional) The image version specified on creation. Changing this forces a new resource to be created.

---

A `non_admin_user` block supports the following:

* `username` - (Required) The username to use when signing in to Lab Service Lab VMs. Changing this forces a new resource to be created.

* `password` - (Optional) The password for the user. Changing this forces a new resource to be created.

---

A `sku` block supports the following:

* `name` - (Required) The name of the SKU. Changing this forces a new resource to be created.

* `capacity` - (Required) The capacity for the SKU. Possible values are between `0` and `400`.

---

A `network_profile` block supports the following:

* `subnet_id` - (Optional) The resource ID of the Subnet for the network profile of the Lab Service Lab. Changing this forces a new resource to be created.

---

A `roster_profile` block supports the following:

* `active_directory_group_id` - (Optional) The AAD group ID which this Lab Service Lab roster is populated from.

* `lms_instance` - (Optional) The base URI identifying the lms instance.

* `lti_client_id` - (Optional) The unique id of the Azure Lab Service tool in the lms.

* `lti_context_id` - (Optional) The unique context identifier for the Lab Service Lab in the lms.

* `lti_roster_endpoint` - (Optional) The URI of the names and roles service endpoint on the lms for the class attached to this Lab Service Lab.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Lab Service Lab.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Lab Service Lab.
* `read` - (Defaults to 5 minutes) Used when retrieving the Lab Service Lab.
* `update` - (Defaults to 90 minutes) Used when updating the Lab Service Lab.
* `delete` - (Defaults to 90 minutes) Used when deleting the Lab Service Lab.

## Import

Lab Service Labs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lab_service_lab.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.LabServices/labs/lab1
```
