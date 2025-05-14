---
subcategory: "Azure VMware Solution"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vmware_netapp_volume_attachment"
description: |-
  Manages an Azure VMware Solution Private Cloud Netapp File Attachment.
---

# azurerm_vmware_netapp_volume_attachment

Manages an Azure VMware Solution Private Cloud Netapp File Attachment.

## Example Usage

~> **Note:** For Azure Azure VMware Solution Private Cloud, normal `terraform apply` could ignore this note. Please disable correlation request id for continuous operations in one build (like acctest). The continuous operations like `update` or `delete` could not be triggered when it shares the same `correlation-id` with its previous operation.

```hcl
provider "azurerm" {
  features {}
  disable_correlation_request_id = true
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_public_ip" "test" {
  name                = "example-public-ip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network" "test" {
  name                = "example-VirtualNetwork"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.88.0.0/16"]
}

resource "azurerm_subnet" "netappSubnet" {
  name                 = "example-Subnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.88.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet" "gatewaySubnet" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.88.1.0/24"]
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "example-vnet-gateway"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type = "ExpressRoute"
  sku  = "Standard"

  ip_configuration {
    name                 = "vnetGatewayConfig"
    public_ip_address_id = azurerm_public_ip.test.id
    subnet_id            = azurerm_subnet.gatewaySubnet.id
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "example-NetAppAccount"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_netapp_pool" "test" {
  name                = "example-NetAppPool"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 4
}

resource "azurerm_netapp_volume" "test" {
  name                            = "example-NetAppVolume"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  account_name                    = azurerm_netapp_account.test.name
  pool_name                       = azurerm_netapp_pool.test.name
  volume_path                     = "my-unique-file-path-%d"
  service_level                   = "Standard"
  subnet_id                       = azurerm_subnet.netappSubnet.id
  protocols                       = ["NFSv3"]
  storage_quota_in_gb             = 100
  azure_vmware_data_store_enabled = true

  export_policy_rule {
    rule_index          = 1
    allowed_clients     = ["0.0.0.0/0"]
    protocols_enabled   = ["NFSv3"]
    unix_read_only      = false
    unix_read_write     = true
    root_access_enabled = true
  }
}

resource "azurerm_vmware_private_cloud" "test" {
  name                = "example-PC"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "av36"

  management_cluster {
    size = 3
  }
  network_subnet_cidr = "192.168.48.0/22"
}

resource "azurerm_vmware_cluster" "test" {
  name               = "example-vm-cluster"
  vmware_cloud_id    = azurerm_vmware_private_cloud.test.id
  cluster_node_count = 3
  sku_name           = "av36"
}

resource "azurerm_vmware_express_route_authorization" "test" {
  name             = "example-VmwareAuthorization"
  private_cloud_id = azurerm_vmware_private_cloud.test.id
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                = "example-vnetgwconn"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "ExpressRoute"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  express_route_circuit_id   = azurerm_vmware_private_cloud.test.circuit[0].express_route_id
  authorization_key          = azurerm_vmware_express_route_authorization.test.express_route_authorization_key
}

resource "azurerm_vmware_netapp_volume_attachment" "test" {
  name              = "example-vmwareattachment"
  netapp_volume_id  = azurerm_netapp_volume.test.id
  vmware_cluster_id = azurerm_vmware_cluster.test.id

  depends_on = [azurerm_virtual_network_gateway_connection.test]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure VMware Solution Private Cloud Netapp File Volume Attachment. Changing this forces a new Azure VMware Solution Private Cloud Netapp File Volume Attachment to be created.

* `netapp_volume_id` - (Required) The netapp file volume for this Azure VMware Solution Private Cloud Netapp File Volume Attachment to connect to. Changing this forces a new Azure VMware Solution Private Cloud Netapp File Volume Attachment to be created.

* `vmware_cluster_id` - (Required) The vmware cluster for this Azure VMware Solution Private Cloud Netapp File Volume Attachment to associated to. Changing this forces a new Azure VMware Solution Private Cloud Netapp File Volume Attachment to be created.

~> **Note:** please follow the prerequisites mentioned in this [article](https://learn.microsoft.com/en-us/azure/azure-vmware/attach-azure-netapp-files-to-azure-vmware-solution-hosts?tabs=azure-portal#prerequisites) before associating the netapp file volume to the Azure VMware Solution hosts.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure VMware Solution Private Cloud Netapp File Volume Attachment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure VMware Solution Private Cloud Netapp File Volume Attachment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure VMware Solution Private Cloud Netapp File Volume Attachment.

## Import

Azure VMware Solution Private Cloud Netapp File Volume Attachments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vmware_netapp_volume_attachment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/privateCloud1/clusters/Cluster1/dataStores/datastore1
```
