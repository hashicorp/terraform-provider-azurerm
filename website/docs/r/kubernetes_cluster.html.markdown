---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster"
sidebar_current: "docs-azurerm-resource-container-kubernetes-cluster"
description: |-
  Creates a managed Kubernetes Cluster (AKS)
---

# azurerm_kubernetes_cluster

Creates a managed Kubernetes Cluster (AKS)

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).


## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG1"
  location = "West US"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                   = "acctestaks1"
  location               = "${azurerm_resource_group.test.location}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  kubernetes_version     = "1.8.2"
  dns_prefix             = "acctestagent1"
  
  linux_profile {
    admin_username = "acctestuser1"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }
  
  agent_pool_profile {
    name            = "default"
    count           = 1
    vm_size         = "Standard_A0"
    os_type         = "Linux"
  }

  service_principal {
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "00000000000000000000000000000000"
  }

  tags {
    Environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the AKS Managed Cluster instance to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the AKS Managed Cluster instance should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `dns_prefix` - (Optional) DNS prefix specified when creating the managed cluster.

* `kubernetes_version` - (Optional) Version of Kubernetes specified when creating the AKS managed cluster.

* `linux_profile` - (Required) A Linux Profile block as documented below.

* `agent_pool_profile` - (Required) One or more Agent Pool Profile's block as documented below.

* `service_principal` - (Required) A Service Principal block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

`linux_profile` supports the following:

* `admin_username` - (Required) The Admin Username for the Cluster. Changing this forces a new resource to be created.
* `ssh_key` - (Required) An SSH Key block as documented below.

`ssh_key` supports the following:

* `key_data` - (Required) The Public SSH Key used to access the cluster. Changing this forces a new resource to be created.

`agent_pool_profile` supports the following:

* `name` - (Required) Unique name of the Agent Pool Profile in the context of the Subscription and Resource Group.
* `count` - (Required) Number of Agents (VMs) in the Pool. Possible values must be in the range of 1 to 50 (inclusive). Defaults to `1`.
* `vm_size` - (Required) The size of each VM in the Agent Pool (e.g. `Standard_F1`).
* `os_type` - (Optional) The Operating System used for the Agents. Possible values are `Linux` and `Windows`. Defaults to `Linux`.
* `vnet_subnet_id` - (Optional) The ID of the Subnet where the Agents in the Pool should be provisioned.

`service_principal` supports the following:

* `client_id` - (Required) The Client ID for the Service Principal.
* `client_secret` - (Required) The Client Secret for the Service Principal.

## Attributes Reference

The following attributes are exported:

* `id` - The Kubernetes Managed Cluster ID.

* `agent_pool_profile.#.fqdn` - The FQDN of the Azure Kubernetes Managed Cluster.

## Import

Kubernetes Managed Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1
```
