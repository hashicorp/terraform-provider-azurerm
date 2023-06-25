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

resource "azurerm_lab_service_lab" "example" {
  name                = "example-lab"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  title               = "Test Title"

  security {
    open_access_enabled = false
  }

  virtual_machine {
    admin_user {
      username = "testadmin"
      password = "Password1234!"
    }

    image_reference {
      offer     = "0001-com-ubuntu-server-focal"
      publisher = "canonical"
      sku       = "20_04-lts"
      version   = "latest"
    }

    sku {
      name     = "Classic_Fsv2_2_4GB_128_S_SSD"
      capacity = 0
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Lab Service Lab. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Lab Service Lab should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Lab Service Lab should exist. Changing this forces a new resource to be created.

* `security` - (Required) A `security` block as defined below.

* `title` - (Required) The title of the Lab Service Lab.

* `virtual_machine` - (Required) A `virtual_machine` block as defined below.

* `auto_shutdown` - (Optional) An `auto_shutdown` block as defined below.

* `connection_setting` - (Optional) A `connection_setting` block as defined below.

* `description` - (Optional) The description of the Lab Service Lab.

* `lab_plan_id` - (Optional) The resource ID of the Lab Plan that is used during resource creation to provide defaults and acts as a permission container when creating a Lab Service Lab via `labs.azure.com`.

* `network` - (Optional) A `network` block as defined below.

* `roster` - (Optional) A `roster` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Lab Service Lab.

---

An `auto_shutdown` block supports the following:

* `disconnect_delay` - (Optional) The amount of time a VM will stay running after a user disconnects if this behavior is enabled. This value must be formatted as an ISO 8601 string.

~> **NOTE:** The `shutdownOnDisconnect` is `Disabled` when `disconnect_delay` isn't specified.

* `idle_delay` - (Optional) The amount of time a VM will idle before it is shutdown if this behavior is enabled. This value must be formatted as an ISO 8601 string.

* `no_connect_delay` - (Optional) The amount of time a VM will stay running before it is shutdown if no connection is made and this behavior is enabled. This value must be formatted as an ISO 8601 string.

~> **NOTE:** The `shutdownWhenNotConnected` is `Disabled` when `no_connect_delay` isn't specified.

* `shutdown_on_idle` - (Optional) A VM will get shutdown when it has idled for a period of time. Possible values are `LowUsage` and `UserAbsence`.

~> **NOTE:** This property is `None` when it isn't specified. No need to set `idle_delay` when `shutdown_on_idle` isn't specified.

---

A `connection_setting` block supports the following:

* `client_rdp_access` - (Optional) The enabled access level for Client Access over RDP. Possible value is `Public`.

~> **NOTE:** This property is `None` when it isn't specified.

* `client_ssh_access` - (Optional) The enabled access level for Client Access over SSH. Possible value is `Public`.

~> **NOTE:** This property is `None` when it isn't specified.

---

A `security` block supports the following:

* `open_access_enabled` - (Required) Is open access enabled to allow any user or only specified users to register to a Lab Service Lab?

---

A `virtual_machine` block supports the following:

* `admin_user` - (Required) An `admin_user` block as defined below.

* `image_reference` - (Required) An `image_reference` block as defined below.

* `sku` - (Required) A `sku` block as defined below.

* `additional_capability_gpu_drivers_installed` - (Optional) Is flagged to pre-install dedicated GPU drivers? Defaults to `false`. Changing this forces a new resource to be created.

* `create_option` - (Optional) The create option to indicate what Lab Service Lab VMs are created from. Possible values are `Image` and `TemplateVM`. Defaults to `Image`. Changing this forces a new resource to be created.

* `non_admin_user` - (Optional) A `non_admin_user` block as defined below.

* `shared_password_enabled` - (Optional) Is the shared password enabled with the same password for all user VMs? Defaults to `false`. Changing this forces a new resource to be created.

* `usage_quota` - (Optional) The initial quota allocated to each Lab Service Lab user. Defaults to `PT0S`. This value must be formatted as an ISO 8601 string.

---

An `admin_user` block supports the following:

* `username` - (Required) The username to use when signing in to Lab Service Lab VMs. Changing this forces a new resource to be created.

* `password` - (Required) The password for the Lab user. Changing this forces a new resource to be created.

---

An `image_reference` block supports the following:

* `id` - (Optional) The resource ID of the image. Changing this forces a new resource to be created.

* `offer` - (Optional) The image offer if applicable. Changing this forces a new resource to be created.

* `publisher` - (Optional) The image publisher. Changing this forces a new resource to be created.

* `sku` - (Optional) The image SKU. Changing this forces a new resource to be created.

* `version` - (Optional) The image version specified on creation. Changing this forces a new resource to be created.

---

A `non_admin_user` block supports the following:

* `username` - (Required) The username to use when signing in to Lab Service Lab VMs.

* `password` - (Required) The password for the user.

---

A `sku` block supports the following:

* `name` - (Required) The name of the SKU. Changing this forces a new resource to be created.

* `capacity` - (Required) The capacity for the SKU. Possible values are between `0` and `400`.

~> **NOTE:** Once `active_directory_group_id` is enabled, `capacity` wouldn't take effect, and it would be automatically set to the number of members in AAD Group by service API. So it has to use `ignore_changes` to avoid the difference of tf plan.

---

A `network` block supports the following:

* `subnet_id` - (Optional) The resource ID of the Subnet for the network profile of the Lab Service Lab.

---

A `roster` block supports the following:

* `active_directory_group_id` - (Optional) The AAD group ID which this Lab Service Lab roster is populated from.

* `lms_instance` - (Optional) The base URI identifying the lms instance.

* `lti_client_id` - (Optional) The unique id of the Azure Lab Service tool in the lms.

* `lti_context_id` - (Optional) The unique context identifier for the Lab Service Lab in the lms.

* `lti_roster_endpoint` - (Optional) The URI of the names and roles service endpoint on the lms for the class attached to this Lab Service Lab.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Lab Service Lab.

* `security` - A `security` block as defined below.

* `network` - A `network` block as defined below.

---

A `security` block supports the following:

* `registration_code` - The registration code for the Lab Service Lab.

---

A `network` block supports the following:

* `load_balancer_id` - The resource ID of the Load Balancer for the network profile of the Lab Service Lab.

* `public_ip_id` - The resource ID of the Public IP for the network profile of the Lab Service Lab.

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
