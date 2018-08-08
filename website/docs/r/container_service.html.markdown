---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_service"
sidebar_current: "docs-azurerm-resource-container-service"
description: |-
  Manages an Azure Container Service instance.
---

# azurerm_container_service

Manages an Azure Container Service Instance

~> **NOTE:** All arguments including the client secret will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

~> **NOTE:** If you're working with Kubernetes - we'd recommend using [the `azurerm_kubernetes_cluster` resource](kubernetes_cluster.html) instead of this resource.

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_container_service" "example" {
  name                   = "dcoscontainersvc"
  location               = "${azurerm_resource_group.example.location}"
  resource_group_name    = "${azurerm_resource_group.example.name}"
  orchestration_platform = "DCOS"

  master_profile {
    count      = 1
    dns_prefix = "acctestmaster1"
  }

  linux_profile {
    admin_username = "acctestuser1"

    ssh_key {
      key_data = "ssh-rsa public-key-goes-here terraform@demo.tld"
    }
  }

  agent_pool_profile {
    name       = "default"
    count      = 1
    dns_prefix = "acctestagent1"
    vm_size    = "Standard_F2"
  }

  diagnostics_profile {
    enabled = false
  }

  tags {
    Environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Container Service instance to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Container Service instance should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `orchestration_platform` - (Required) Specifies the Container Orchestration Platform to use. Currently can be either `DCOS`, `Kubernetes` or `Swarm`. Changing this forces a new resource to be created.

~> **NOTE:** If you're working with Kubernetes - we'd recommend using [the `azurerm_kubernetes_cluster` resource](kubernetes_cluster.html) instead of this resource.

* `master_profile` - (Required) A `master_profile` block as documented below.

* `linux_profile` - (Required) A `linux_profile` block as documented below.

* `agent_pool_profile` - (Required) One or more `agent_pool_profile` blocks as documented below.

* `service_principal` - (only Required when you're using `Kubernetes` as an Orchestration Platform) A Service Principal block as documented below.

* `diagnostics_profile` - (Required) A `diagnostics_profile` block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.


`master_profile` supports the following:

* `count` - (Required) Number of masters (VMs) in the container service cluster. Allowed values are 1, 3, and 5. The default value is 1.
* `dns_prefix` - (Required) The DNS Prefix to use for the Container Service master nodes.

`linux_profile` supports the following:

* `admin_username` - (Required) The Admin Username for the Cluster.
* `ssh_key` - (Required) An SSH Key block as documented below.

`ssh_key` supports the following:

* `key_data` - (Required) The Public SSH Key used to access the cluster.

`agent_pool_profile` supports the following:

* `name` - (Required) Unique name of the agent pool profile in the context of the subscription and resource group.
* `count` - (Required) Number of agents (VMs) to host docker containers. Allowed values must be in the range of 1 to 100 (inclusive). The default value is 1.
* `dns_prefix` - (Required) The DNS Prefix given to Agents in this Agent Pool.
* `vm_size` - (Required) The VM Size of each of the Agent Pool VM's (e.g. Standard_F1 / Standard_D2v2).

`service_principal` supports the following:

* `client_id` - (Required) The ID for the Service Principal.
* `client_secret` - (Required) The secret password associated with the service principal.

`diagnostics_profile` supports the following:

* `enabled` - (Required) Should VM Diagnostics be enabled for the Container Service VM's

## Attributes Reference

The following attributes are exported:

* `id` - The Container Service ID.

* `master_profile.0.fqdn` - FDQN for the master.

* `agent_pool_profile.0.fqdn` - FDQN for the agent pool.

* `diagnostics_profile.0.storage_uri` - The URI of the storage account where diagnostics are stored.
