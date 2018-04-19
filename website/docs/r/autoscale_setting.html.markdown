---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_autoscale_setting"
sidebar_current: "docs-azurerm-resource-autoscale-setting"
description: |-
  Create an autoscale setting that can be applied to Virtual Machine Scale Sets, Cloud Services and App Service - Web Apps.
---

# azurerm\_autoscale\_setting

Create an [autoscale setting](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-overview-autoscale) that can be 
applied to [Virtual Machine Scale Sets](virtual_machine_scale_set.html), Cloud Services and App Service - Web Apps.

## Example Usage with Default Profile

```hcl
variable "adminPwd" {
  type    = "string"
  default = "^l&f_ZwX2L4C"
}

resource "azurerm_resource_group" "test" {
  name     = "autoscalingTest"
  location = "West US"
}

resource "azurerm_autoscale_setting" "test" {
  name                = "myAutoscaleSetting"
  enabled             = true
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
  profile {
    name = "defaultProfile"
    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }
    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "GreaterThan"
        threshold           = 75
      }
      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }
    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "LessThan"
        threshold           = 25
      }
      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }
  }

  notification {
    operation = "Scale"
    email {
      send_to_subscription_administrator    = true
      send_to_subscription_co_administrator = true
    }
  }
}

resource "azurerm_virtual_network" "test" { 
  name                = "acctvn" 
  address_space       = ["10.0.0.0/16"] 
  location            = "${azurerm_resource_group.test.location}" 
  resource_group_name = "${azurerm_resource_group.test.name}" 
} 
 
resource "azurerm_subnet" "test" { 
  name                 = "acctsub" 
  resource_group_name  = "${azurerm_resource_group.test.name}" 
  virtual_network_name = "${azurerm_virtual_network.test.name}" 
  address_prefix       = "10.0.2.0/24" 
} 
 
resource "azurerm_public_ip" "test" { 
  name                         = "acctPublicIp" 
  location                     = "${azurerm_resource_group.test.location}" 
  resource_group_name          = "${azurerm_resource_group.test.name}" 
  public_ip_address_allocation = "static" 
} 
 
resource "azurerm_lb" "test" { 
  name                = "acctlb" 
  location            = "${azurerm_resource_group.test.location}" 
  resource_group_name = "${azurerm_resource_group.test.name}" 

  frontend_ip_configuration { 
    name                 = "PublicIPAddress" 
    public_ip_address_id = "${azurerm_public_ip.test.id}" 
  } 
} 
 
resource "azurerm_lb_backend_address_pool" "bpepool" { 
  name                = "BackEndAddressPool" 
  resource_group_name = "${azurerm_resource_group.test.name}" 
  loadbalancer_id     = "${azurerm_lb.test.id}" 
} 
 
resource "azurerm_lb_nat_pool" "lbnatpool" { 
  count                          = 3 
  resource_group_name            = "${azurerm_resource_group.test.name}" 
  name                           = "ssh" 
  loadbalancer_id                = "${azurerm_lb.test.id}" 
  protocol                       = "Tcp" 
  frontend_port_start            = 50000 
  frontend_port_end              = 50119 
  backend_port                   = 22 
  frontend_ip_configuration_name = "PublicIPAddress" 
} 
 
resource "azurerm_virtual_machine_scale_set" "test" { 
  name                = "myTestScaleset-1" 
  location            = "${azurerm_resource_group.test.location}" 
  resource_group_name = "${azurerm_resource_group.test.name}" 
  upgrade_policy_mode = "Manual" 
 
  sku { 
    name     = "Standard_A0" 
    tier     = "Standard" 
    capacity = 2 
  } 
 
  storage_profile_image_reference { 
    publisher = "Canonical" 
    offer     = "UbuntuServer" 
    sku       = "14.04.2-LTS" 
    version   = "latest" 
  } 
 
  storage_profile_os_disk { 
    name              = "" 
    caching           = "ReadWrite" 
    create_option     = "FromImage" 
    managed_disk_type = "Standard_LRS" 
  } 
 
  storage_profile_data_disk { 
      lun          = 0 
    caching        = "ReadWrite" 
    create_option  = "Empty" 
    disk_size_gb   = 10  
  } 
 
  os_profile { 
    computer_name_prefix = "testvm" 
    admin_username       = "myadmin" 
    admin_password       = "${var.adminPwd}" 
  } 
 
  os_profile_linux_config { 
    disable_password_authentication = false 
  } 
 
  network_profile { 
    name    = "terraformnetworkprofile" 
    primary = true 
 
    ip_configuration { 
      name                                   = "TestIPConfiguration" 
      subnet_id                              = "${azurerm_subnet.test.id}" 
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.bpepool.id}"] 
      load_balancer_inbound_nat_rules_ids    = ["${element(azurerm_lb_nat_pool.lbnatpool.*.id, count.index)}"] 
    } 
  } 
 
  tags { 
    environment = "staging" 
  } 
}
```

