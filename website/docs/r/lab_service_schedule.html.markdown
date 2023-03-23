---
subcategory: "Lab Service"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lab_service_schedule"
description: |-
  Manages a Lab Service Schedule.
---

# azurerm_lab_service_schedule

Manages a Lab Service Schedule.

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
      capacity = 1
    }
  }
}

resource "azurerm_lab_service_schedule" "example" {
  name      = "example-labschedule"
  lab_id    = azurerm_lab_service_lab.example.id
  stop_time = "2022-11-28T00:00:00Z"
  time_zone = "America/Los_Angeles"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Lab Service Schedule. Changing this forces a new resource to be created.

* `lab_id` - (Required) The resource ID of the Lab Service Schedule. Changing this forces a new resource to be created.

* `stop_time` - (Required) When Lab User Virtual Machines will be stopped in RFC-3339 format.

* `time_zone` - (Required) The IANA Time Zone ID for the Schedule.

* `notes` - (Optional) The notes for the Schedule.

* `recurrence` - (Optional) A `recurrence` block as defined below.

* `start_time` - (Optional) When Lab User Virtual Machines will be started in RFC-3339 format.

---

A `recurrence` block supports the following:

* `expiration_date` - (Required) When the recurrence will expire in RFC-3339 format.

* `frequency` - (Required) The frequency of the recurrence. Possible values are `Daily` and `Weekly`.

* `interval` - (Optional) The interval to invoke the schedule on. Possible values are between `1` and `365`.

* `week_days` - (Optional) The interval to invoke the schedule on. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` and `Saturday`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Lab Service Schedule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Lab Service Schedule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Lab Service Schedule.
* `update` - (Defaults to 30 minutes) Used when updating the Lab Service Schedule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Lab Service Schedule.

## Import

Lab Service Schedules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lab_service_schedule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.LabServices/labs/lab1/schedules/schedule1
```
