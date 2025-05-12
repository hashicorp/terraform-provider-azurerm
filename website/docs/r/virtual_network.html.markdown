---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network"
description: |-
  Manages a virtual network including any configured subnets. Each subnet can optionally be configured with a security group to be associated with the subnet.
---

# azurerm_virtual_network

Manages a virtual network including any configured subnets. Each subnet can
optionally be configured with a security group to be associated with the subnet.

~> **Note:** Terraform currently provides both a standalone [Subnet resource](subnet.html), and allows for Subnets to be defined in-line within the [Virtual Network resource](virtual_network.html).
At this time you cannot use a Virtual Network with in-line Subnets in conjunction with any Subnet resources. Doing so will cause a conflict of Subnet configurations and will overwrite subnets.

~> **Note:** Terraform currently provides both a standalone [virtual network DNS Servers resource](virtual_network_dns_servers.html), and allows for DNS servers to be defined in-line within the [Virtual Network resource](virtual_network.html).
At this time you cannot use a Virtual Network with in-line DNS servers in conjunction with any Virtual Network DNS Servers resources. Doing so will cause a conflict of Virtual Network DNS Servers configurations and will overwrite virtual networks DNS servers.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_security_group" "example" {
  name                = "example-security-group"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]

  subnet {
    name             = "subnet1"
    address_prefixes = ["10.0.1.0/24"]
  }

  subnet {
    name             = "subnet2"
    address_prefixes = ["10.0.2.0/24"]
    security_group   = azurerm_network_security_group.example.id
  }

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the virtual network. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the virtual network. Changing this forces a new resource to be created.

* `address_space` - (Required) The address space that is used the virtual network. You can supply more than one address space.

* `location` - (Required) The location/region where the virtual network is created. Changing this forces a new resource to be created.

---

* `bgp_community` - (Optional) The BGP community attribute in format `<as-number>:<community-value>`.

-> **Note:** The `as-number` segment is the Microsoft ASN, which is always `12076` for now.

* `ddos_protection_plan` - (Optional) A `ddos_protection_plan` block as documented below.

* `encryption` - (Optional) A `encryption` block as defined below.

* `dns_servers` - (Optional) List of IP addresses of DNS servers

-> **Note:** Since `dns_servers` can be configured both inline and via the separate `azurerm_virtual_network_dns_servers` resource, we have to explicitly set it to empty slice (`[]`) to remove it.

