---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_cluster"
sidebar_current: "docs-azurerm-resource-managed-cluster"
description: |-
  Creates an AKS Managed Cluster instance.
---

# azurerm\_managed\_cluster

Creates an AKS Managed Cluster instance

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).


## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG1"
  location = "West US"
}

resource "azurerm_managed_cluster" "test" {
  name                   = "acctestaks1"
  location               = "${azurerm_resource_group.test.location}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  kubernetes_version     = "1.8.2"

  linux_profile {
    admin_username = "acctestuser1"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  agent_pool_profile {
    name            = "default"
    count           = 1
    dns_prefix      = "acctestagent1"
    vm_size         = "Standard_A0"
    storage_profile = "ManagedDisks"
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

* `dns_prefix` - (Required) DNS prefix specified when creating the managed cluster.

* `kubernetes_version` - (Optional) Version of Kubernetes specified when creating the AKS managed cluster.

* `linux_profile` - (Required) A Linux Profile block as documented below.

* `agent_pool_profile` - (Required) One or more Agent Pool Profile's block as documented below.

* `service_principal` - (Required) A Service Principal block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

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
* `storage_profile` - (Optional) Storage profile specifies what kind of storage used. Choose from StorageAccount and ManagedDisks. Leave it empty, we will choose for you based on the orchestrator choice.
* `os_type` - (Optional) OsType to be used to specify os type. Choose from Linux and Windows. Default to Linux.

`service_principal` supports the following:

* `client_id` - (Required) The ID for the Service Principal.
* `client_secret` - (Required) The secret password associated with the service principal.

## Attributes Reference

The following attributes are exported:

* `id` - The AKS Managed Cluster ID.

* `fqdn` - FDQN for the master pool.