## Example Usage with A Profile Repeats on Saturdays and Sundays

```hcl
variable "adminPwd" {
  type    = "string"
  default = "^l&f_ZwX2L4C"
}

resource "azurerm_resource_group" "test" {
  name     = "autoscalingTest"
  location = "West US"
}

resource "azurerm_autoscale_setting" "test" {
  name                = "myAutoscaleSetting"
  enabled             = true
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
  profile {
    name = "{\"name\":\"defaultProfile\",\"for\":\"Weekends\"}"
    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }
    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "GreaterThan"
        threshold           = 80
      }
      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }
    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "LessThan"
        threshold           = 20
      }
      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }
    recurrence {
      frequency = "Week"
      schedule {
        time_zone = "Pacific Standard Time"
        days      = ["Saturday", "Sunday"]
        hours     = [13]
        minutes   = [0]
      }
    }
  }
  profile {
    name = "Weekends"
    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }
    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "GreaterThan"
        threshold           = 90
      }
      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = "2"
        cooldown  = "PT1M"
      }
    }
    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "LessThan"
        threshold           = 10
      }
      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = "2"
        cooldown  = "PT1M"
      }
    }
    recurrence {
      frequency = "Week"
      schedule {
        time_zone = "Pacific Standard Time"
        days      = ["Saturday", "Sunday"]
        hours     = [12]
        minutes   = [0]
      }
    }
  }

  notification {
    operation = "Scale"
    email {
      send_to_subscription_administrator    = true
      send_to_subscription_co_administrator = true
    }
  }
}

resource "azurerm_virtual_network" "test" { 
  name                = "acctvn" 
  address_space       = ["10.0.0.0/16"] 
  location            = "${azurerm_resource_group.test.location}" 
  resource_group_name = "${azurerm_resource_group.test.name}" 
} 
 
resource "azurerm_subnet" "test" { 
  name                 = "acctsub" 
  resource_group_name  = "${azurerm_resource_group.test.name}" 
  virtual_network_name = "${azurerm_virtual_network.test.name}" 
  address_prefix       = "10.0.2.0/24" 
} 
 
resource "azurerm_public_ip" "test" { 
  name                         = "acctPublicIp" 
  location                     = "${azurerm_resource_group.test.location}" 
  resource_group_name          = "${azurerm_resource_group.test.name}" 
  public_ip_address_allocation = "static" 
} 
 
resource "azurerm_lb" "test" { 
  name                = "acctlb" 
  location            = "${azurerm_resource_group.test.location}" 
  resource_group_name = "${azurerm_resource_group.test.name}" 

  frontend_ip_configuration { 
    name                 = "PublicIPAddress" 
    public_ip_address_id = "${azurerm_public_ip.test.id}" 
  } 
} 
 
resource "azurerm_lb_backend_address_pool" "bpepool" { 
  name                = "BackEndAddressPool" 
  resource_group_name = "${azurerm_resource_group.test.name}" 
  loadbalancer_id     = "${azurerm_lb.test.id}" 
} 
 
resource "azurerm_lb_nat_pool" "lbnatpool" { 
  count                          = 3 
  resource_group_name            = "${azurerm_resource_group.test.name}" 
  name                           = "ssh" 
  loadbalancer_id                = "${azurerm_lb.test.id}" 
  protocol                       = "Tcp" 
  frontend_port_start            = 50000 
  frontend_port_end              = 50119 
  backend_port                   = 22 
  frontend_ip_configuration_name = "PublicIPAddress" 
} 
 
resource "azurerm_virtual_machine_scale_set" "test" { 
  name                = "myTestScaleset-1" 
  location            = "${azurerm_resource_group.test.location}" 
  resource_group_name = "${azurerm_resource_group.test.name}" 
  upgrade_policy_mode = "Manual" 
 
  sku { 
    name     = "Standard_A0" 
    tier     = "Standard" 
    capacity = 2 
  } 
 
  storage_profile_image_reference { 
    publisher = "Canonical" 
    offer     = "UbuntuServer" 
    sku       = "14.04.2-LTS" 
    version   = "latest" 
  } 
 
  storage_profile_os_disk { 
    name              = "" 
    caching           = "ReadWrite" 
    create_option     = "FromImage" 
    managed_disk_type = "Standard_LRS" 
  } 
 
  storage_profile_data_disk { 
      lun          = 0 
    caching        = "ReadWrite" 
    create_option  = "Empty" 
    disk_size_gb   = 10  
  } 
 
  os_profile { 
    computer_name_prefix = "testvm" 
    admin_username       = "myadmin" 
    admin_password       = "${var.adminPwd}" 
  } 
 
  os_profile_linux_config { 
    disable_password_authentication = false 
  } 
 
  network_profile { 
    name    = "terraformnetworkprofile" 
    primary = true 
 
    ip_configuration { 
      name                                   = "TestIPConfiguration" 
      subnet_id                              = "${azurerm_subnet.test.id}" 
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.bpepool.id}"] 
      load_balancer_inbound_nat_rules_ids    = ["${element(azurerm_lb_nat_pool.lbnatpool.*.id, count.index)}"] 
    } 
  } 
 
  tags { 
    environment = "staging" 
  } 
}
```

