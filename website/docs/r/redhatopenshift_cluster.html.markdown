---
subcategory: "Red Hat Openshift"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redhatopenshift_cluster"
description: |-
  Manages fully managed Azure Red Hat Openshift Cluster (also known as ARO)
---

# azurerm_redhatopenshift_cluster

Manages a fully managed Azure Red Hat Openshift Cluster (also known as ARO).

-> **Note:** Due to the fast-moving nature of ARO, we recommend using the latest version of the Azure Provider when using ARO - you can find [the latest version of the Azure Provider here](https://registry.terraform.io/providers/hashicorp/azurerm/latest).

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

This example provisions a basic Azure Red Hat Openshift Cluster.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_virtual_network" "example" {
  name                = "aro-vnet"
  address_space       = ["10.0.0.0/22"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "master_subnet" {
  name                 = "master-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/23"]
  service_endpoints    = ["Microsoft.ContainerRegistry"]
  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "worker_subnet" {
  name                 = "worker-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/23"]
  service_endpoints    = ["Microsoft.ContainerRegistry"]
}

resource "azurerm_redhatopenshift_cluster" "example" {
  name                = "example-redhatopenshift1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  
  service_principal {
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "MyCl1eNtSeCr3t"
  }

  master_profile {
    vm_size = "Standard_D8s_v3"
    subnet_id = azurerm_subnet.master_subnet.id
  }
  
  worker_profile {
    vm_size      = "Standard_D4s_v3"
    subnet_id      = azurerm_subnet.worker_subnet.id
  }

  tags = {
    Environment = "Production"
  }
}

output "openshift_version" {
    value = azurerm_redhatopenshift_cluster.example.version
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Red Hat Openshift Cluster to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Azure Red Hat Openshift Cluster should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Azure Red Hat Openshift Cluster should exist. Changing this forces a new resource to be created.

* `service_principal` - (Required) A `service_principal` block as defined below.

* `master_profile` - (Required) A `master_profile` block as defined below.

* `worker_profile` - (Required) A `worker_profile` block as defined below.

* `cluster_profile` - (Optional) A `cluster_profile` block as defined below.

* `network_profile` - (Optional) A `network_profile` block as defined below.

* `api_server_profile` - (Optional) A `api_server_profile` block as defined below.

* `ingress_profile` - (Optional) A `ingress_profile` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `service_principal` block supports the following:

* `client_id` - (Required) The Client ID for the Service Principal.

* `client_secret` - (Required) The Client Secret for the Service Principal.

---

A `master_profile` block supports the following:

* `vm_size` - (Required) The size of the Virtual Machines for the master nodes. Currently supported values are `Standard_D2s_v3`, `Standard_D4s_v3` and `Standard_D8s_v3`. Changing this forces a new resource to be created.
* `subnet_id` - (Required) The ID of the subnet where master nodes will be hosted.

-> **NOTE** The subnet which master nodes will be associated must met the following requirements:

  * Private subnet access granted to `Microsoft.ContainerRegistry` service endpoint.
  * Subnet private endpoint policies disabled. For more info, see [Disable network policies for Private Link service source IP](https://docs.microsoft.com/azure/private-link/disable-private-link-service-network-policy).

---

A `worker_profile` block supports the following:

* `name` - (Optional) The worker profile name. Defaults to `worker`. Changing this forces a new resource to be created.
* `vm_size` - (Required) The size of the Virtual Machines for the worker nodes. Currently supported values are `Standard_D2s_v3`, `Standard_D4s_v3` and `Standard_D8s_v3`. Changing this forces a new resource to be created.
* `subnet_id` - (Required) The ID of the subnet where worker nodes will be hosted.
* `disk_size_gb` - (Optional) The internal OS disk size of the worker Virtual Machines in GB. Must be `128` or greater. Defaults to `128`.
* `node_count` - (Optional) The initial number of worker nodes which should exist in the cluster. If specified this must be between `3` and `20`. Defaults to `3`.

-> **NOTE** The subnet which worker nodes will be associated must have private subnet access granted to `Microsoft.ContainerRegistry` service endpoint.

---

A `cluster_profile` block supports the following:

* `pull_secret` - (Optional) The Red Hat pull secret for the cluster.
* `domain` - (Optional) The custom domain for the cluster. Defaults to `<random>.<location>.aroapp.io`. For more info, see [Prepare a custom domain for your cluster](https://docs.microsoft.com/azure/openshift/tutorial-create-cluster#prepare-a-custom-domain-for-your-cluster-optional).
* `resource_group_id` - (Optional) The Red Hat Openshift cluster resource group ID. Defaults to `aro-<random>`.

---

A `network_profile` block supports the following:

* `pod_cidr` - (Optional) The CIDR to use for pod IP addresses. Defaults to `10.128.0.0/1`. Changing this forces a new resource to be created.
* `service_cidr` - (Optional) The network range used by the Openshift service. Defaults to `172.30.0.0/16`. Changing this forces a new resource to be created.

---

A `api_server_profile` block supports the following:

* `visibility` - (Optional) Cluster API server visibility. Supported values are `Public` and `Private`. Defaults to `Public`. Changing this forces a new resource to be created.

---

A `ingress_profile` block supports the following:

* `visibility` - (Optional) Cluster Ingress visibility. Supported values are `Public` and `Private`. Defaults to `Public`. Changing this forces a new resource to be created.

---

## Attributes Reference

The following attributes are exported:

* `version` - The Red Hat Openshift version used by the cluster.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/language/resources/syntax.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Kubernetes Cluster.
* `update` - (Defaults to 90 minutes) Used when updating the Kubernetes Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Cluster.
* `delete` - (Defaults to 90 minutes) Used when deleting the Kubernetes Cluster.

## Import

Red Hat Openshift Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redhatopenshift_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.RedHatOpenShift/openShiftClusters/cluster1
```
