---
subcategory: "Chaos Studio"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_chaos_studio_experiment"
description: |-
  Manages a Chaos Studio Experiment.
---

# azurerm_chaos_studio_experiment

Manages a Chaos Studio Experiment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "westeurope"
}

resource "azurerm_user_assigned_identity" "example" {
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  name                = "example"
}

resource "azurerm_virtual_network" "example" {
  name                = "example"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "example"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "example" {
  name                            = "example"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "example"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}

resource "azurerm_chaos_studio_target" "example" {
  location           = azurerm_resource_group.example.location
  target_resource_id = azurerm_linux_virtual_machine.example.id
  target_type        = "Microsoft-VirtualMachine"
}

resource "azurerm_chaos_studio_capability" "example" {
  chaos_studio_target_id = azurerm_chaos_studio_target.example.id
  capability_type        = "Shutdown-1.0"
}

resource "azurerm_chaos_studio_experiment" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type = "SystemAssigned"
  }

  selectors {
    name                    = "Selector1"
    chaos_studio_target_ids = [azurerm_chaos_studio_target.example.id]
  }

  steps {
    name = "example"
    branch {
      name = "example"
      actions {
        urn           = azurerm_chaos_studio_capability.example.urn
        selector_name = "Selector1"
        parameters = {
          abruptShutdown = "false"
        }
        action_type = "continuous"
        duration    = "PT10M"
      }
    }
  }
}


```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Chaos Studio Experiment should exist. Changing this forces a new Chaos Studio Experiment to be created.

* `name` - (Required) The name which should be used for this Chaos Studio Experiment. Changing this forces a new Chaos Studio Experiment to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Chaos Studio Experiment should exist. Changing this forces a new Chaos Studio Experiment to be created.

* `selectors` - (Required) One or more `selectors` blocks as defined below.

* `steps` - (Required) One or more `steps` blocks as defined below.

---

* `identity` - (Optional) A `identity` block as defined below.

---

A `actions` block supports the following:

* `action_type` - (Required) The type of action that should be added to the experiment. Possible values are `continuous`, `delay` and `discrete`. 

* `duration` - (Optional) An ISO8601 formatted string specifying the duration for a `delay` or `continuous` action.

* `parameters` - (Optional) A key-value map of additional parameters to configure the action. The values that are accepted by this depend on the `urn` i.e. the capability/fault that is applied. Possible parameter values can be found in this [documentation](https://learn.microsoft.com/azure/chaos-studio/chaos-studio-fault-library)

* `selector_name` - (Optional) The name of the Selector to which this action should apply to. This must be specified if the `action_type` is `continuous` or `discrete`.

* `urn` - (Optional) The Unique Resource Name of the action, this value is provided by the `azurerm_chaos_studio_capability` resource e.g. `azurerm_chaos_studio_capability.example.urn`. This must be specified if the `action_type` is `continuous` or `discrete`.

---

A `branch` block supports the following:

* `actions` - (Required) One or more `actions` blocks as defined above.

* `name` - (Required) The name of the branch.

---

A `identity` block supports the following:

* `type` - (Required) The Type of Managed Identity which should be added to this Policy Definition. Possible values are `SystemAssigned` and `UserAssigned`.

* `identity_ids` - (Optional) A list of User Managed Identity IDs which should be assigned to the Policy Definition.

~> **Note:** This is required when `type` is set to `UserAssigned`.

---

A `selectors` block supports the following:

* `chaos_studio_target_ids` - (Required) A list of Chaos Studio Target IDs that should be part of this Selector.

* `name` - (Required) The name of this Selector.

---

A `steps` block supports the following:

* `branch` - (Required) One or more `branch` blocks as defined above.

* `name` - (Required) The name of the Step.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Chaos Studio Experiment.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Chaos Studio Experiment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Chaos Studio Experiment.
* `update` - (Defaults to 30 minutes) Used when updating the Chaos Studio Experiment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Chaos Studio Experiment.

## Import

Chaos Studio Experiments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_chaos_studio_experiment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Chaos/experiments/experiment1
```
