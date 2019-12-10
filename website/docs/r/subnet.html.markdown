---
subcategory: "Network"
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
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_virtual_network" "example" {
  name                = "acceptanceTestVirtualNetwork1"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_subnet" "example" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.0.1.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action", "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action"]
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

* `enforce_private_link_service_network_policies` - (Optional) Enable or Disable network policies on the `private link service` in the subnet. Default is `false`.

-> **NOTE:** Network policies like network security groups (NSG) are not supported for the private link service. In order to deploy a private link service on a given subnet, an explicit disable setting is required on that subnet(e.g. `enforce_private_link_service_network_policies` = `true`). This setting is only applicable for the private link service. For other resources in the subnet, access is controlled based on Network Security Groups (NSG) security rules definition.

* `network_security_group_id` - (Optional / **Deprecated**) The ID of the Network Security Group to associate with the subnet.

-> **NOTE:** At this time Subnet `<->` Network Security Group associations need to be configured both using this field (which is now Deprecated) and using the `azurerm_subnet_network_security_group_association` resource. This field is deprecated and will be removed in favour of that resource in the next major version (2.0) of the AzureRM Provider.

* `route_table_id` - (Optional / **Deprecated**) The ID of the Route Table to associate with the subnet.

-> **NOTE:** At this time Subnet `<->` Route Table associations need to be configured both using this field (which is now Deprecated) and using the `azurerm_subnet_route_table_association` resource. This field is deprecated and will be removed in favour of that resource in the next major version (2.0) of the AzureRM Provider.

* `service_endpoints` - (Optional) The list of Service endpoints to associate with the subnet. Possible values include: `Microsoft.AzureActiveDirectory`, `Microsoft.AzureCosmosDB`, `Microsoft.ContainerRegistry`, `Microsoft.EventHub`, `Microsoft.KeyVault`, `Microsoft.ServiceBus`, `Microsoft.Sql`, `Microsoft.Storage` and `Microsoft.Web`.

* `delegation` - (Optional) One or more `delegation` blocks as defined below.

* `enforce_private_link_endpoint_network_policies` - (Optional) Enable or Disable network policies for the private link endpoint on the subnet. Default valule is `false`. Conflicts with enforce_private_link_service_network_policies.

-> **NOTE:** Network policies, like network security groups (NSG), are not supported for Private Link Endpoints or Private Link Services. In order to deploy a Private Link Endpoint on a given subnet, you must set the `enforce_private_link_endpoint_network_policies` attribute to `true`. This setting is only applicable for the Private Link Endpoint, for all other resources in the subnet access is controlled based on the `network_security_group_id`.

* `enforce_private_link_service_network_policies` - (Optional) Enable or Disable network policies for the private link service on the subnet. Default valule is `false`. Conflicts with enforce_private_link_endpoint_network_policies.

-> **NOTE:** In order to deploy a Private Link Service on a given subnet, you must set the `enforce_private_link_service_network_policies` attribute to `true`. This setting is only applicable for the Private Link Service, for all other resources in the subnet access is controlled based on the `network_security_group_id`. 

---

A `delegation` block supports the following:

* `name` (Required) A name for this delegation.

* `service_delegation` (Required) A `service_delegation` block as defined below.

---

A `service_delegation` block supports the following:

-> **NOTE:** Delegating to services may not be available in all regions. Check that the service you are delegating to is available in your region using the [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/network/vnet/subnet?view=azure-cli-latest#az-network-vnet-subnet-list-available-delegations). Also, `actions` is specific to each service type. The exact list of `actions` needs to be retrieved using the aforementioned [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/network/vnet/subnet?view=azure-cli-latest#az-network-vnet-subnet-list-available-delegations).

* `name` - (Required) The name of service to delegate to. Possible values include `Microsoft.BareMetal/AzureVMware`, `Microsoft.BareMetal/CrayServers`, `Microsoft.Batch/batchAccounts`, `Microsoft.ContainerInstance/containerGroups`, `Microsoft.Databricks/workspaces`, `Microsoft.DBforPostgreSQL/serversv2`, `Microsoft.HardwareSecurityModules/dedicatedHSMs`, `Microsoft.Logic/integrationServiceEnvironments`, `Microsoft.Netapp/volumes`, `Microsoft.ServiceFabricMesh/networks`, `Microsoft.Sql/managedInstances`, `Microsoft.Sql/servers`, `Microsoft.StreamAnalytics/streamingJobs`, `Microsoft.Web/hostingEnvironments` and `Microsoft.Web/serverFarms`.

* `actions` - (Optional) A list of Actions which should be delegated. This list is specific to the service to delegate to. Possible values include `Microsoft.Network/networkinterfaces/*`, `Microsoft.Network/virtualNetworks/subnets/action`, `Microsoft.Network/virtualNetworks/subnets/join/action`, `Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action` and `Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action`.

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
terraform import azurerm_subnet.exampleSubnet /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1
```
