---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_assessment"
description: |-
  Manages the Security Center Assessment for Azure Security Center.
---

# azurerm_security_center_assessment

Manages the Security Center Assessment for Azure Security Center.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "internal" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_linux_virtual_machine_scale_set" "example" {
  name                = "example-vmss"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"

  admin_ssh_key {
    username   = "adminuser"
    public_key = file("~/.ssh/id_rsa.pub")
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.internal.id
    }
  }
}

resource "azurerm_security_center_assessment_policy" "example" {
  display_name = "Test Display Name"
  severity     = "Medium"
  description  = "Test Description"
}

resource "azurerm_security_center_assessment" "example" {
  assessment_policy_id = azurerm_security_center_assessment_policy.example.id
  target_resource_id   = azurerm_linux_virtual_machine_scale_set.example.id

  status {
    code = "Healthy"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `assessment_policy_id` - (Required) The ID of the security Assessment policy to apply to this resource. Changing this forces a new security Assessment to be created.

* `target_resource_id` - (Required) The ID of the target resource. Changing this forces a new security Assessment to be created.

* `status` - (Required) A `status` block as defined below.

* `additional_data` - (Optional) A map of additional data to associate with the assessment.

---

A `status` block supports the following:

* `code` - (Required) Specifies the programmatic code of the assessment status. Possible values are `Healthy`, `Unhealthy` and `NotApplicable`.

* `cause` - (Optional) Specifies the cause of the assessment status.

* `description` - (Optional) Specifies the human readable description of the assessment status.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Security Center Assessment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Security Center Assessment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Center Assessment.
* `update` - (Defaults to 30 minutes) Used when updating the Security Center Assessment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Security Center Assessment.

## Import

Security Assessment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_assessment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/providers/Microsoft.Security/assessments/00000000-0000-0000-0000-000000000000
```
