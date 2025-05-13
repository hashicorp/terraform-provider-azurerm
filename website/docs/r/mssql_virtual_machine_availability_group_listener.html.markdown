---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_virtual_machine_availability_group_listener"
description: |-
  Manages a Microsoft SQL Virtual Machine Availability Group Listener.
---

# azurerm_mssql_virtual_machine_availability_group_listener

Manages a Microsoft SQL Virtual Machine Availability Group Listener.

## Example Usage

```hcl
data "azurerm_subnet" "example" {
  name                 = "examplesubnet"
  virtual_network_name = "examplevnet"
  resource_group_name  = "example-resources"
}

data "azurerm_lb" "example" {
  name                = "example-lb"
  resource_group_name = "example-resources"
}

data "azurerm_virtual_machine" "example" {
  count = 2

  name                = "example-vm"
  resource_group_name = "example-resources"
}

resource "azurerm_mssql_virtual_machine_group" "example" {
  name                = "examplegroup"
  resource_group_name = "example-resources"
  location            = "West Europe"

  sql_image_offer = "SQL2017-WS2016"
  sql_image_sku   = "Developer"

  wsfc_domain_profile {
    fqdn                = "testdomain.com"
    cluster_subnet_type = "SingleSubnet"
  }
}

resource "azurerm_mssql_virtual_machine" "example" {
  count = 2

  virtual_machine_id           = data.azurerm_virtual_machine.example[count.index].id
  sql_license_type             = "PAYG"
  sql_virtual_machine_group_id = azurerm_mssql_virtual_machine_group.example.id

  wsfc_domain_credential {
    cluster_bootstrap_account_password = "P@ssw0rd1234!"
    cluster_operator_account_password  = "P@ssw0rd1234!"
    sql_service_account_password       = "P@ssw0rd1234!"
  }
}

resource "azurerm_mssql_virtual_machine_availability_group_listener" "example" {
  name                         = "listener1"
  availability_group_name      = "availabilitygroup1"
  port                         = 1433
  sql_virtual_machine_group_id = azurerm_mssql_virtual_machine_group.example.id

  load_balancer_configuration {
    load_balancer_id   = data.azurerm_lb.example.id
    private_ip_address = "10.0.2.11"
    probe_port         = 51572
    subnet_id          = data.azurerm_subnet.example.id

    sql_virtual_machine_ids = [
      azurerm_mssql_virtual_machine.example[0].id,
      azurerm_mssql_virtual_machine.example[1].id
    ]
  }

  replica {
    sql_virtual_machine_id = azurerm_mssql_virtual_machine.example[0].id
    role                   = "Primary"
    commit                 = "Synchronous_Commit"
    failover               = "Automatic"
    readable_secondary     = "All"
  }

  replica {
    sql_virtual_machine_id = azurerm_mssql_virtual_machine.example[1].id
    role                   = "Secondary"
    commit                 = "Asynchronous_Commit"
    failover               = "Manual"
    readable_secondary     = "No"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for the Microsoft SQL Virtual Machine Availability Group Listener. Changing this forces a new resource to be created.

* `sql_virtual_machine_group_id` - (Required) The ID of the SQL Virtual Machine Group to create the listener. Changing this forces a new resource to be created.

* `availability_group_name` - (Optional) The name of the Availability Group. Changing this forces a new resource to be created.

* `load_balancer_configuration` - (Optional) A `load_balancer_configuration` block as defined below. Changing this forces a new resource to be created.

~> **Note:** Either one of `load_balancer_configuration` or `multi_subnet_ip_configuration` must be specified.

* `multi_subnet_ip_configuration` - (Optional) One or more `multi_subnet_ip_configuration` blocks as defined below. Changing this forces a new resource to be created.

* `port` - (Optional) The port of the listener. Changing this forces a new resource to be created.

* `replica` - (Required) One or more `replica` blocks as defined below. Changing this forces a new resource to be created.

---

A `load_balancer_configuration` block supports the following:

* `load_balancer_id` - (Required) The ID of the Load Balancer. Changing this forces a new resource to be created.

* `private_ip_address` - (Required) The private IP Address of the listener. Changing this forces a new resource to be created.

* `probe_port` - (Required) The probe port of the listener. Changing this forces a new resource to be created.

* `sql_virtual_machine_ids` - (Required) Specifies a list of SQL Virtual Machine IDs. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet to create the listener. Changing this forces a new resource to be created.

~> **Note:** `sql_virtual_machine_ids` should match with the SQL Virtual Machines specified in `replica`.

---

A `multi_subnet_ip_configuration` block supports the following:

* `private_ip_address` - (Required) The private IP Address of the listener. Changing this forces a new resource to be created.

* `sql_virtual_machine_id` - (Required) The ID of the Sql Virtual Machine. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet to create the listener. Changing this forces a new resource to be created.

~> **Note:** `sql_virtual_machine_id` should match with the SQL Virtual Machines specified in `replica`.

---

A `replica` block supports the following:

* `commit` - (Required) The replica commit mode for the availability group. Possible values are `Synchronous_Commit` and `Asynchronous_Commit`. Changing this forces a new resource to be created.

* `failover_mode` - (Required) The replica failover mode for the availability group. Possible values are `Manual` and `Automatic`. Changing this forces a new resource to be created.

* `readable_secondary` - (Required) The replica readable secondary mode for the availability group. Possible values are `No`, `Read_Only` and `All`. Changing this forces a new resource to be created.

* `role` - (Required) The replica role for the availability group. Possible values are `Primary` and `Secondary`. Changing this forces a new resource to be created.

* `sql_virtual_machine_id` - (Required) The ID of the SQL Virtual Machine. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Microsoft SQL Virtual Machine Availability Group Listener.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Microsoft SQL Virtual Machine Availability Group Listener.
* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft SQL Virtual Machine Availability Group Listener.
* `delete` - (Defaults to 30 minutes) Used when deleting the Microsoft SQL Virtual Machine Availability Group Listener.

## Import

Microsoft SQL Virtual Machine Availability Group Listeners can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_virtual_machine_availability_group_listener.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/vmgroup1/availabilityGroupListeners/listener1
```