* `edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Virtual Network should exist. Changing this forces a new Virtual Network to be created.

* `flow_timeout_in_minutes` - (Optional) The flow timeout in minutes for the Virtual Network, which is used to enable connection tracking for intra-VM flows. Possible values are between `4` and `30` minutes.

* `subnet` - (Optional) Can be specified multiple times to define multiple subnets. Each `subnet` block supports fields documented below.

-> **Note:** Since `subnet` can be configured both inline and via the separate `azurerm_subnet` resource, we have to explicitly set it to empty slice (`[]`) to remove it.

* `private_endpoint_vnet_policies` - (Optional) The Private Endpoint VNet Policies for the Virtual Network. Possible values are `Disabled` and `Basic`. Defaults to `Disabled`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `ddos_protection_plan` block supports the following:

* `id` - (Required) The ID of DDoS Protection Plan.

* `enable` - (Required) Enable/disable DDoS Protection Plan on Virtual Network.

---

A `encryption` block supports the following:

* `enforcement` - (Required) Specifies if the encrypted Virtual Network allows VM that does not support encryption. Possible values are `DropUnencrypted` and `AllowUnencrypted`.

-> **Note:** Currently `AllowUnencrypted` is the only supported value for the `enforcement` property as `DropUnencrypted` is not yet in public preview or general availability. Please see the [official documentation](https://learn.microsoft.com/en-us/azure/virtual-network/virtual-network-encryption-overview#limitations) for more information.

---

The `subnet` block supports:

* `name` - (Required) The name of the subnet.

* `address_prefixes` - (Required) The address prefixes to use for the subnet.

* `security_group` - (Optional) The Network Security Group to associate with the subnet. (Referenced by `id`, i.e. `azurerm_network_security_group.example.id`)

* `default_outbound_access_enabled` - (Optional) Enable default outbound access to the internet for the subnet. Defaults to `true`.

* `delegation` - (Optional) One or more `delegation` blocks as defined below.

* `private_endpoint_network_policies` - (Optional) Enable or Disable network policies for the private endpoint on the subnet. Possible values are `Disabled`, `Enabled`, `NetworkSecurityGroupEnabled` and `RouteTableEnabled`. Defaults to `Disabled`.

-> **Note:** If you don't want to use network policies like user-defined Routes and Network Security Groups, you need to set `private_endpoint_network_policies` in the subnet to `Disabled`. This setting only applies to Private Endpoints in the Subnet and affects all Private Endpoints in the Subnet.

-> **Note:** If you want to use network policies like user-defined Routes and Network Security Groups, you need to set the `private_endpoint_network_policies` in the Subnet to `Enabled`/`NetworkSecurityGroupEnabled`/`RouteTableEnabled`. This setting only applies to Private Endpoints in the Subnet and affects all Private Endpoints in the Subnet.

-> **Note:** See more details from [Manage network policies for Private Endpoints](https://learn.microsoft.com/en-gb/azure/private-link/disable-private-endpoint-network-policy?tabs=network-policy-portal).

* `private_link_service_network_policies_enabled` - (Optional) Enable or Disable network policies for the private link service on the subnet. Defaults to `true`.

-> **Note:** When configuring Azure Private Link service, the explicit setting `private_link_service_network_policies_enabled` must be set to `false` in the subnet since Private Link Service does not support network policies like user-defined Routes and Network Security Groups. This setting only affects the Private Link service. For other resources in the subnet, access is controlled based on the Network Security Group which can be configured using the `azurerm_subnet_network_security_group_association` resource. See more details from [Manage network policies for Private Link Services](https://learn.microsoft.com/en-gb/azure/private-link/disable-private-link-service-network-policy?tabs=private-link-network-policy-powershell).

* `route_table_id` - (Optional) The ID of the Route Table that should be associated with this subnet.

* `service_endpoints` - (Optional) The list of Service endpoints to associate with the subnet. Possible values include: `Microsoft.AzureActiveDirectory`, `Microsoft.AzureCosmosDB`, `Microsoft.ContainerRegistry`, `Microsoft.EventHub`, `Microsoft.KeyVault`, `Microsoft.ServiceBus`, `Microsoft.Sql`, `Microsoft.Storage`, `Microsoft.Storage.Global` and `Microsoft.Web`.

* `service_endpoint_policy_ids` - (Optional) The list of IDs of Service Endpoint Policies to associate with the subnet.

---

A `delegation` block supports the following:

* `name` - (Required) A name for this delegation.

* `service_delegation` - (Required) A `service_delegation` block as defined below.

---

A `service_delegation` block supports the following:

-> **Note:** Delegating to services may not be available in all regions. Check that the service you are delegating to is available in your region using the [Azure CLI](https://docs.microsoft.com/cli/azure/network/vnet/subnet?view=azure-cli-latest#az-network-vnet-subnet-list-available-delegations). Also, `actions` is specific to each service type. The exact list of `actions` needs to be retrieved using the aforementioned [Azure CLI](https://docs.microsoft.com/cli/azure/network/vnet/subnet?view=azure-cli-latest#az-network-vnet-subnet-list-available-delegations).

* `name` - (Required) The name of service to delegate to. Possible values are `GitHub.Network/networkSettings`, `Informatica.DataManagement/organizations`, `Microsoft.ApiManagement/service`, `Microsoft.Apollo/npu`, `Microsoft.App/environments`, `Microsoft.App/testClients`, `Microsoft.AVS/PrivateClouds`, `Microsoft.AzureCosmosDB/clusters`, `Microsoft.BareMetal/AzureHostedService`, `Microsoft.BareMetal/AzureHPC`, `Microsoft.BareMetal/AzurePaymentHSM`, `Microsoft.BareMetal/AzureVMware`, `Microsoft.BareMetal/CrayServers`, `Microsoft.BareMetal/MonitoringServers`, `Microsoft.Batch/batchAccounts`, `Microsoft.CloudTest/hostedpools`, `Microsoft.CloudTest/images`, `Microsoft.CloudTest/pools`, `Microsoft.Codespaces/plans`, `Microsoft.ContainerInstance/containerGroups`, `Microsoft.ContainerService/managedClusters`, `Microsoft.ContainerService/TestClients`, `Microsoft.Databricks/workspaces`, `Microsoft.DBforMySQL/flexibleServers`, `Microsoft.DBforMySQL/servers`, `Microsoft.DBforMySQL/serversv2`, `Microsoft.DBforPostgreSQL/flexibleServers`, `Microsoft.DBforPostgreSQL/serversv2`, `Microsoft.DBforPostgreSQL/singleServers`, `Microsoft.DelegatedNetwork/controller`, `Microsoft.DevCenter/networkConnection`, `Microsoft.DevOpsInfrastructure/pools`, `Microsoft.DocumentDB/cassandraClusters`, `Microsoft.Fidalgo/networkSettings`, `Microsoft.HardwareSecurityModules/dedicatedHSMs`, `Microsoft.Kusto/clusters`, `Microsoft.LabServices/labplans`, `Microsoft.Logic/integrationServiceEnvironments`, `Microsoft.MachineLearningServices/workspaces`, `Microsoft.Netapp/volumes`, `Microsoft.Network/applicationGateways`, `Microsoft.Network/dnsResolvers`, `Microsoft.Network/managedResolvers`, `Microsoft.Network/fpgaNetworkInterfaces`, `Microsoft.Network/networkWatchers.`, `Microsoft.Network/virtualNetworkGateways`, `Microsoft.Orbital/orbitalGateways`, `Microsoft.PowerAutomate/hostedRpa`, `Microsoft.PowerPlatform/enterprisePolicies`, `Microsoft.PowerPlatform/vnetaccesslinks`, `Microsoft.ServiceFabricMesh/networks`, `Microsoft.ServiceNetworking/trafficControllers`, `Microsoft.Singularity/accounts/networks`, `Microsoft.Singularity/accounts/npu`, `Microsoft.Sql/managedInstances`, `Microsoft.Sql/managedInstancesOnebox`, `Microsoft.Sql/managedInstancesStage`, `Microsoft.Sql/managedInstancesTest`, `Microsoft.Sql/servers`, `Microsoft.StoragePool/diskPools`, `Microsoft.StreamAnalytics/streamingJobs`, `Microsoft.Synapse/workspaces`, `Microsoft.Web/hostingEnvironments`, `Microsoft.Web/serverFarms`, `NGINX.NGINXPLUS/nginxDeployments`, `PaloAltoNetworks.Cloudngfw/firewalls`, `Qumulo.Storage/fileSystems`, and `Oracle.Database/networkAttachments`.

* `actions` - (Optional) A list of Actions which should be delegated. This list is specific to the service to delegate to. Possible values are `Microsoft.Network/networkinterfaces/*`, `Microsoft.Network/publicIPAddresses/join/action`, `Microsoft.Network/publicIPAddresses/read`, `Microsoft.Network/virtualNetworks/read`, `Microsoft.Network/virtualNetworks/subnets/action`, `Microsoft.Network/virtualNetworks/subnets/join/action`, `Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action`, and `Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action`.

-> **Note:** Azure may add default actions depending on the service delegation name and they can't be changed.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The virtual NetworkConfiguration ID.

* `name` - (Required) The name of the virtual network. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the virtual network.

* `location` - (Required) The location/region where the virtual network is created. Changing this forces a new resource to be created.

* `address_space` - (Required) The list of address spaces used by the virtual network.

* `guid` - The GUID of the virtual network.

* `subnet` - One or more `subnet` blocks as defined below.

---

The `subnet` block exports:

* `id` - The ID of this subnet.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Network.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Network.

## Import

Virtual Networks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_network.exampleNetwork /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1
```
