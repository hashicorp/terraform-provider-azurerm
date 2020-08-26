---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet"
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
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

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

* `address_prefix` - (Optional / **Deprecated in favour of `address_prefixes`**) The address prefix to use for the subnet.

* `address_prefixes` - (Optional) The address prefixes to use for the subnet.

-> **NOTE:** One of `address_prefix` or `address_prefixes` is required.

---

* `delegation` - (Optional) One or more `delegation` blocks as defined below.

* `enforce_private_link_endpoint_network_policies` - (Optional) Enable or Disable network policies for the private link endpoint on the subnet. Default value is `false`. Conflicts with enforce_private_link_service_network_policies.

-> **NOTE:** Network policies, like network security groups (NSG), are not supported for Private Link Endpoints or Private Link Services. In order to deploy a Private Link Endpoint on a given subnet, you must set the `enforce_private_link_endpoint_network_policies` attribute to `true`. This setting is only applicable for the Private Link Endpoint, for all other resources in the subnet access is controlled based via the Network Security Group which can be configured using the `azurerm_subnet_network_security_group_association` resource.

* `enforce_private_link_service_network_policies` - (Optional) Enable or Disable network policies for the private link service on the subnet. Default valule is `false`. Conflicts with `enforce_private_link_endpoint_network_policies`.

-> **NOTE:** In order to deploy a Private Link Service on a given subnet, you must set the `enforce_private_link_service_network_policies` attribute to `true`. This setting is only applicable for the Private Link Service, for all other resources in the subnet access is controlled based on the Network Security Group which can be configured using the `azurerm_subnet_network_security_group_association` resource.

* `service_endpoints` - (Optional) The list of Service endpoints to associate with the subnet. Possible values include: `Microsoft.AzureActiveDirectory`, `Microsoft.AzureCosmosDB`, `Microsoft.ContainerRegistry`, `Microsoft.EventHub`, `Microsoft.KeyVault`, `Microsoft.ServiceBus`, `Microsoft.Sql`, `Microsoft.Storage` and `Microsoft.Web`.

---

A `delegation` block supports the following:

* `name` (Required) A name for this delegation.

* `service_delegation` (Required) A `service_delegation` block as defined below.

---

A `service_delegation` block supports the following:

-> **NOTE:** Delegating to services may not be available in all regions. Check that the service you are delegating to is available in your region using the [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/network/vnet/subnet?view=azure-cli-latest#az-network-vnet-subnet-list-available-delegations). Also, `actions` is specific to each service type. The exact list of `actions` needs to be retrieved using the aforementioned [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/network/vnet/subnet?view=azure-cli-latest#az-network-vnet-subnet-list-available-delegations).

* `name` - (Required) The name of service to delegate to. Possible values include `Microsoft.ApiManagement/service`, `Microsoft.AzureCosmosDB/clusters`, `Microsoft.BareMetal/AzureVMware`, `Microsoft.BareMetal/CrayServers`, `Microsoft.Batch/batchAccounts`, `Microsoft.ContainerInstance/containerGroups`, `Microsoft.Databricks/workspaces`, `Microsoft.DBforMySQL/flexibleServers`, `Microsoft.DBforMySQL/serversv2`, `Microsoft.DBforPostgreSQL/flexibleServers`, `Microsoft.DBforPostgreSQL/serversv2`, `Microsoft.DBforPostgreSQL/singleServers`, `Microsoft.HardwareSecurityModules/dedicatedHSMs`, `Microsoft.Kusto/clusters`, `Microsoft.Logic/integrationServiceEnvironments`, `Microsoft.MachineLearningServices/workspaces`,  `Microsoft.Netapp/volumes`, `Microsoft.Network/managedResolvers`, `Microsoft.PowerPlatform/vnetaccesslinks`, `Microsoft.ServiceFabricMesh/networks`, `Microsoft.Sql/managedInstances`, `Microsoft.Sql/servers`, `Microsoft.StreamAnalytics/streamingJobs`, `Microsoft.Synapse/workspaces`, `Microsoft.Web/hostingEnvironments`, and `Microsoft.Web/serverFarms`.

* `actions` - (Optional) A list of Actions which should be delegated. This list is specific to the service to delegate to. Possible values include `Microsoft.Network/networkinterfaces/*`, `Microsoft.Network/virtualNetworks/subnets/action`, `Microsoft.Network/virtualNetworks/subnets/join/action`, `Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action` and `Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action`.

-> **NOTE:** Azure may add default actions depending on the service delegation name and they can't be changed.

## Attributes Reference

The following attributes are exported:

* `id` - The subnet ID.
* `name` - The name of the subnet.
* `resource_group_name` - The name of the resource group in which the subnet is created in.
* `virtual_network_name` - The name of the virtual network in which the subnet is created in
* `address_prefix` - (Deprecated) The address prefix for the subnet
* `address_prefixes` - The address prefixes for the subnet

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Subnet.
* `update` - (Defaults to 30 minutes) Used when updating the Subnet.
* `read` - (Defaults to 5 minutes) Used when retrieving the Subnet.
* `delete` - (Defaults to 30 minutes) Used when deleting the Subnet.

## Import

Subnets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_subnet.exampleSubnet /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1
```