## Example Usage with Specific Start and End Dates

```hcl
variable "adminPwd" {
  type    = "string"
  default = "^l&f_ZwX2L4C"
}

resource "azurerm_resource_group" "test" {
  name     = "autoscalingTest"
  location = "West US"
}

resource "azurerm_autoscale_setting" "test" {
  name                = "myAutoscaleSetting"
  enabled             = true
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
  profile {
    name = "defaultProfile"
    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }
    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "GreaterThan"
        threshold           = 80
      }
      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }
    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "LessThan"
        threshold           = 20
      }
      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }
  }
  profile {
    name = "forJuly"
    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }
    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "GreaterThan"
        threshold           = 90
      }
      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = "2"
        cooldown  = "PT1M"
      }
    }
    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "LessThan"
        threshold           = 10
      }
      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = "2"
        cooldown  = "PT1M"
      }
    }
    fixed_date {
      time_zone = "Pacific Standard Time"
      start     = "2017-07-01T00:00:00Z"
      end       = "2017-07-31T23:59:59Z"
    }
  }

  notification {
    operation = "Scale"
    email {
      send_to_subscription_administrator    = true
      send_to_subscription_co_administrator = true
    }
  }
}

resource "azurerm_virtual_network" "test" { 
  name                = "acctvn" 
  address_space       = ["10.0.0.0/16"] 
  location            = "${azurerm_resource_group.test.location}" 
  resource_group_name = "${azurerm_resource_group.test.name}" 
} 
 
resource "azurerm_subnet" "test" { 
  name                 = "acctsub" 
  resource_group_name  = "${azurerm_resource_group.test.name}" 
  virtual_network_name = "${azurerm_virtual_network.test.name}" 
  address_prefix       = "10.0.2.0/24" 
} 
 
resource "azurerm_public_ip" "test" { 
  name                         = "acctPublicIp" 
  location                     = "${azurerm_resource_group.test.location}" 
  resource_group_name          = "${azurerm_resource_group.test.name}" 
  public_ip_address_allocation = "static" 
} 
 
resource "azurerm_lb" "test" { 
  name                = "acctlb" 
  location            = "${azurerm_resource_group.test.location}" 
  resource_group_name = "${azurerm_resource_group.test.name}" 

  frontend_ip_configuration { 
    name                 = "PublicIPAddress" 
    public_ip_address_id = "${azurerm_public_ip.test.id}" 
  } 
} 
 
resource "azurerm_lb_backend_address_pool" "bpepool" { 
  name                = "BackEndAddressPool" 
  resource_group_name = "${azurerm_resource_group.test.name}" 
  loadbalancer_id     = "${azurerm_lb.test.id}" 
} 
 
resource "azurerm_lb_nat_pool" "lbnatpool" { 
  count                          = 3 
  resource_group_name            = "${azurerm_resource_group.test.name}" 
  name                           = "ssh" 
  loadbalancer_id                = "${azurerm_lb.test.id}" 
  protocol                       = "Tcp" 
  frontend_port_start            = 50000 
  frontend_port_end              = 50119 
  backend_port                   = 22 
  frontend_ip_configuration_name = "PublicIPAddress" 
} 
 
resource "azurerm_virtual_machine_scale_set" "test" { 
  name                = "myTestScaleset-1" 
  location            = "${azurerm_resource_group.test.location}" 
  resource_group_name = "${azurerm_resource_group.test.name}" 
  upgrade_policy_mode = "Manual" 
 
  sku { 
    name     = "Standard_A0" 
    tier     = "Standard" 
    capacity = 2 
  } 
 
  storage_profile_image_reference { 
    publisher = "Canonical" 
    offer     = "UbuntuServer" 
    sku       = "14.04.2-LTS" 
    version   = "latest" 
  } 
 
  storage_profile_os_disk { 
    name              = "" 
    caching           = "ReadWrite" 
    create_option     = "FromImage" 
    managed_disk_type = "Standard_LRS" 
  } 
 
  storage_profile_data_disk { 
      lun          = 0 
    caching        = "ReadWrite" 
    create_option  = "Empty" 
    disk_size_gb   = 10  
  } 
 
  os_profile { 
    computer_name_prefix = "testvm" 
    admin_username       = "myadmin" 
    admin_password       = "${var.adminPwd}" 
  } 
 
  os_profile_linux_config { 
    disable_password_authentication = false 
  } 
 
  network_profile { 
    name    = "terraformnetworkprofile" 
    primary = true 
 
    ip_configuration { 
      name                                   = "TestIPConfiguration" 
      subnet_id                              = "${azurerm_subnet.test.id}" 
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.bpepool.id}"] 
      load_balancer_inbound_nat_rules_ids    = ["${element(azurerm_lb_nat_pool.lbnatpool.*.id, count.index)}"] 
    } 
  } 
 
  tags { 
    environment = "staging" 
  } 
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the autoscale setting. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the autoscale setting. Changing this 
forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `enabled` - (Required) Specifies whether automatic scaling is enabled for the target resource.

* `target_resource_id` - (Required) Specifies the resource ID of the resource that the autoscale setting should be added to.

* `profile` - (Required) Specifies a collection of automatic scaling profiles that specify different scaling parameters for 
different time periods. A maximum of 20 profiles can be specified.

* `notification` - (Optional) Specifies notification settings for autoscale events.

* `tags` - (Optional) A mapping of tags to assign to the resource.

`profile` supports the following:

* `name` - (Required) Specifies the name of the profile.

* `capacity` - (Required) Specifies the number of instances that can be used during this profile.

* `rule` - (Required) Specifies a collection of rules that provide the triggers and parameters for the scaling action. 
A maximum of 10 rules can be specified.

* `fixed_date` - (Optional) Specifies a date for the profile. This element is not used if the `recurrence` element is used.

* `recurrence` - (Optional) Specifies the repeating times at which this profile begins. This element is not used if the `fixed_date` element is used.

`capacity` supports the following:

* `minimum` - (Required) Specifies the minimum number of instances that are available for the scaling action.

* `maximum` - (Required) Specifies the maximum number of instances that are available for the scaling action. 
The maximum number of instances is limited by the cores that are available in the subscription. You can use the portal to see 
the number of cores that are available in your subscription.

* `default` - (Required) Specifies the number of instances that are available for scaling if metrics are not available for 
evaluation. The default is only used if the current instance count is lower than the default.

`rule` supports the following:

* `metric_trigger` - (Required) Specifies the trigger that results in a scaling action.

* `scale_action` - (Required) Specifies parameters for the scaling action.

`metric_trigger` supports the following:

* `metric_name` - (Required) Specifies the name of the metric that defines what the rule monitors.

* `metric_resource_id` - (Required) Specifies the resource identifier of the resource the rule monitors.

* `time_grain` - (Required) Specifies the granularity of metrics the rule monitors. Must be one of the predefined 
values returned from metric definitions for the metric. Must be between 1 minute and 12 hours. ISO 8601 duration format.

* `statistic` - (Required) Specifies how the metrics from multiple instances are combined. Possible values are: 
`Average`, `Min`, `Max`.

* `time_window` - (Required) Specifies the range of time in which instance data is collected. This value must be 
greater than the delay in metric collection, which can vary from resource-to-resource. Must be between 5 minutes 
and 12 hours. ISO 8601 duration format.

* `time_aggregation` - (Required) Specifies how the data that is collected should be combined over time. The default 
value is Average. Possible values are: `Average`, `Minimum`, `Maximum`, `Last`, `Total`, `Count`.

* `operator` - (Required) Specifies the operator that is used to compare the metric data and the threshold. Possible 
values are: `Equals`, `NotEquals`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan`, `LessThanOrEqual`.

