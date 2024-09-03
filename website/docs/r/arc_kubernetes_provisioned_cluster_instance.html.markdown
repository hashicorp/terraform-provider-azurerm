---
subcategory: "ArcKubernetes"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_arc_kubernetes_provisioned_cluster"
description: |-
  Manages an Arc Kubernetes Provisioned Cluster Instance.
---

# azurerm_arc_kubernetes_provisioned_cluster

Manages an Arc Kubernetes Provisioned Cluster Instance.

## Example Usage

```hcl
resource "tls_private_key" "rsaKey" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "azurerm_resource_group" "example" {
  name     = "example-akpci"
  location = "West Europe"
}

resource "azurerm_stack_hci_logical_network" "example" {
  name                = "example-ln"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["192.168.1.254"]

  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "192.168.1.0/24"
    ip_pool {
      start = "192.168.1.171"
      end   = "192.168.1.190"
    }
    route {
      name                = "example-route"
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "192.168.1.1"
    }
  }
}

resource "azurerm_arc_kubernetes_cluster" "example" {
  name                = "example-akc"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  kind                = "ProvisionedCluster"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_arc_kubernetes_provisioned_cluster" "example" {
  cluster_id         = azurerm_arc_kubernetes_cluster.example.id
  custom_location_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"
  kubernetes_version = "1.28.5"

  agent_pool_profile {
    auto_scaling_enabled = true
    count                = 1
    max_count            = 2
    min_count            = 1
    max_pods             = 20
    name                 = "nodepool1"
    os_sku               = "CBLMariner"
    os_type              = "Linux"
    vm_size              = "Standard_A4_v2"
    node_taints          = ["env=prod:NoSchedule"]

    node_labels = {
      foo = "bar"
      env = "dev"
    }
  }

  cloud_provider_profile {
    infra_network_profile {
      vnet_subnet_ids = [azurerm_stack_hci_logical_network.example.id]
    }
  }

  cluster_vm_access_profile {
    authorized_ip_ranges = "192.168.0.1,192.168.0.2"
  }

  control_plane_profile {
    count   = 3
    vm_size = "Standard_A4_v2"
    host_ip = "192.168.1.190"
  }

  license_profile {
    azure_hybrid_benefit = "False"
  }

  linux_profile {
    ssh_key = [tls_private_key.rsaKey.public_key_openssh]
  }

  network_profile {
    network_policy = "calico"
    pod_cidr       = "10.244.0.0/16"
  }

  storage_profile {
    smb_csi_driver_enabled = true
    nfs_csi_driver_enabled = true
  }
}
```

## Arguments Reference

The following arguments are supported:

* `agent_pool_profile` - (Required) One or more `agent_pool_profile` blocks as defined below.

* `cloud_provider_profile` - (Required) A `cloud_provider_profile` block as defined below. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `cluster_id` - (Required) The ID of the Arc Kubernetes Cluster. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `control_plane_profile` - (Required) A `control_plane_profile` block as defined below.

* `custom_location_id` - (Required) The ID of the Custom Location. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `kubernetes_version` - (Required) The Kubernetes version to use for the Arc Kubernetes Provisioned Cluster Instance. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `network_profile` - (Required) A `network_profile` block as defined below. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

---

* `cluster_vm_access_profile` - (Optional) A `cluster_vm_access_profile` block as defined below. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `license_profile` - (Optional) A `license_profile` block as defined below.

* `linux_profile` - (Optional) A `linux_profile` block as defined below. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `storage_profile` - (Optional) A `storage_profile` block as defined below.

---

A `agent_pool_profile` block supports the following:

* `name` - (Required) The name which should be used for this Agent Pool. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `vm_size` - (Required) The VM sku size of the agent pool node VMs. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `os_sku` - (Required) The OS SKU used by the agent pool nodes. Possible values are `CBLMariner`, `Windows2019` and `Windows2022`. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `os_type` - (Required) The OS type for the agent pool nodes. Possible values are `Windows` and `Linux`. Defaults to `Linux`. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `auto_scaling_enabled` - (Optional) Whether to enable auto scaling. Defaults to `false`.

* `count` - (Optional) The number of nodes in the agent pool. Defaults to `1`.

* `max_count` - (Optional) The maximum number of nodes for auto scaling.

* `max_pods` - (Optional) The maximum number of pods that can run on a node.

* `min_count` - (Optional) The minimum number of nodes for auto scaling.

* `node_labels` - (Optional) The node labels to be persisted across all nodes in agent pool.

* `node_taints` - (Optional) The taints added to new nodes during node pool create and scale. For example, `key=value:NoSchedule`.

---

A `cloud_provider_profile` block supports the following:

* `infra_network_profile` - (Required) A `infra_network_profile` block as defined below. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

---

A `cluster_vm_access_profile` block supports the following:

* `authorized_ip_ranges` - (Required) The ranges of IP Address or CIDR for SSH access to VMs in the provisioned cluster. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

---

A `control_plane_profile` block supports the following:

* `vm_size` - (Required) The VM sku size of the control plane nodes. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `count` - (Optional) The number of control plane nodes. The count should be an odd number. Defaults to `1`.

* `host_ip` - (Optional) The IP Address of the Kubernetes API server. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

~> **NOTE:** The `control_plane_profile.host_ip` and `load_balancer_profile` cannot be specified together.

---

A `infra_network_profile` block supports the following:

* `vnet_subnet_ids` - (Required) Specifies a list of ARM resource IDs for the infrastructure network object with `Microsoft.AzureStackHCI/logicalNetworks` or `Microsoft.HybridContainerService/virtualNetworks` resource type.

---

A `license_profile` block supports the following:

* `azure_hybrid_benefit` - (Optional) Whether Azure Hybrid Benefit is opted in. Defaults to `NotApplicable`.

---

A `linux_profile` block supports the following:

* `ssh_key` - (Required) A list of certificate public keys used to authenticate with VMs through SSH. The certificate must be in PEM format with or without headers. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

---

A `load_balancer_profile` block supports the following:

* `count` - (Required) The number of HA Proxy load balancer VMs. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

~> **NOTE:** The `load_balancer_profile` and `control_plane_profile.host_ip` cannot be specified together.

---

A `network_profile` block supports the following:

* `network_policy` - (Required) The network policy used for building Kubernetes network. The only possible value is `calico`. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `pod_cidr` - (Required) A CIDR notation IP Address range from which to assign pod IPs. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

* `load_balancer_profile` - (Optional) A `load_balancer_profile` block as defined above. Changing this forces a new Arc Kubernetes Provisioned Cluster Instance to be created.

---

A `storage_profile` block supports the following:

* `nfs_csi_driver_enabled` - (Optional) Whether to enable the NFS CSI Driver. Defaults to `true`.

* `smb_csi_driver_enabled` - (Optional) Whether to enable the SMB CSI Driver. Defaults to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Arc Kubernetes Provisioned Cluster Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Arc Kubernetes Provisioned Cluster Instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the Arc Kubernetes Provisioned Cluster Instance.
* `update` - (Defaults to 30 minutes) Used when updating the Arc Kubernetes Provisioned Cluster Instance.
* `delete` - (Defaults to 30 minutes) Used when deleting the Arc Kubernetes Provisioned Cluster Instance.

## Import

Arc Kubernetes Provisioned Cluster Instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_arc_kubernetes_provisioned_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Kubernetes/connectedClusters/cluster1/providers/Microsoft.HybridContainerService/provisionedClusterInstances/default
```
