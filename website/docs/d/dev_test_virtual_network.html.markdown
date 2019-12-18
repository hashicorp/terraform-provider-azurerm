---
subcategory: "Dev Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_test_virtual_network"
sidebar_current: "docs-azurerm-datasource-dev-test-virtual-network"
description: |-
  Gets information about an existing Dev Test Lab Virtual Network.
---

# Data Source: azurerm_dev_test_virtual_network

Use this data source to access information about an existing Dev Test Lab Virtual Network.

## Example Usage

```hcl
data "azurerm_dev_test_virtual_network" "example" {
  name                = "example-network"
  lab_name            = "examplelab"
  resource_group_name = "example-resource"
}

output "lab_subnet_name" {
  value = "${data.azurerm_dev_test_virtual_network.example.allowed_subnets.0.lab_subnet_name}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Virtual Network.
* `lab_name` - (Required) Specifies the name of the Dev Test Lab.
* `resource_group_name` - (Required) Specifies the name of the resource group that contains the Virtual Network.

## Attributes Reference

* `allowed_subnets` - The list of subnets enabled for the virtual network as defined below.
* `subnet_overrides` - The list of permission overrides for the subnets as defined below.
* `unique_identifier` - The unique immutable identifier of the virtual network.

---

An `allowed_subnets` block supports the following:

* `allow_public_ip` - Indicates if this subnet allows public IP addresses. Possible values are `Allow`, `Default` and `Deny`.

* `lab_subnet_name` - The name of the subnet.

* `resource_id` - The resource identifier for the subnet.

---

An `subnets_override` block supports the following:

* `lab_subnet_name` - The name of the subnet.

* `resource_id` - The resource identifier for the subnet.

* `use_in_vm_creation_permission` - Indicates if the subnet can be used for VM creation.  Possible values are `Allow`, `Default` and `Deny`.

* `use_public_ip_permission` - Indicates if the subnet can be assigned public IP addresses.  Possible values are `Allow`, `Default` and `Deny`.

* `virtual_network_pool_name` - The virtual network pool associated with this subnet.
