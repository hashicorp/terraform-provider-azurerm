---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet"
sidebar_current: "docs-azurerm-resource-network-subnet-x"
description: |-
  Manages a subnet. Subnets represent network segments within the IP space defined by the virtual network.

---

# azurerm_subnet

Manages a subnet. Subnets represent network segments within the IP space defined by the virtual network.

~> **NOTE on Virtual Networks and Subnet's:** Terraform currently
provides both a standalone [Subnet resource](subnet.html), and allows for Subnets to be defined in-line within the [Virtual Network resource](virtual_network.html).
At this time you cannot use a Virtual Network with in-line Subnets in conjunction with any Subnet resources. Doing so will cause a conflict of Subnet configurations and will overwrite Subnet's.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_virtual_network" "test" {
  name                = "acceptanceTestVirtualNetwork1"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the subnet. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the subnet. Changing this forces a new resource to be created.

* `virtual_network_name` - (Required) The name of the virtual network to which to attach the subnet. Changing this forces a new resource to be created.

* `address_prefix` - (Required) The address prefix to use for the subnet.

* `network_security_group_id` - (Optional / **Deprecated**) The ID of the Network Security Group to associate with the subnet.

-> **NOTE:** At this time Subnet `<->` Network Security Group associations need to be configured both using this field (which is now Deprecated) and/or using the `azurerm_subnet_network_security_group_association` resource. This field is deprecated and will be removed in favour of that resource in the next major version (2.0) of the AzureRM Provider.

* `route_table_id` - (Optional / **Deprecated**) The ID of the Route Table to associate with the subnet.

-> **NOTE:** At this time Subnet `<->` Route Table associations need to be configured both using this field (which is now Deprecated) and/or using the `azurerm_subnet_route_table_association` resource. This field is deprecated and will be removed in favour of that resource in the next major version (2.0) of the AzureRM Provider.

* `service_endpoints` - (Optional) The list of Service endpoints to associate with the subnet. Possible values include: `Microsoft.AzureActiveDirectory`, `Microsoft.AzureCosmosDB`, `Microsoft.ContainerRegistry`, `Microsoft.EventHub`, `Microsoft.KeyVault`, `Microsoft.ServiceBus`, `Microsoft.Sql`, `Microsoft.Storage` and `Microsoft.Web`.

* `delegation` - (Optional) One or more `delegation` blocks as defined below.

---

A `delegation` block supports the following:
* `name` (Required) A name for this delegation.
* `service_delegation` (Required) A `service_delegation` block as defined below.

---

A `service_delegation` block supports the following:

-> **NOTE:** Delegating to services may not be available in all regions. Check that the service you are delegating to is available in your region using the [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/network/vnet/subnet?view=azure-cli-latest#az-network-vnet-subnet-list-available-delegations)

* `name` - (Required) The name of service to delegate to. Possible values include: `Microsoft.BareMetal/AzureVMware`, 
`Microsoft.BareMetal/CrayServers`, `Microsoft.Batch/batchAccounts`, `Microsoft.ContainerInstance/containerGroups`, `Microsoft.Databricks/workspaces`, `Microsoft.HardwareSecurityModules/dedicatedHSMs`, `Microsoft.Logic/integrationServiceEnvironments`, `Microsoft.Netapp/volumes`, `Microsoft.ServiceFabricMesh/networks`, `Microsoft.Sql/managedInstances`, `Microsoft.Sql/servers`, `Microsoft.Web/hostingEnvironments` or `Microsoft.Web/serverFarms`.
* `actions` - (Optional) A list of Actions which should be delegated. Possible values include: `Microsoft.Network/virtualNetworks/subnets/action`.

## Attributes Reference

The following attributes are exported:

* `id` - The subnet ID.
* `ip_configurations` - The collection of IP Configurations with IPs within this subnet.
* `name` - The name of the subnet.
* `resource_group_name` - The name of the resource group in which the subnet is created in.
* `virtual_network_name` - The name of the virtual network in which the subnet is created in
* `address_prefix` - The address prefix for the subnet

## Import

Subnets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_subnet.testSubnet /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1
```
