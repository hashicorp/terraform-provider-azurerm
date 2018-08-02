---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster"
sidebar_current: "docs-azurerm-resource-container-kubernetes-cluster"
description: |-
  Manages a managed Kubernetes Cluster (AKS)
---

# azurerm_kubernetes_cluster

Manages a managed Kubernetes Cluster (AKS)

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage - Basic

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG1"
  location = "East US"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "acctestagent1"

  linux_profile {
    admin_username = "acctestuser1"

    ssh_key {
      key_data = "ssh-rsa ..."
    }
  }

  agent_pool_profile {
    name            = "default"
    count           = 1
    vm_size         = "Standard_D1_v2"
    os_type         = "Linux"
    os_disk_size_gb = 30
  }

  service_principal {
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "00000000000000000000000000000000"
  }

  tags {
    Environment = "Production"
  }
}

output "id" {
    value = "${azurerm_kubernetes_cluster.test.id}"
}

output "kube_config" {
  value = "${azurerm_kubernetes_cluster.test.kube_config_raw}"
}

output "client_key" {
  value = "${azurerm_kubernetes_cluster.test.kube_config.0.client_key}"
}

output "client_certificate" {
  value = "${azurerm_kubernetes_cluster.test.kube_config.0.client_certificate}"
}

output "cluster_ca_certificate" {
  value = "${azurerm_kubernetes_cluster.test.kube_config.0.cluster_ca_certificate}"
}

output "host" {
  value = "${azurerm_kubernetes_cluster.test.kube_config.0.host}"
}
```

## Example Usage - Advanced Networking

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG1"
  location = "East US"
}

resource azurerm_network_security_group "test_advanced_network" {
  name                = "akc-1-nsg"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_network" "test_advanced_network" {
  name                = "akc-1-vnet"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test_subnet" {
  name                      = "akc-1-subnet"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  network_security_group_id = "${azurerm_network_security_group.test_advanced_network.id}"
  address_prefix            = "10.1.0.0/24"
  virtual_network_name      = "${azurerm_virtual_network.test_advanced_network.name}"
}

resource "azurerm_kubernetes_cluster" "test" {
  name       = "akc-1"
  location   = "${azurerm_resource_group.test.location}"
  dns_prefix = "akc-1"

  resource_group_name = "${azurerm_resource_group.test.name}"

  linux_profile {
    admin_username = "acctestuser1"

    ssh_key {
      key_data = "ssh-rsa ..."
    }
  }

  agent_pool_profile {
    name    = "agentpool"
    count   = "2"
    vm_size = "Standard_DS2_v2"
    os_type = "Linux"

    # Required for advanced networking
    vnet_subnet_id = "${azurerm_subnet.test_subnet.id}"
  }

  service_principal {
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "00000000000000000000000000000000"
  }

  network_profile {
    network_plugin = "azure"
  }
}

output "subnet_id" {
  value = "${azurerm_kubernetes_cluster.test.agent_pool_profile.0.vnet_subnet_id}"
}

output "network_plugin" {
  value = "${azurerm_kubernetes_cluster.test.network_profile.0.network_plugin}"
}

output "service_cidr" {
  value = "${azurerm_kubernetes_cluster.test.network_profile.0.service_cidr}"
}

output "dns_service_ip" {
  value = "${azurerm_kubernetes_cluster.test.network_profile.0.dns_service_ip}"
}

output "docker_bridge_cidr" {
  value = "${azurerm_kubernetes_cluster.test.network_profile.0.docker_bridge_cidr}"
}

output "pod_cidr" {
  value = "${azurerm_kubernetes_cluster.test.network_profile.0.pod_cidr}"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the AKS Managed Cluster instance to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the AKS Managed Cluster instance should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `dns_prefix` - (Required) DNS prefix specified when creating the managed cluster.

* `kubernetes_version` - (Optional) Version of Kubernetes specified when creating the AKS managed cluster. If not specified, the latest recommended version will be used at provisioning time (but won't auto-upgrade).

* `linux_profile` - (Required) A Linux Profile block as documented below.

* `agent_pool_profile` - (Required) One or more Agent Pool Profile's block as documented below.

* `service_principal` - (Required) A Service Principal block as documented below.

* `network_profile` - (Optional) A Network Profile block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

`linux_profile` supports the following:

* `admin_username` - (Required) The Admin Username for the Cluster. Changing this forces a new resource to be created.
* `ssh_key` - (Required) An SSH Key block as documented below.

`ssh_key` supports the following:

* `key_data` - (Required) The Public SSH Key used to access the cluster. Changing this forces a new resource to be created.

`agent_pool_profile` supports the following:

* `name` - (Required) Unique name of the Agent Pool Profile in the context of the Subscription and Resource Group. Changing this forces a new resource to be created.
* `count` - (Required) Number of Agents (VMs) in the Pool. Possible values must be in the range of 1 to 50 (inclusive). Defaults to `1`.
* `vm_size` - (Required) The size of each VM in the Agent Pool (e.g. `Standard_F1`). Changing this forces a new resource to be created.
* `os_disk_size_gb` - (Optional) The Agent Operating System disk size in GB. Changing this forces a new resource to be created.
* `os_type` - (Optional) The Operating System used for the Agents. Possible values are `Linux` and `Windows`.  Changing this forces a new resource to be created. Defaults to `Linux`.
* `vnet_subnet_id` - (Optional) The ID of the Subnet where the Agents in the Pool should be provisioned. Changing this forces a new resource to be created.

`service_principal` supports the following:

* `client_id` - (Required) The Client ID for the Service Principal.
* `client_secret` - (Required) The Client Secret for the Service Principal.

`network_profile` supports the following:

* `network_plugin` - (Required) Network plugin to use for networking. Currently supported values are 'azure' and 'kubenet'. Changing this forces a new resource to be created.

-> **NOTE:** When `network_plugin` is set to `azure` - the `vnet_subnet_id` field in the `agent_pool_profile` block must be set.

* `service_cidr` - (Optional) The Network Range used by the Kubernetes service. This is required when `network_plugin` is set to `kubenet`. Changing this forces a new resource to be created.

~> **NOTE:** This range should not be used by any network element on or connected to this VNet. Service address CIDR must be smaller than /12.

* `dns_service_ip` - (Optional) IP address within the Kubernetes service address range that will be used by cluster service discovery (kube-dns). This is required when `network_plugin` is set to `kubenet`. Changing this forces a new resource to be created.

* `docker_bridge_cidr` - (Optional) IP address (in CIDR notation) used as the Docker bridge IP address on nodes. This is required when `network_plugin` is set to `kubenet`. Changing this forces a new resource to be created.

* `pod_cidr` - (Optional) The CIDR to use for pod IP addresses. Changing this forces a new resource to be created.

Here's an example of configuring the `kubenet` Networking Profile:

```
resource "azurerm_subnet" "test" {
  # ...
}