* `threshold` - (Required) Specifies the threshold of the metric that triggers the scale action.

`scale_action` supports the following:

* `direction` - (Required) Specifies the scale direction. Possible values are: `Increase`, `Decrease`.

* `type` - (Required) Specifies the type of action that should occur, this must be set to `ChangeCount` or `PercentChangeCount`.

* `value` - (Required) Specifies the number that is involved in the scaling action. This value must be 1 or greater. The default 
value is 1.

* `cooldown` - (Required) Specifies the amount of time to wait since the last scaling action before this action occurs. Must 
be between 1 minute and 1 week. ISO 8601 duration format.

`fixed_date` supports the following:

* `time_zone` - (Required) Specifies the time zone of the start and end times for the profile. Click [here](https://msdn.microsoft.com/en-us/library/azure/dn931928.aspx) for the complete list of possible values.

* `start` - (Required) Specifies the start time for the profile. RFC3339 format.

* `end` - (Required) Specifies the end time for the profile. RFC3339 format.

`recurrence` supports the following:

* `frequency` - (Required) Specifies how often the schedule profile should take effect. This value must be `Week`, meaning each week 
will have the same set profile.

* `schedule` - (Required) Specifies the scheduling constraints for when the profile begins.

`schedule` supports the following:

* `time_zone` - (Required) Specifies the time zone for the hours of the profile. Click [here](https://msdn.microsoft.com/en-us/library/azure/dn931928.aspx) for the complete list of possible values.

* `days` - (Required) Specifies a list of days that the profile takes effect on. Possible values are Sunday through Saturday.

* `hours` - (Required) Specifies a list of hours that the profile takes effect on. Values supported are 0 to 23 on the 24-hour 
clock (AM/PM times are not supported).

* `minutes` - (Required) Specifies a list of minutes at which the profile takes effect at.

`notification` supports the following:

* `operation` - (Required) Specifies the type of operation which triggers the notification. Possible value is `Scale`.

* `email` - (Required) Specifies the email notification settings.

* `webhook` - (Optional) Specifies the webhook notification settings.

`email` supports the following:

* `send_to_subscription_administrator` - (Required) Specifies whether to send email notifications to the subscription administrator.

* `send_to_subscription_co_administrator` - (Required) Specifies whether to send email notifications to the subscription co-administrator.

* `custom_emails` - (Optional) Specifies a list of custom email addresses to which the email notifications will be sent.

`webhook` supports the following:

* `service_uri` - (Required) Specifies a valid HTTPS URI for the webhook call.

* `properties` - (Optional) Specifies properties in key-value pairs.

## Attributes Reference

The following attributes are exported:

* `id` - The autoscale setting resource ID.


## Import

Autoscale settings can be imported using the `resource id`, e.g.

```
terraform import azurerm_autoscale_setting.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/autoscalingTest/providers/microsoft.insights/autoscalesettings/myAutoscaleSetting
```