resource "azurerm_kubernetes_cluster" "test" {
  # ...

  agent_pool_profile {
    # ...
    vnet_subnet_id = "${azurerm_subnet.test.id}"
  }

  network_profile {
    network_plugin     = "kubenet"
    pod_cidr           = "10.244.0.0/24"
    dns_service_ip     = "10.10.0.10"
    docker_bridge_cidr = "172.17.0.1/16"
    service_cidr       = "10.10.0.0/16"
  }
}
```


[**Find out more about AKS Advanced Networking**](https://docs.microsoft.com/en-us/azure/aks/networking-overview#advanced-networking)

## Attributes Reference

The following attributes are exported:

* `id` - The Kubernetes Managed Cluster ID.

* `fqdn` - The FQDN of the Azure Kubernetes Managed Cluster.

* `node_resource_group` - Auto-generated Resource Group containing AKS Cluster resources.

* `kube_config_raw` - Raw Kubernetes config to be used by
    [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and
    other compatible tools

* `kube_config` - Kubernetes configuration, sub-attributes defined below:

  * `client_key` - Base64 encoded private key used by clients to authenticate to the Kubernetes cluster.

  * `client_certificate` - Base64 encoded public certificate used by clients to authenticate to the Kubernetes cluster.

  * `cluster_ca_certificate` - Base64 encoded public CA certificate used as the root of trust for the Kubernetes cluster.

  * `host` - The Kubernetes cluster server host.

  * `username` - A username used to authenticate to the Kubernetes cluster.

  * `password` - A password or token used to authenticate to the Kubernetes cluster.

-> **NOTE:** It's possible to use these credentials with [the Kubernetes Provider](/docs/providers/kubernetes/index.html) like so:

```
provider "kubernetes" {
  host                   = "${azurerm_kubernetes_cluster.main.kube_config.0.host}"
  username               = "${azurerm_kubernetes_cluster.main.kube_config.0.username}"
  password               = "${azurerm_kubernetes_cluster.main.kube_config.0.password}"
  client_certificate     = "${base64decode(azurerm_kubernetes_cluster.main.kube_config.0.client_certificate)}"
  client_key             = "${base64decode(azurerm_kubernetes_cluster.main.kube_config.0.client_key)}"
  cluster_ca_certificate = "${base64decode(azurerm_kubernetes_cluster.main.kube_config.0.cluster_ca_certificate)}"
}
```

## Import

Kubernetes Managed Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1
```
